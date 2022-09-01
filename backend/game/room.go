package game

import (
	"math/rand"
	"time"
	"timesup/structures"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GameState string

const (
	WaitingRoom   GameState = "waitingRoom"
	CardSelection GameState = "cardSelection"
	TeamsRecap    GameState = "teamsRecap"
	RulesRecap    GameState = "rulesRecap"
	WaitTurnStart GameState = "waitTurnStart"
	Turn          GameState = "turn"
	TurnRecap     GameState = "turnRecap"
	ScoreRecap    GameState = "scoreRecap"
)

type Room struct {
	Id    string
	host  *Player
	state GameState

	players       []*Player
	teams         [2]Team
	speaker       *Player
	speakerNumber int

	usernameSet structures.Set

	guessedCards     []string
	turnGuessedCards []string
	remainingCards   *structures.CardPile
	roundEndSignal   chan struct{}

	cardPartitions [][]string

	round int
}

func NewRoom(host *Player) *Room {
	r := &Room{
		Id:             uuid.NewString(),
		host:           host,
		players:        []*Player{},
		state:          WaitingRoom,
		teams:          [2]Team{newTeam(BLUE), newTeam(PURPLE)},
		guessedCards:   make([]string, 0, 40),
		remainingCards: &structures.CardPile{},
		round:          1,
		usernameSet:    structures.Set{},
	}

	r.AddPlayer(host)

	rooms[r.Id] = r

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
	p.Room = r
	r.usernameSet.Add(p.Username)
	r.players = append(r.players, p)

	r.addPlayerToSmallestTeam(p)
}

func (r Room) Players() []*PlayerInfos {
	infos := make([]*PlayerInfos, 0, len(r.players))
	for _, p := range r.players {
		infos = append(infos, p.Infos())
	}

	return infos
}

func (r Room) GetFullRoomState() gin.H {
	return gin.H{
		"state":   r.state,
		"host":    r.host.Infos(),
		"speaker": r.speaker.Infos(),
		"players": r.Players(),
		"round":   r.round,
	}
}

func (r *Room) Brodcast(data gin.H, except ...*Player) {
	exceptSet := make(structures.Set, len(except))

	for _, p := range except {
		exceptSet.Add(p.Token)
	}

	for _, p := range r.players {
		if exceptSet.Has(p.Token) == false {
			p.SendMessage(data)
		}
	}
}

func (r *Room) SetState(state GameState) {
	r.state = state

	data := gin.H{"state": state}

	switch state {
	case TeamsRecap:
		for _, p := range r.players {
			p.SendMessage(gin.H{"state": state, "team": p.team.Color()})
		}

		go func() {
			time.Sleep(2 * time.Second)
			r.SetState(CardSelection)
		}()

	case CardSelection:
		for i, p := range r.players {
			cards := r.cardPartitions[i]

			p.Cards = &PlayerCards{Selected: cards[:len(cards)-3], Stock: cards[len(cards)-3:]}
			p.SendCardsToSelect()
		}

	case RulesRecap:
		data["round"] = r.round
		r.Brodcast(data)

	case WaitTurnStart:
		data["round"] = r.round
		data["speaker"] = r.speaker.Infos()
		data["cards"] = []string{}
		r.Brodcast(data)

	case Turn:
		r.turnGuessedCards = make([]string, 0, len(r.guessedCards))
		r.Brodcast(data, r.speaker)
		data["cards"] = r.remainingCards

		r.speaker.SendMessage(data)

		r.roundEndSignal = make(chan struct{})
		go func() {
			select {
			case <-time.After(30 * time.Second):
				r.SetState(TurnRecap)

			case <-r.roundEndSignal:
				r.SetState(TurnRecap)
			}
		}()

	case TurnRecap:
		data["cards"] = r.turnGuessedCards
		r.Brodcast(data)

	case ScoreRecap:
		data["score"] = r.Score()
		r.Brodcast(data)

		if r.round == 3 {
			r.Close()
		}
	}
}

func (r *Room) Score() []TeamPoints {
	s := make([]TeamPoints, 0, 2)

	for _, t := range r.teams {
		s = append(s, TeamPoints{
			Team:   t.Color(),
			Points: t.Points(),
		})
	}

	return s
}

func (r *Room) ChangeSpeaker() {
	r.speakerNumber++
	t := r.teams[r.speakerNumber%2]

	r.speaker = t.GetPlayer(r.speakerNumber / 2)
}

func (r *Room) Start() {
	r.speakerNumber = rand.Intn(len(r.players))
	r.ChangeSpeaker()

	cards := RandomCards(40 + 3*len(r.players))
	r.cardPartitions = partitionCards(cards, len(r.players))

	r.remainingCards = &structures.CardPile{}

	r.SetState(TeamsRecap)
}

func (r *Room) ValidateCard() {
	c := r.remainingCards.Remove()

	r.guessedCards = append(r.guessedCards, c.Value())
	r.turnGuessedCards = append(r.turnGuessedCards, c.Value())

	r.speaker.team.IncrementRoundScore(r.round)

	r.Brodcast(gin.H{"cards": r.turnGuessedCards}, r.speaker)

	if r.remainingCards.Len() == 0 {
		close(r.roundEndSignal)
	}
}

func (r *Room) PassCard() {
	c := r.remainingCards.Remove()
	r.remainingCards.Add(c.Value())
}

func (r *Room) StartNextRound() {
	r.round++
	r.remainingCards.Add(Shuffle(r.guessedCards)...)
	r.guessedCards = make([]string, 0, 40)
	r.SetState(RulesRecap)
}

func (r *Room) Close() {
	for _, p := range r.players {
		delete(players, p.Token)

		p.Conn.WriteMessage(websocket.CloseMessage, []byte("  gameOver  "))
		p.Conn.Close()
	}

	delete(rooms, r.Id)
}
