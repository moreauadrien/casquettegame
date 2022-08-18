package game

import (
	"log"

	"github.com/gorilla/websocket"
)

var Players = make(map[string]*Player)

type PlayerInfos struct {
	Username string `json:"username"`
	Id       string `json:"id"`
	Team     int    `json:"team"`
}

type Player struct {
	Id       string `json:"id"`
	Token    string `json:"token"`
	Username string `json:"username"`
	team     int
	room     *Room
	wsConn   *websocket.Conn
}

func (p Player) Room() *Room {
	return p.room
}

func (p Player) SendRoomInfos() {
	infosEvent := InfosEvent{
		Players: p.room.Players(),
		Host:    p.room.host.Id,
	}

	err := p.SendEvent(infosEvent.Wrap())
	if err != nil {
		log.Printf("Room infos could not be sent to %v: %v", p.Username, err.Error())
	}
}

func (p Player) Infos() PlayerInfos {
	return PlayerInfos{
		Username: p.Username,
		Team:     p.team,
		Id:       p.Id,
	}
}

func (p *Player) SetRoom(r *Room) {
	p.room = r
}

func (p *Player) SetWsConn(conn *websocket.Conn) {
	p.wsConn = conn
}

func (p *Player) SetTeam(t *Team) {
	p.team = t.Id
}

func (p Player) SendEvent(e Event) error {
	return p.wsConn.WriteMessage(websocket.TextMessage, e.Marshal())
}
