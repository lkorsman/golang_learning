package chat

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

// ServerWS handles websocket requests from clients
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	// Get username from query param
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "Anonymous"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade connection: %v\n", err)
		return
	}

	client := &Client{
		hub: 		h.hub,
		conn:		conn,
		send: 		make(chan *Message, 256),
		username: 	username,
	}

	client.hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// Stats returns chat statistics
func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "applicaiton/json")
	fmt.Fprintf(w, `{"connected_users": %d}`, h.hub.GetClientCount())
}