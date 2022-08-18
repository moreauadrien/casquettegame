package game

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func HandleRoomCreation(c *gin.Context) {
	o := RoomOptions{}
	err := c.BindJSON(&o)

	if err != nil {
		c.String(http.StatusBadRequest, "body is not correct")
		return
	}

	_, exist := Players[o.HostToken]

	if exist {
		c.JSON(http.StatusMethodNotAllowed, ErrorResponse{Code: http.StatusMethodNotAllowed, Message: "player already in a room"})
		return
	}

	p := Player{
		Token:    o.HostToken,
		Id:       o.HostId,
		Username: o.HostUsername,
		room:     nil,
	}

	r := Room{
		Id:            generateRandomId(),
		players:       []*Player{},
		host:          &p,
		numberOfTeams: 2,
		state:         WaitingRoom,
		teams:         [2]*Team{NewTeam(TeamBlue), NewTeam(TeamPurple)},
	}

	Rooms[r.Id] = &r
	Players[p.Token] = &p

	log.Printf("%v create the room %v", p.Username, r.Id)

	c.JSON(http.StatusOK, r)
}

func generateRandomId() string {
	id := ksuid.New()
	return id.String()
}
