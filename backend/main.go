package main

import (
	"io"
	"log"
	"os"
	"timesup/game"

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

	gin.DefaultWriter = wrt

	router := gin.Default()
	router.Static("/_app", "./static/_app")
	router.StaticFile("/", "./static/app.html")
	router.StaticFile("/favicon.png", "./static/favicon.png")
	router.StaticFile("/manifest.json", "./static/manifest.json")

	router.NoRoute(func(c *gin.Context) {
		c.File("./static/app.html")
	})

	router.GET("/ws", game.WsHandler)

	api := router.Group("/api")
	{
		api.POST("createRoom", game.CreateRoomHandler)
		api.POST("joinRoom", game.JoinRoomHandler)
	}

	router.Run(":8080")
}
