package chat

import (
	"fmt"
	"sync"
)

type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan *Message

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	mu sync.RWMutex
}

type Message struct {
	Username	string `json:"username"`
	Content		string `json:"content"`
	Type 		string `json:"type"` 	// "message", "join", "leave"
}

func NewHub() *Hub {
	return &Hub{
		clients:	make(map[*Client]bool),
		broadcast: 	make(chan *Message, 256),
		register: 	make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			
			fmt.Printf("Client registered: %s (total: %d)\n", client.username, len(h.clients))

			// Broadcast join message
			h.broadcast <- &Message{
				Username: 	client.username,
				Content: 	"joined the chat",
				Type: 		"join",
			}

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

			fmt.Printf("Client unregistered: %s (total: %d)\n", client.username, len(h.clients))
            
            // Broadcast leave message
            h.broadcast <- &Message{
                Username: client.username,
                Content:  "left the chat",
                Type:     "leave",
            }

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client is slow, unregister them
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.Unlock()
	return len(h.clients)
}