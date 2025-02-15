package main

import (
	"log"
	"net/http"
	"os"
	Customwebsocket "practice-chat/websocket"
)

func serverWs(pool *Customwebsocket.Pool, w http.ResponseWriter, r *http.Request) {
	connection, err := Customwebsocket.Upgrade(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Customwebsocket.Client{
		Connection: connection,
		Pool:       pool,
	}
	pool.Register <- client
	client.Read()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	setUpRoutes()

	log.Printf("Starting the server %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func setUpRoutes() {
	log.Println("This is working...")
	pool := Customwebsocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWs(pool, w, r)
	})
}
