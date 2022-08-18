package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"timesup/game"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadCards() []string {
	data, err := os.ReadFile("cards.txt")
	check(err)

	return strings.Split(string(data), "\n")
}

var players = map[string]game.Player{}

func main() {

	cards := loadCards()
	fmt.Println(cards)

	router := gin.Default()
	router.Static("/_app", "./static/_app")
	router.StaticFile("/", "./static/app.html")
	router.StaticFile("/favicon.png", "./static/favicon.png")
	router.StaticFile("/manifest.json", "./static/manifest.json")

	router.NoRoute(func(c *gin.Context) {
		c.File("./static/app.html")
	})

	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})

	router.POST("/api/createRoom", game.HandleRoomCreation)

	router.Run(":8080")
}
