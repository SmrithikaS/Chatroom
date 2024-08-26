package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	ChatRoomid string `json:"chatid"`
}

var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(create *gin.Context) {
	chatroomid := create.Param("id")
	connect, err := upgrader.Upgrade(create.Writer, create.Request, nil)
	if err != nil {
		http.NotFound(create.Writer, create.Request)
		return
	}
	defer connect.Close()
	for {
		var msg Message
		err := connect.ReadJSON(&msg)
		if err != nil {
			return
		}
		var Counter int
		Counter++
		msg.ID = strconv.Itoa(Counter)
		msg.ChatRoomid = chatroomid
		clients[connect] = chatroomid
		broadcast <- msg
	}

}

func BroadcastMessages() {
	for {
		messages := <-broadcast
		for client, chatroomid := range clients {
			if messages.ChatRoomid == chatroomid {
				err := client.WriteJSON(messages)
				if err != nil {
					delete(clients, client)
				}
			}
		}
	}
}
