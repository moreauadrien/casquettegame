package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"timesup/ws/client"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

var Wrapper = wrapper{map[string]EventHandler{}}

type EventHandler func(client.Client, *Payload)

type wrapper struct {
	eventHandlers map[string]EventHandler
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Event struct {
	Type string `json:"type" validate:"required,min=3"`
	To   string `json:"to"`
	Data string `json:"data"`
}

type Credentials struct {
	Token    string `json:"token" validate:"required,min=3"`
	Id       string `json:"id" validate:"required,min=3"`
	Username string `json:"username" validate:"required,min=3"`
}

type Payload struct {
	MessageId   string       `json:"messageId"`
	Event       *Event       `json:"event" validate:"required"`
	Credentials *Credentials `json:"credentials"`
}

var payloadValidator = validator.New()
var credentialsValidator = validator.New()

func (p *Payload) unmarshalFrom(msg []byte) error {
	err := json.Unmarshal(msg, p)

	if err != nil {
		return err
	}

	if err := payloadValidator.Struct(p); err != nil {
		return err
	}

	if p.Event.Type == "response" {
		if len(p.Event.To) == 0 {
			return fmt.Errorf("\"to field is required on a response event")
		}
	} else {
		if p.Credentials == nil {
			return fmt.Errorf("\"credentials\" field is required")
		}

		if err := credentialsValidator.Struct(*p.Credentials); err != nil {
			return err
		}

		if len(p.MessageId) == 0 {
			return fmt.Errorf("\"messageId\" field is required")
		}
	}

	return nil
}

func (wr *wrapper) On(eventName string, handler EventHandler) {
	if _, exist := wr.eventHandlers[eventName]; exist {
		panic(fmt.Errorf("\"%v\" already has a handler ", eventName))
	}

	wr.eventHandlers[eventName] = handler
	log.Printf("\"%v\" handler registered", eventName)
}

func (wr *wrapper) handleIncomingMessage(c client.Client, msgType int, msg []byte) {
	if msgType != websocket.TextMessage {
		log.Printf("message type %v is not supported", msgType)
		return
	}

	p := new(Payload)
	if err := p.unmarshalFrom(msg); err != nil {
		c.SendRawError(err)
		return
	}

	eventType := p.Event.Type
	if eventType == "response" {
	} else {
		handler := wr.eventHandlers[eventType]

		if handler != nil {
			handler(c, p)
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

		wr.handleIncomingMessage(client.New(conn), t, msg)
	}
}
