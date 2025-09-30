package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

type Client struct {
	url  string
	conn *websocket.Conn
}

// NewClient creates a new client instance
func NewClient(addr string) (*Client, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
	return &Client{
		url: u.String(),
	}, nil
}

// Connect establishes a connection to the WebSocket server.
func (c *Client) Connect() error {
	log.Printf("Connecting to %s", c.url)
	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Close() {
	log.Println("Closing connection.")
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, ""))
	if err != nil {
		log.Println("Write close error:", err)
	}
}

// Run starts the client's read and write loops.
func (c *Client) Run() error {
	// First, establish the connection.
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Close()
	defer c.conn.Close() // Also ensure raw connection is closed.

	// Channel to signal that the read goroutine has finished.
	done := make(chan struct{})

	// Start the read goroutine.
	go func() {
		defer close(done)
		for {
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				// Check for clean and unclean close messages.
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Read error: %v", err)
				}
				return
			}
			log.Printf("Received: %s", message)
			fmt.Print("Enter text: ")
		}
	}()

	// A dedicated goroutine for reading stdin
	userInputChan := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading stdin:", err)
				close(userInputChan)
				return
			}
			userInputChan <- text
		}
	}()

	// Channel for handling OS interrupt signals (Ctrl+C).
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	fmt.Print("Enter text: ")

	for {
		select {
		case <-done:
			// The read loop ended (e.g., server disconnected), so we exit.
			return nil

		case text := <-userInputChan:
			if text == "" { // Channel was closed
				return nil
			}
			err := c.conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				log.Println("Write error:", err)
				return err
			}

		case <-interrupt:
			// User pressed Ctrl+C. The deferred Close() methods will be called.
			return nil
		}
	}
}

func main() {
	// Use the flag package to make the server address configurable.
	addr := flag.String("addr", "localhost:8080", "websocket server address")
	flag.Parse()

	client, err := NewClient(*addr)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Run the client. This will block until the client is done.
	if err := client.Run(); err != nil {
		log.Fatalf("Client run error: %v", err)
	}

	log.Println("Client finished.")
}
