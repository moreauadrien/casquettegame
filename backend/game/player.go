package game

import (
	"log"
	"timesup/events"
	"timesup/ws"
)

type Player struct {
	Id       string
	Token    string
	Username string
	Room     *Room
	Client   *ws.Client
}

func (p *Player) GetInfos() events.PlayerInfos {
	return events.PlayerInfos{
		Username: p.Username,
		Id:       p.Id,
	}
}

func (p *Player) SendEvent(eventType string, eventData events.EventData, handler ws.ResponseHandler) {
	err := p.Client.SendEvent(eventType, eventData, handler)

	if err != nil {
		log.Printf("the event could not be sent to %v: %v", p.Username, err.Error())
	}
}
