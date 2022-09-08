package game

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type PlayerCards struct {
	Selected []string
	Stock    []string
}

type Player struct {
	Token    string
	Username string
	Room     *Room
	Conn     *websocket.Conn
	team     *Team
	Cards    *PlayerCards
}

type PlayerInfos struct {
	Username string    `json:"username"`
	Team     TeamColor `json:"team"`
}

func NewPlayer(username string) *Player {
	p := Player{
		Token:    uuid.NewString(),
		Username: username,
	}

	players[p.Token] = &p

	return &p
}

func (p *Player) SetTeam(t *Team) {
	p.team = t
}

func (p *Player) Infos() *PlayerInfos {
	if p == nil {
		return nil
	}

	return &PlayerInfos{
		Username: p.Username,
		Team:     p.team.Color(),
	}
}

func (p Player) SendMessage(data gin.H) {
	msg, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal data %+v", data)
		return
	}

	if err := p.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("the message \"%v\" could not be sent", string(msg))
	}
}

func (p Player) IsHost() bool {
	return p.Room.host.Token == p.Token
}

func (p Player) IsSpeaker() bool {
	return p.Room.speaker.Token == p.Token
}

func (p Player) SendFullState() {
	data := p.Room.GetFullRoomState()

	switch p.Room.state {
	case CardSelection:
		if p.Cards == nil {
			data["state"] = "waitPlayers"
		} else {
			data["cards"] = p.Cards.Selected
			data["swapsRemaining"] = len(p.Cards.Stock)
		}

	case Turn:
		if p.IsSpeaker() {
			data["cards"] = p.Room.remainingCards
		} else {
			data["cards"] = p.Room.turnGuessedCards
		}

	case TurnRecap:
		data["cards"] = p.Room.turnGuessedCards

	case ScoreRecap:
		data["score"] = p.Room.Score()
	}

	data["team"] = p.team.Color()
	data["username"] = p.Username

	p.SendMessage(data)
}

func (p Player) SendCardsToSelect() {
	data := gin.H{"state": CardSelection, "cards": p.Cards.Selected, "swapsRemaining": len(p.Cards.Stock)}
	p.SendMessage(data)
}
