package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"timesup/events"
	"timesup/ws"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
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

func main() {
	wrt := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "./logs/timesup.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     29,
		Compress:   false,
	})
	log.SetOutput(wrt)

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

	router.Run(":8080")
}
