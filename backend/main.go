package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"timesup/game"
	"timesup/ws"

	"github.com/gin-gonic/gin"
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
	logFileName := fmt.Sprintf("%s.%s", "./logs/", "%Y-%m-%d.%H:%M:%S")
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	log.Println(" Orders API Called")

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

	router.GET("/ws", gin.WrapF(ws.Wrapper.HttpHandler))

	ws.Wrapper.On("test", func(c ws.Client, p *ws.Payload) {
		fmt.Println("test")

		c.SendEvent(game.Event{Type: "myevent"}, nil)
	})

	router.POST("/api/createRoom", game.HandleRoomCreation)

	router.Run(":8080")
}
