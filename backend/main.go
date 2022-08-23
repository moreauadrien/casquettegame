package main

import (
	"io"
	"log"
	"os"
	"timesup/game"
	"timesup/ws"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	wrt := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "./logs/timesup.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     29,
		Compress:   false,
	})
	log.SetOutput(wrt)

	router := gin.Default()
	router.Static("/_app", "./static/_app")
	router.StaticFile("/", "./static/app.html")
	router.StaticFile("/favicon.png", "./static/favicon.png")
	router.StaticFile("/manifest.json", "./static/manifest.json")

	router.NoRoute(func(c *gin.Context) {
		c.File("./static/app.html")
	})

	wrapper := ws.NewWrapper()

	router.GET("/ws", gin.WrapF(wrapper.HttpHandler))

	game.InitWsHandlers(wrapper)

	router.Run(":8080")
}
