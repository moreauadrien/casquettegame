package client

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func New(conn *websocket.Conn) Client {
	return Client{conn}
}

func (c *Client) SendRawError(e error) {
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(e.Error())); err != nil {
		log.Printf("failed to send error: %v", e.Error())
	}
}
