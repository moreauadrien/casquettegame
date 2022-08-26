package game

import (
	"fmt"
	"time"
	"timesup/events"
	"timesup/structures"

	"github.com/google/uuid"
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
	Id    string
	host  *Player
	state GameState

	players       []*Player
	teams         [2]Team
	speaker       *Player
	speakerNumber int

	guessedCards     []string
	turnGuessedCards []string
	remainingCards   *structures.CardPile
	roundEndSignal   chan struct{}

	round int
}

func NewRoom(host *Player) *Room {
	r := &Room{
		Id:             uuid.NewString(),
		host:           host,
		players:        []*Player{host},
		state:          WaitingRoom,
		teams:          [2]Team{newTeam(BLUE), newTeam(PURPLE)},
		guessedCards:   make([]string, 0, 40),
		remainingCards: &structures.CardPile{},
		round:          1,
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

func (r *Room) BrodcastEvent(eventType string, eventData events.EventData, except ...*Player) {
	exceptSet := make(map[*Player]bool, len(except))
	for _, p := range except {
		exceptSet[p] = true
	}

	for _, p := range r.players {
		if exceptSet[p] == false {
			p.SendEvent(eventType, eventData, nil)
		}
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

	case Turn:
		r.turnGuessedCards = make([]string, 0, len(r.guessedCards))

		r.BrodcastEvent("stateUpdate", events.StateUpdateData{
			State: string(state),
		}, r.speaker)

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
		r.BrodcastEvent("stateUpdate", events.StateUpdateData{
			State: string(state),
			Cards: r.turnGuessedCards,
		})

	case ScoreRecap:
		r.BrodcastEvent("stateUpdate", events.StateUpdateData{
			State:  string(state),
			Scores: r.Scores(),
		})
	}

}

func (r *Room) StartNextRound() {
	r.round++
	r.remainingCards.Add(Shuffle(r.guessedCards)...)
	r.guessedCards = make([]string, 0, 40)
	r.SetState(WaitTurnStart)
}

func (r *Room) Scores() []events.TeamPoints {
	s := make([]events.TeamPoints, 0, 2)

	for _, t := range r.teams {
		s = append(s, events.TeamPoints{
			Team:   t.Color(),
			Points: t.Points(),
		})
	}

	return s
}

func (r *Room) ValidateCard() {
	c := r.remainingCards.Remove()

	r.guessedCards = append(r.guessedCards, c.Value())
	r.turnGuessedCards = append(r.turnGuessedCards, c.Value())

	r.speaker.team.IncrementRoundScore(r.round)

	r.BrodcastEvent("turnUpdate", events.TurnUpdate{Cards: r.turnGuessedCards}, r.speaker)

	if r.remainingCards.Len() == 0 {
		close(r.roundEndSignal)
	}
}

func (r *Room) PassCard() {
	c := r.remainingCards.Remove()
	r.remainingCards.Add(c.Value())
}

func (r *Room) ChangeSpeaker() {
	r.speakerNumber++
	t := r.teams[r.speakerNumber%2]

	r.speaker = t.GetPlayer(r.speakerNumber / 2)
}

func (r *Room) Start() error {

	if r.state != WaitingRoom {
		return fmt.Errorf("room #%v has already started", r.Id)
	}

	r.SetState(TeamsRecap)
	r.remainingCards.Add(RandomCards(40)...)

	r.speaker = r.host

	go func() {
		time.Sleep(1 * time.Second)
		r.SetState(WaitTurnStart)
	}()

	return nil
}
