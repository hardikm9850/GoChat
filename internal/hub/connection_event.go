package hub

import "github.com/gorilla/websocket"

// ConnEvent represents a single WebSocket connection for a user.
type ConnEvent struct {
	UserID string
	Conn   *websocket.Conn
}
