package game

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var rooms = map[string]*Room{}
var players = map[string]*Player{}

type CreateRoomPayload struct {
	Username string `json:"username" binding:"required,min=4"`
}

type JoinRoomPayload struct {
	Username string `json:"username" binding:"required,min=4"`
	RoomId   string `json:"roomId" binding:"required,uuid4_rfc4122"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type CreateRoomSuccess struct {
	Status string `json:"status"`
	RoomId string `json:"roomId"`
}

type JoinRoomSuccess struct {
	Status string `json:"status"`
}

func CreateRoomHandler(c *gin.Context) {
	var payload CreateRoomPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: "username must contain at least 4 characters",
		})
		return
	}

	p := NewPlayer(payload.Username)
	r := NewRoom(p)

	c.SetCookie("token", p.Token, 2*3600, "/", "", false, true)

	c.JSON(http.StatusOK, CreateRoomSuccess{
		Status: "ok",
		RoomId: r.Id,
	})
}

func JoinRoomHandler(c *gin.Context) {
	var payload JoinRoomPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		var m string
		if len(payload.Username) < 4 {
			m = "username must contain at least 4 characters"
		} else {
			m = "roomId must be uuid4_rfc4122"
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: m,
		})
		return
	}

	r := rooms[payload.RoomId]
	if r == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("room %v does not exist", payload.RoomId),
		})
		return
	}

	if r.state != WaitingRoom {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("room %v has already stated", payload.RoomId),
		})
		return
	}

	if r.usernameSet.Has(payload.Username) {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("the username %v is already taken", payload.Username),
		})
		return
	}

	p := NewPlayer(payload.Username)
	r.AddPlayer(p)

	r.Brodcast(gin.H{"players": r.Players()}, p)

	c.SetCookie("token", p.Token, 2*3600, "/", "", false, true)

	c.JSON(http.StatusOK, JoinRoomSuccess{
		Status: "ok",
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(c *gin.Context) {
	token, _ := c.Cookie("token")
	p := players[token]
	if p == nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	p.Conn = conn

	p.SendFullState()

	for {
		t, msg, err := conn.ReadMessage()

		if err != nil {
			break
		}

		handleIncomingMessage(p, t, msg)
	}
}

func handleIncomingMessage(p *Player, msgType int, msg []byte) {
	if msgType != websocket.TextMessage {
		log.Printf("message type %v is not supported", msgType)
		return
	}

	r := p.Room

	switch string(msg) {
	case "startGame":
		if p.IsHost() == false || r.state != WaitingRoom {
			return
		}

		r.Start()

	case "startTurn":
		if p.IsSpeaker() == false || r.state != WaitTurnStart {
			return
		}

		r.SetState(Turn)

	case "validateCard":
		if p.IsSpeaker() == false || r.state != Turn {
			return
		}

		r.ValidateCard()

	case "passCard":
		if p.IsSpeaker() == false || r.state != Turn {
			return
		}

		r.PassCard()

	case "handOver":
		r.ChangeSpeaker()

		if r.remainingCards.Len() == 0 {
			r.SetState(ScoreRecap)
		} else {
			r.SetState(WaitTurnStart)
		}

	case "startNextRound":
		if p.IsHost() == false || r.state != ScoreRecap {
			return
		}

		r.StartNextRound()

	case "acceptRules":
		if p.IsHost() == false || r.state != RulesRecap {
			return
		}

		r.SetState(WaitTurnStart)

	case "validateCardSwitch":
		if r.state != CardSelection || p.Cards == nil {
			return
		}

		r.remainingCards.Add(p.Cards.Selected...)
		p.Cards = nil

		if r.remainingCards.Len() == 40 {
			r.SetState(RulesRecap)
		} else {
			data := gin.H{"state": "waitPlayers", "cards": []string{}, "swapsRemaining": 0}
			p.SendMessage(data)
		}

	default:
		if strings.HasPrefix(string(msg), "changeCard:") {
			s := strings.Split(string(msg), ":")
			if len(s) < 2 {
				return
			}

			i64, err := strconv.ParseInt(s[1], 10, 32)
			if err != nil {
				return
			}

			i := int(i64)

			if p.Cards != nil && len(p.Cards.Stock) > 0 && len(p.Cards.Selected) > i {
				p.Cards.Selected[i] = p.Cards.Stock[0]
				p.Cards.Stock = p.Cards.Stock[1:]
				p.SendCardsToSelect()
			}
		}

	}

}
