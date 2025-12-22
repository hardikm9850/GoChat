package repository

import (
	"time"

	"github.com/hardikm9850/GoChat/internal/chat/domain"
)

type MessageOrder int

const (
	OrderAsc  MessageOrder = iota // oldest → newest
	OrderDesc                     // newest → oldest
)

type MessageRepository interface {
	// Save persists a message in a conversation.
	// Returns ErrConversationNotFound if the conversation does not exist.
	Save(message domain.Message) error

	// Find retrieves messages for a conversation.
	//
	// - Results are ordered by CreatedAt.
	// - Pagination is cursor-based using time.
	// - If cursor is nil:
	//     - OrderDesc → latest messages
	//     - OrderAsc  → oldest messages
	//
	// Returns:
	// - ErrConversationNotFound if conversation does not exist
	// - empty slice if no messages exist
	Find(
		conversationID domain.ConversationID,
		limit int,
		after *time.Time,
		order MessageOrder,
	) ([]domain.Message, error)
}
