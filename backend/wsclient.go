package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"timesup/game"

	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type EventData struct {
	RoomId string `json:"roomId"`
}

type Event struct {
	Type string    `json:"type"`
	Data EventData `json:"data"`
}

type Payload struct {
	PlayerUsername string `json:"username"`
	PlayerId       string `json:"id"`
	PlayerToken    string `json:"token"`
	Event          Event  `json:"event"`
}

func buildErrorPayload(msg string) error {
	return fmt.Errorf("{\"type\":\"error\",\"message\":\"%s\"}", msg)
}

func parsePayload(b []byte) (*Payload, error) {
	p := Payload{}
	if err := json.Unmarshal(b, &p); err != nil {
		fmt.Println(err.Error())
		return nil, buildErrorPayload("the payload could not be parsed")
	}

	if len(p.PlayerToken) == 0 {
		return nil, buildErrorPayload("token is needed")
	}

	if len(p.PlayerUsername) == 0 {
		return nil, buildErrorPayload("username is needed")
	}

	return &p, nil
}

func handlePayload(p Payload, conn *websocket.Conn) {
	switch p.Event.Type {
	case "join":
		r, exist := game.Rooms[p.Event.Data.RoomId]
		if !exist {
			conn.WriteMessage(websocket.TextMessage, []byte(buildErrorPayload("this room does not exist").Error()))
			return
		}

		player, exist := game.Players[p.PlayerToken]

		if !exist {
			player = &game.Player{Token: p.PlayerToken, Username: p.PlayerUsername, Id: p.PlayerId}
			game.Players[p.PlayerToken] = player
		}

		player.SetWsConn(conn)

		if player.Room() != r {
			player.SetRoom(r)
			r.AddPlayer(player)

			log.Printf("%v joined the room %v", player.Username, r.Id)
		} else {
			player.SendRoomInfos()
			log.Printf("%v re-joined the room %v", player.Username, r.Id)
		}

	case "start":
		r, exist := game.Rooms[p.Event.Data.RoomId]
		if !exist {
			conn.WriteMessage(websocket.TextMessage, []byte(buildErrorPayload("this room does not exist").Error()))
			return
		}

		player, _ := game.Players[p.PlayerToken]

		if r.Host() != player {
			conn.WriteMessage(websocket.TextMessage, []byte(buildErrorPayload("your not the owner of this room").Error()))
			return
		}

		if r.State() != game.WaitingRoom {
			return
		}

		err := r.StartGame()
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(buildErrorPayload(err.Error()).Error()))
		}
	}
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			break
		}

		p, err := parsePayload(msg)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			continue
		}

		handlePayload(*p, conn)
	}
}
