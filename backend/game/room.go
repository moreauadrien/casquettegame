package game

import (
	"fmt"
	"time"
	"timesup/events"

	_ "github.com/google/uuid"
)

type GameState string

const (
	WaitingRoom   GameState = "waitingRoom"
	TeamsRecap    GameState = "teamsRecap"
	RulesRecap    GameState = "rulesRecap"
	WaitTurnStart GameState = "waitTurnStart"
	Turn          GameState = "turn"
	TurnRecap     GameState = "turnRecap"
	ScoreRecap    GameState = "scoreRecap"
)

type Room struct {
	Id      string
	players []*Player
	host    *Player
	state   GameState
	teams   [2]Team
	speaker *Player
}

func NewRoom(host *Player) *Room {
	r := &Room{
		//Id:      uuid.NewString(),
		Id:      "abcdefg",
		host:    host,
		players: []*Player{host},
		state:   WaitingRoom,
		teams:   [2]Team{newTeam(BLUE), newTeam(PURPLE)},
	}

	r.addPlayerToSmallestTeam(host)

	return r
}

func (r *Room) addPlayerToSmallestTeam(p *Player) {
	if r.teams[0].Len() <= r.teams[1].Len() {
		r.teams[0].AddPlayer(p)
		p.SetTeam(&r.teams[0])
	} else {
		r.teams[1].AddPlayer(p)
		p.SetTeam(&r.teams[1])
	}
}

func (r *Room) AddPlayer(p *Player) {
	r.addPlayerToSmallestTeam(p)

	playerList := append(r.ListPlayers(), *p.GetInfos())
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
		list = append(list, *p.GetInfos())
	}

	return list
}

func (r *Room) SetState(state GameState) {
	r.state = state
	switch state {
	case TeamsRecap:
		r.BrodcastEvent("stateUpdate", events.StateUpdateData{
			State:   string(state),
			Players: r.ListPlayers(),
		})

	case WaitTurnStart:
		r.BrodcastEvent("stateUpdate", events.StateUpdateData{
			State:   string(state),
			Speaker: r.speaker.GetInfos(),
		})
	}
}

func (r *Room) Start() error {
	if r.state != WaitingRoom {
		return fmt.Errorf("room #%v has already started", r.Id)
	}

	r.SetState(TeamsRecap)

	r.speaker = r.host

	go func() {
		time.Sleep(5 * time.Second)
		r.SetState(WaitTurnStart)
	}()

	return nil
}
