package game

import "timesup/events"

const (
	BLUE   events.TeamColor = "blue"
	PURPLE                  = "purple"
)

type Team struct {
	color   events.TeamColor
	players []*Player
	points  []int
}

func newTeam(color events.TeamColor) Team {
	return Team{
		color:   color,
		players: []*Player{},
		points:  make([]int, 3, 3),
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

func (t Team) GetPlayer(i int) *Player {
	return t.players[i%t.Len()]
}

func (t Team) Points() []int {
	return t.points
}

func (t *Team) IncrementRoundScore(roundNumber int) {
	t.points[roundNumber-1] += 1
}
