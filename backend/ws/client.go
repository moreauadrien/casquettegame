package ws

import (
	"encoding/json"
	"log"
	"timesup/game"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Message struct {
	MessageId string     `json:"messageId"`
	Event     game.Event `json:"event"`
}

func (m Message) Marshal() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return b
}

func newMessage(event game.Event) Message {
	return Message{
		Event:     event,
		MessageId: uuid.New().String(),
	}
}

type Response struct {
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

func (c *Client) SendEvent(event game.Event, handler ResponseHandler) {
	m := newMessage(event)
	if err := c.conn.WriteMessage(websocket.TextMessage, m.Marshal()); err != nil {
		log.Printf("the message %v could not be sent", m)
		return
	}

	c.wrapper.OnResponse(m.MessageId, handler)
}
