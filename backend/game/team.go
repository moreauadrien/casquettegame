package game

const (
	TeamBlue = iota
	TeamPurple
	TeamRed
	TeamYellow
)

type Team struct {
	Id      int
	players []*Player
}

func NewTeam(id int) *Team {
	return &Team{Id: id, players: []*Player{}}
}

func (t *Team) AddPlayer(p *Player) {
	t.players = append(t.players, p)
}

func (t Team) Len() int {
	return len(t.players)
}
