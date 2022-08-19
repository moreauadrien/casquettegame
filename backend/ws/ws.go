package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"timesup/events"
	"timesup/payloads"

	"github.com/gorilla/websocket"
)

var Wrapper = wrapper{eventHandlers: map[string]EventHandler{}, responseHandlers: map[string]ResponseHandler{}}

type EventHandler func(Client, events.EventData) *events.ResponseData
type ResponseHandler func(Client, events.ResponseData)

type wrapper struct {
	eventHandlers    map[string]EventHandler
	responseHandlers map[string]ResponseHandler
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wr *wrapper) On(eventName string, handler EventHandler) {
	if _, exist := wr.eventHandlers[eventName]; exist {
		panic(fmt.Errorf("\"%v\" already has a handler ", eventName))
	}

	wr.eventHandlers[eventName] = handler
	log.Printf("\"%v\" handler registered", eventName)
}

func (wr *wrapper) OnResponse(responseId string, handler ResponseHandler) {
	if _, exist := wr.responseHandlers[responseId]; exist {
		panic(fmt.Errorf("message \"%v\" already has a handler ", responseId))
	}

	wr.responseHandlers[responseId] = handler

	go func() {
		time.Sleep(20 * time.Second)
		delete(wr.responseHandlers, responseId)
	}()
}

func (wr *wrapper) handleIncomingMessage(c Client, msgType int, msg []byte) {
	if msgType != websocket.TextMessage {
		log.Printf("message type %v is not supported", msgType)
		return
	}

	p := new(payloads.GenericPayload)
	if err := p.UnmarshalFrom(msg); err != nil {
		c.sendRawError(err)
		return
	}

	if p.Event.Type == "response" {
		p := new(payloads.ResponsePayload)

		if err := json.Unmarshal(msg, p); err != nil {
			c.sendRawError(err)
			return
		}
		handler := wr.responseHandlers[p.Event.To]

		if handler != nil {
			handler(c, p.Event.Data)
			delete(wr.responseHandlers, p.Event.To)
		}
	} else {
		p := new(payloads.NormalEventPayload)

		if err := json.Unmarshal(msg, p); err != nil {
			c.sendRawError(err)
			return
		}

		handler := wr.eventHandlers[p.Event.Type]

		if handler != nil {
			resp := handler(c, p.Event.Data)
			if resp != nil {
				c.sendResponse(p.MessageId, *resp)
			}
		}
	}
}

func (wr *wrapper) HttpHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()

		if err != nil {
			break
		}

		wr.handleIncomingMessage(New(conn, *wr), t, msg)
	}
}
