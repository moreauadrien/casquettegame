package ws

import (
	"fmt"
	"log"
	"timesup/events"
	"timesup/payloads"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	wrapper wrapper
}

func New(conn *websocket.Conn, w wrapper) Client {
	return Client{conn, w}
}

func (c *Client) sendRawError(e error) {
	msg := fmt.Sprintf("ERROR: %v", e.Error())
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Printf("failed to send error: %v", e.Error())
	}
}

func (c *Client) sendResponse(to string, data events.EventData) {
	p := payloads.ResponsePayload{
		Event: events.ResponseEvent{
			Type: "response",
			To:   to,
			Data: data,
		},
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, p.Marshal()); err != nil {
		log.Printf("the response to %v could not be sent", to)
		return
	}
}

func (c *Client) SendEvent(eventType string, eventData events.EventData, handler ResponseHandler) error {
	p := payloads.GenericPayload{
		MessageId: uuid.NewString(),
		Event: &events.GenericEvent{
			Type: eventType,
			Data: eventData,
		},
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, p.Marshal()); err != nil {
		return fmt.Errorf("the message %v could not be sent", p)
	}

	c.wrapper.OnResponse(p.MessageId, handler)

	return nil
}

/*
import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type PlayerInfos struct {
	Username string `json:"username"`
	Id       string `json:"id"`
	Team     int    `json:"team"`
}

type RoomInfos struct {
	Host    string        `json:"host"`
	Players []PlayerInfos `json:"players"`
}

type ResponseData struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	RoomInfos RoomInfos `json:"roomInfos"`
}

type Message struct {
	MessageId string `json:"messageId"`
	Event     GeEvent  `json:"event"`
}

func (m Message) Marshal() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return b
}

func newMessage(event Event) Message {
	return Message{
		Event:     event,
		MessageId: uuid.NewString(),
	}
}

type ResponseEvent struct {
	Type string       `json:"type"`
	To   string       `json:"to"`
	Data ResponseData `json:"data,omitempty"`
}

type Response struct {
	Event ResponseEvent `json:"event"`
}

func (r Response) Marshal() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return b
}

func newResponse(to string, data ResponseData) Response {
	return Response{
		Event: ResponseEvent{
			Type: "response",
			To:   to,
			Data: data,
		},
	}
}

type Client struct {
	conn    *websocket.Conn
	wrapper wrapper
}

func New(conn *websocket.Conn, w wrapper) Client {
	return Client{conn, w}
}

func (c *Client) SendRawError(e error) {
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(e.Error())); err != nil {
		log.Printf("failed to send error: %v", e.Error())
	}
}

func (c *Client) SendEvent(event Event, handler ResponseHandler) error {
	m := newMessage(event)
	if err := c.conn.WriteMessage(websocket.TextMessage, m.Marshal()); err != nil {
		return fmt.Errorf("the message %v could not be sent", m)
	}

	c.wrapper.OnResponse(m.MessageId, handler)

	return nil
}

func (c *Client) sendResponse(to string, data ResponseData) {
	resp := newResponse(to, data)

	if err := c.conn.WriteMessage(websocket.TextMessage, resp.Marshal()); err != nil {
		log.Printf("the response to %v could not be sent", to)
		return
	}
}
*/
