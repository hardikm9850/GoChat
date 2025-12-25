package hub

import "github.com/hardikm9850/GoChat/internal/chat/domain"

// MessageEvent represents a message to be broadcast to one or more users.
type MessageEvent struct {
	Message    domain.Message
	Recipients []string // list of userIDs to receive the message
}
