package domain

// 1 conversation can have * messages [1 to many]

import "time"

type MessageID string

type Message struct {
	ID             MessageID
	ConversationID ConversationID
	SenderID       UserID
	Content        string
	CreatedAt      time.Time
}
