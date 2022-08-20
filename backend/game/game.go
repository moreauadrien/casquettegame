package game

import (
	"fmt"
	"timesup/events"
	"timesup/payloads"
	"timesup/ws"
)

var rooms = map[string]*Room{}
var players = map[string]*Player{}

func InitWsHandlers(wrapper ws.Wrapper) {
	wrapper.On("create", func(c ws.Client, cred payloads.Credentials, ed events.EventData) *events.ResponseData {
		p := players[cred.Token]
		if p == nil {
			p = &Player{
				Id:       cred.Id,
				Token:    cred.Token,
				Username: cred.Username,
				Client:   &c,
			}
		}

		players[p.Token] = p

		r := p.Room
		if r != nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is already in a room",
			}
		}

		r = NewRoom(p)
		rooms[r.Id] = r

		p.Room = r

		return &events.ResponseData{
			Status:  "success",
			RoomId:  r.Id,
			Players: r.ListPlayers(),
		}
	})

	wrapper.On("join", func(c ws.Client, creds payloads.Credentials, ed events.EventData) *events.ResponseData {
		data, ok := ed.(*events.JoinData)
		if !ok {
			return &events.ResponseData{
				Status:  "error",
				Message: "event data is not of type *JoinData",
			}
		}

		p := players[creds.Token]
		if p == nil {
			p = &Player{
				Id:       creds.Id,
				Token:    creds.Token,
				Username: creds.Username,
				Client:   &c,
			}
		}

		players[p.Token] = p

		r := p.Room
		if r != nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is already in a room",
			}
		}

		r = rooms[data.RoomId]
		if r == nil {
			return &events.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("room #%v does not exist", data.RoomId),
			}
		}

		if r.state != NotStarted {
			return &events.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("room #%v has already started", r.Id),
			}
		}

		r.AddPlayer(p)
		p.Room = r

		return &events.ResponseData{
			Status:  "success",
			Host:    r.host.Id,
			Players: r.ListPlayers(),
		}
	})
}
