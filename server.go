package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	name string
}

var clients = make(map[string]*Client)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.HandleFunc("/ws", handleConnections)

	fmt.Println("Server running on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	defer ws.Close()

	_, usernameBytes, err := ws.ReadMessage()

	if err != nil {
		return
	}

	username := string(usernameBytes)

	client := &Client{
		conn: ws,
		name: username,
	}

	clients[username] = client

	fmt.Println(username, "connected")

	for {

		_, msgBytes, err := ws.ReadMessage()

		if err != nil {
			delete(clients, username)
			fmt.Println(username, "disconnected")
			break
		}

		msg := string(msgBytes)

		handleMessage(username, msg)
	}
}

func handleMessage(sender string, msg string) {

	// format:
	// target|message

	var target string
	var message string

	for i := 0; i < len(msg); i++ {

		if msg[i] == '|' {
			target = msg[:i]
			message = msg[i+1:]
			break
		}
	}

	client, ok := clients[target]

	if ok {

		finalMsg := fmt.Sprintf("%s: %s", sender, message)

		client.conn.WriteMessage(
			websocket.TextMessage,
			[]byte(finalMsg),
		)
	}
}
