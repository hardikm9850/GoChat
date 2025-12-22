package hub

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hardikm9850/GoChat/internal/chat/domain"
)

// ConnEvent represents a single WebSocket connection for a user.
type ConnEvent struct {
	UserID string
	Conn   *websocket.Conn
}

// MessageEvent represents a message to be broadcast to one or more users.
type MessageEvent struct {
	Message     domain.Message
	Recipients  []string // list of userIDs to receive the message
}

// Hub manages all active WebSocket connections and routes messages.
// Responsibilities:
// 1. Track active connections per user (multiple connections per user allowed)
// 2. Broadcast messages to all active connections of recipients
// 3. Buffer messages for offline users
// 4. Handle register/unregister events for connections
// 5. Ensure concurrency safety using channels + mutex
type Hub struct {
	clients    map[string]map[*websocket.Conn]bool // userID -> set of connections
	buffer     map[string][]domain.Message         // userID -> buffered messages for offline delivery
	register   chan ConnEvent                       // channel for new connections
	unregister chan ConnEvent                       // channel for disconnected connections
	broadcast  chan MessageEvent                    // channel for messages to broadcast
	mu         sync.Mutex                           // protects clients and buffer maps
}

// NewHub creates a new instance of Hub.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*websocket.Conn]bool),
		buffer:     make(map[string][]domain.Message),
		register:   make(chan ConnEvent),
		unregister: make(chan ConnEvent),
		broadcast:  make(chan MessageEvent),
	}
}

// Run starts the Hub's main loop.
// It continuously listens on the register, unregister, and broadcast channels
// and updates the internal state accordingly.
func (h *Hub) Run() {
	for {
		select {
		// Register a new WebSocket connection
		case connEvent := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[connEvent.UserID]; !ok {
				h.clients[connEvent.UserID] = make(map[*websocket.Conn]bool)
			}
			h.clients[connEvent.UserID][connEvent.Conn] = true

			// If the user has buffered messages, send them now
			if buffered, exists := h.buffer[connEvent.UserID]; exists {
				for _, msg := range buffered {
					err := connEvent.Conn.WriteJSON(msg)
					if err != nil {
						return
					}
				}
				delete(h.buffer, connEvent.UserID) // clear buffer after sending
			}
			h.mu.Unlock()

		// Unregister a disconnected WebSocket connection
		case connEvent := <-h.unregister:
			h.mu.Lock()
			if conns, ok := h.clients[connEvent.UserID]; ok {
				if _, exists := conns[connEvent.Conn]; exists {
					delete(conns, connEvent.Conn)
					err := connEvent.Conn.Close()
					if err != nil {
						return
					} // close the WebSocket
				}
				// Remove the user from clients map if no connections left
				if len(conns) == 0 {
					delete(h.clients, connEvent.UserID)
				}
			}
			h.mu.Unlock()

		// Broadcast a message to recipients
		case msgEvent := <-h.broadcast:
			h.mu.Lock()
			for _, userID := range msgEvent.Recipients {
				if conns, ok := h.clients[userID]; ok && len(conns) > 0 {
					// Send message to all active connections for this user
					for c := range conns {
						err := c.WriteJSON(msgEvent.Message)
						if err != nil {
							return
						}
					}
				} else {
					// User is offline, buffer the message for later delivery
					h.buffer[userID] = append(h.buffer[userID], msgEvent.Message)
				}
			}
			h.mu.Unlock()
		}
	}
}
