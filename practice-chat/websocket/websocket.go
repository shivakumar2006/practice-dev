package Customwebsocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Connection, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket Connection Error: ", err)
		return nil, err
	}

	log.Println("websocket connection established successfully...")
	return connection, nil
}
