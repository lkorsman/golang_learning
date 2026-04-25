package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"
)

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	Type     string `json:"type"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/chat-client/main.go <username>")
		os.Exit(1)
	}

	username := os.Args[1]

	// Connect to WebSocket
	url := fmt.Sprintf("ws://localhost:8080/ws?username=%s", username)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	fmt.Printf("Connected as %s. Type messages and press Enter. \n", username)
	fmt.Println("Type 'quit' to exit.")
	fmt.Println("---")

	// Channel for interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Goroutine to read messages from server
	go func() {
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				return
			}

			switch msg.Type {
			case "join":
				fmt.Printf("*** %s %s\n", msg.Username, msg.Content)
			case "leave":
				fmt.Printf("*** %s %s\n", msg.Username, msg.Content)
			case "message":
				fmt.Printf("%s: %s\n", msg.Username, msg.Content)
			}
		}
	}()

	// Read from stdin and send to server
	scanner := bufio.NewScanner(os.Stdin)
	done := make(chan struct{})

	go func() {
		for scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())

			if text == "quit" {
				close(done)
				return
			}

			if text == "" {
				continue
			}

			msg := Message{
				Content: text,
			}

			err := conn.WriteJSON(msg)
			if err != nil {
				log.Println("Write error:", err)
				return
			}
		}
	}()

	select {
	case <-done:
		fmt.Println("\nGoodbye!")
	case <-interrupt:
		fmt.Println("\nInterrupted!")
	}
}
