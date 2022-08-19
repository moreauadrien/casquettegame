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

func (c *Client) sendResponse(to string, data events.ResponseData) {
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
