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
	team     *Team
}

func (p *Player) GetInfos() events.PlayerInfos {
	return events.PlayerInfos{
		Username: p.Username,
		Id:       p.Id,
		Team:     p.team.Color(),
	}
}

func (p *Player) SetTeam(t *Team) {
	p.team = t
}

func (p *Player) SendEvent(eventType string, eventData events.EventData, handler ws.ResponseHandler) {
	err := p.Client.SendEvent(eventType, eventData, handler)

	if err != nil {
		log.Printf("the event could not be sent to %v: %v", p.Username, err.Error())
	}
}
