package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = "localhost:8080"
var u = url.URL{Scheme: "ws", Host: addr, Path: "/echo"}

func main() {
	// Create a channel to receive OS interrupt signals (e.g., Ctrl+C).
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error during dial:", err)
	}
	defer conn.Close()

	// Create a channel that will be closed when the read goroutine ends.
	// We use struct{} as the channel's type because we only care about the signal,
	// not the value being sent.
	done := make(chan struct{})

	// Start a goroutine to continuously read messages from the server.
	go func() {
		defer close(done)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error during msg reading:", err)
				return
			}
			log.Printf("Received: %s", msg)
		}
	}()

	// A channel to receive user input from the keyboard.
	userInputChan := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter text: ")
			text, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error during reading stdin:", err)
				close(userInputChan)
				return
			}
			userInputChan <- text
		}
	}()

	// Create a ticker that triggers every second.
	// We can use it for heartbeats or other periodic tasks
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			// If the 'done' channel is closed, it means the read goroutine has finished.
			// This could be due to a server disconnect. We exit the program.
			log.Println("Reader goroutine finished. Exiting.")
			return

		case text := <-userInputChan:
			// Every second, the ticker sends a value. We use this to send a message.
			// err := conn.WriteMessage(websocket.TextMessage, []byte(t.String()))

			// Send the user's message to the WebSocket server
			err := conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				log.Println("Error during msg writing:", err)
				return
			}

		case <-interrupt:
			log.Println("Interrupt signal received. Closing connection.")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during writing close:", err)
				return
			}

			// Wait for the server to close the connection, with a timeout.
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
