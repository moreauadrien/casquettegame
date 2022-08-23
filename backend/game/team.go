package game

import "timesup/events"

const (
	BLUE   events.TeamColor = "blue"
	PURPLE                  = "purple"
)

type Team struct {
	color   events.TeamColor
	players []*Player
}

func newTeam(color events.TeamColor) Team {
	return Team{
		color:   color,
		players: []*Player{},
	}
}

func (t Team) Len() int {
	return len(t.players)
}

func (t Team) Color() events.TeamColor {
	return t.color
}

func (t *Team) AddPlayer(p *Player) {
	t.players = append(t.players, p)
}
