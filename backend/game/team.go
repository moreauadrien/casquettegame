package game

type TeamColor string

const (
	BLUE   TeamColor = "blue"
	PURPLE           = "purple"
)

type TeamPoints struct {
	Team   TeamColor `json:"team"`
	Points []int     `json:"points"`
}

type Team struct {
	color   TeamColor
	players []*Player
	points  []int
}

func newTeam(color TeamColor) Team {
	return Team{
		color:   color,
		players: []*Player{},
		points:  make([]int, 3, 3),
	}
}

func (t Team) Len() int {
	return len(t.players)
}

func (t Team) Color() TeamColor {
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
