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
			Team:    p.team.Color(),
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

		if r.state != WaitingRoom {
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
			Team:    p.team.Color(),
		}
	})

	wrapper.On("startGame", func(c ws.Client, creds payloads.Credentials, ed events.EventData) *events.ResponseData {
		p := players[creds.Token]
		if p == nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "this player does not exist",
			}
		}

		r := p.Room
		if r == nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is not in a room",
			}
		}

		if r.host.Token != p.Token {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is not the host of his room",
			}
		}

		if err := r.Start(); err != nil {
			return &events.ResponseData{
				Status:  "error",
				Message: err.Error(),
			}
		}

		return &events.ResponseData{
			Status: "success",
		}
	})

	wrapper.On("startTurn", func(c ws.Client, creds payloads.Credentials, ed events.EventData) *events.ResponseData {
		p := players[creds.Token]
		if p == nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "this player does not exist",
			}
		}

		r := p.Room
		if r == nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is not in a room",
			}
		}

		if r.speaker.Token != p.Token {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is currently not the speaker",
			}
		}

		r.SetState(Turn)

		return &events.ResponseData{
			Status: "success",
			Cards:  r.remainingCards,
		}
	})

	wrapper.On("validateCard", func(c ws.Client, creds payloads.Credentials, ed events.EventData) *events.ResponseData {

		p := players[creds.Token]
		if p == nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "this player does not exist",
			}
		}

		r := p.Room
		if r == nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is not in a room",
			}
		}

		if r.speaker.Token != p.Token {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is currently not the speaker",
			}
		}

		if r.state != Turn {
			return &events.ResponseData{
				Status:  "error",
				Message: "the turn is over",
			}
		}

		r.ValidateCard()

		return &events.ResponseData{
			Status: "success",
		}
	})

	wrapper.On("passCard", func(c ws.Client, creds payloads.Credentials, ed events.EventData) *events.ResponseData {

		p := players[creds.Token]
		if p == nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "this player does not exist",
			}
		}

		r := p.Room
		if r == nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is not in a room",
			}
		}

		if r.speaker.Token != p.Token {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is currently not the speaker",
			}
		}

		if r.state != Turn {
			return &events.ResponseData{
				Status:  "error",
				Message: "the turn is over",
			}
		}

		r.PassCard()

		return &events.ResponseData{
			Status: "success",
		}
	})

	wrapper.On("handOver", func(c ws.Client, creds payloads.Credentials, ed events.EventData) *events.ResponseData {
		p := players[creds.Token]
		if p == nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "this player does not exist",
			}
		}

		r := p.Room
		if r == nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is not in a room",
			}
		}

		if r.speaker.Token != p.Token {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is currently not the speaker",
			}
		}

		r.ChangeSpeaker()

		if r.remainingCards.Len() == 0 {
			r.SetState(ScoreRecap)
		} else {
			r.SetState(WaitTurnStart)
		}

		return &events.ResponseData{
			Status: "success",
		}
	})

	wrapper.On("nextRound", func(c ws.Client, creds payloads.Credentials, ed events.EventData) *events.ResponseData {
		p := players[creds.Token]
		if p == nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "this player does not exist",
			}
		}

		r := p.Room
		if r == nil {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is not in a room",
			}
		}

		if r.host.Token != p.Token {
			return &events.ResponseData{
				Status:  "error",
				Message: "player is not the host of his room",
			}
		}

		r.StartNextRound()

		return &events.ResponseData{
			Status: "success",
		}
	})
}
