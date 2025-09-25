package main

import (
	"fmt"
	"log"
	"net/http"
	
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	// Upgrade upgrades the HTTP server connection to the WebSocket protocol
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during connection upgrading:", err)
		return
	}
	defer conn.Close()
    for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during msg reading:", err)
			break
		}
		log.Printf("Received: %s", msg)
		err = conn.WriteMessage(msgType, msg)
		if err != nil {
			log.Println("Error during msg writing:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Home page")
}

func main() {
    http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8080", nil))
}