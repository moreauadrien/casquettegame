package game

import (
	"timesup/events"

	_ "github.com/google/uuid"
)

type GameState int

const (
	NotStarted GameState = iota
)

type Room struct {
	Id      string
	players []*Player
	host    *Player
	state   GameState
}

func NewRoom(host *Player) *Room {
	return &Room{
		//Id:      uuid.NewString(),
		Id:      "abcdefg",
		host:    host,
		players: []*Player{host},
		state:   NotStarted,
	}
}

func (r *Room) AddPlayer(p *Player) {
	playerList := append(r.ListPlayers(), p.GetInfos())
	r.BrodcastEvent("playerJoin", struct {
		Players []events.PlayerInfos `json:"players"`
	}{Players: playerList})

	r.players = append(r.players, p)
}

func (r *Room) BrodcastEvent(eventType string, eventData events.EventData) {
	for _, p := range r.players {
		p.SendEvent(eventType, eventData, nil)
	}
}

func (r *Room) ListPlayers() []events.PlayerInfos {
	list := make([]events.PlayerInfos, 0, len(r.players))

	for _, p := range r.players {
		list = append(list, p.GetInfos())
	}

	return list
}
