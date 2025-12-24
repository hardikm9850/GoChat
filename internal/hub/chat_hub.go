package hub

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hardikm9850/GoChat/internal/chat/domain"
)

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
	Register   chan ConnEvent                      // channel for new connections
	Unregister chan ConnEvent                      // channel for disconnected connections
	Broadcast  chan MessageEvent                   // channel for messages to Broadcast
	mu         sync.Mutex                          // protects clients and buffer maps
}

// NewHub creates a new instance of Hub.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*websocket.Conn]bool),
		buffer:     make(map[string][]domain.Message),
		Register:   make(chan ConnEvent),
		Unregister: make(chan ConnEvent),
		Broadcast:  make(chan MessageEvent),
	}
}

// SafeWriteJSON sends a message, logs error, but doesn't stop the Hub
func SafeWriteJSON(conn *websocket.Conn, msg interface{}) {
	if err := conn.WriteJSON(msg); err != nil {
		// log and ignore
		fmt.Println("Hub: failed to write message:", err)
	}
}

func (h *Hub) Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Hub panic:", r)
		}
	}()
	log.Println("Hub.Run started")
	for {
		select {

		case event := <-h.Register:
			h.handleRegister(event)

		case event := <-h.Unregister:
			h.handleUnregister(event)

		case event := <-h.Broadcast:
			h.handleBroadcast(event)
		}
	}
}

func (h *Hub) handleRegister(event ConnEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[event.UserID]; !ok {
		h.clients[event.UserID] = make(map[*websocket.Conn]bool)
	}

	h.clients[event.UserID][event.Conn] = true

	// Deliver buffered messages if any
	if buffered, ok := h.buffer[event.UserID]; ok {
		for _, msg := range buffered {
			if err := event.Conn.WriteJSON(msg); err != nil {
				log.Println("register write error:", err)
				break
			}
		}
		delete(h.buffer, event.UserID)
	}
}

func (h *Hub) handleUnregister(event ConnEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns, ok := h.clients[event.UserID]
	if !ok {
		return
	}

	if _, exists := conns[event.Conn]; exists {
		delete(conns, event.Conn)
		_ = event.Conn.Close()
	}

	if len(conns) == 0 {
		delete(h.clients, event.UserID)
	}
}

func (h *Hub) handleBroadcast(event MessageEvent) {
	h.mu.Lock()

	// Copy connections to avoid holding lock while writing
	userConns := make([]*websocket.Conn, 0)

	for _, userID := range event.Recipients {
		if conns, ok := h.clients[userID]; ok && len(conns) > 0 {
			for c := range conns {
				userConns = append(userConns, c)
			}
		} else {
			// Offline user → buffer message
			h.buffer[userID] = append(h.buffer[userID], event.Message)
		}
	}
	log.Printf("handleBroadcast: message ID %s to %d connections", event.Message.ID, len(userConns))

	h.mu.Unlock()

	// Write outside lock
	for _, conn := range userConns {
		if err := conn.WriteJSON(event.Message); err != nil {
			log.Println("broadcast write error:", err)
		}
	}
}
