package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

// Test function must start with "Test" and take a *testing.T argument.
func TestEcho(t *testing.T) {
	// 1. Create a new test server. 
	// httptest.NewServer will start a server on a random available port.
	server := httptest.NewServer(http.HandlerFunc(echo))
	defer server.Close() // Ensure the server is closed after the test finishes.

	// 2. Convert the server's http:// URL to a ws:// URL.
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// 3. Act as a client and connect to the test server.
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to dial WebSocket server: %v", err)
	}
	defer conn.Close()

	// 4. Define the message we want to send.
	testMessage := "hello world"

	// 5. Write the message to the WebSocket connection.
	if err := conn.WriteMessage(websocket.TextMessage, []byte(testMessage)); err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}

	// 6. Read the response from the server.
	_, receivedMessage, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	// 7. Assert that the received message is the same as the one we sent.
	if string(receivedMessage) != testMessage {
		t.Errorf("Received message is not correct. Got: %s, Want: %s", receivedMessage, testMessage)
	}
}