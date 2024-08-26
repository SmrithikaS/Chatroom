package main

import (
	"log"

	"chatServer.go/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	api := router.Group("/chatroom")

	api.GET("/:id/ws", handler.HandleWebSocket)

	go handler.BroadcastMessages()

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
