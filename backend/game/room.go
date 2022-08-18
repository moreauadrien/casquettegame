package game

import (
	"fmt"
	"log"
)

const (
	WaitingRoom = iota
	InGame
)

const ALREADY_LAUNCHED_ERROR = "the game has already been launched"

var Rooms = make(map[string]*Room)

type Room struct {
	Id            string `json:"roomId"`
	players       []*Player
	host          *Player
	numberOfTeams int
	state         int
	teams         [2]*Team
}

func (r *Room) getSmallestTeam() *Team {
	if r.teams[0].Len() <= r.teams[1].Len() {
		return r.teams[0]
	} else {
		return r.teams[1]
	}
}

func (r *Room) Players() []PlayerInfos {
	s := make([]PlayerInfos, 0, len(r.players))

	for _, p := range r.players {
		s = append(s, p.Infos())
	}

	return s
}

func (r *Room) AddPlayer(p *Player) {
	r.players = append(r.players, p)

	team := r.getSmallestTeam()
	team.AddPlayer(p)

	p.SetTeam(team)

	joinEvent := JoinEvent{
		Players: r.Players(),
		Host:    r.host.Id,
	}

	r.BrodcastEvent(joinEvent.Wrap())
}

func (r Room) Host() *Player {
	return r.host
}

func (r Room) State() int {
	return r.state
}

func (r *Room) StartGame() error {
	if r.state != WaitingRoom {
		return fmt.Errorf(ALREADY_LAUNCHED_ERROR)
	}

	r.state = InGame

	return nil
}

func (r *Room) BrodcastEvent(e Event) {
	for _, p := range r.players {
		err := p.SendEvent(e)

		if err != nil {
			log.Printf("the event could not be sent to %v: %v", p.Username, err.Error())
		}
	}
}

type RoomOptions struct {
	NumberOfTeams int    `json:"numberOfTeams" binding:"required"`
	HostUsername  string `json:"hostUsername" binding:"required"`
	HostToken     string `json:"hostToken" binding:"required"`
	HostId        string `json:"hostId" binding:"required"`
}
