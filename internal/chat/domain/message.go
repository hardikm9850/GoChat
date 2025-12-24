package domain

// 1 conversation can have * messages [1 to many]

import (
	"github.com/google/uuid"
	"time"
)

type MessageID string

type Message struct {
	ID             MessageID
	ConversationID ConversationID
	SenderID       UserID
	Content        string
	CreatedAt      time.Time
}

func NewMessageID() MessageID {
	return MessageID(uuid.NewString())
}
