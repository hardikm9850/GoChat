package database


import "time"

type UserID string
type ConversationID string

type Message struct {
	ID             string         `gorm:"primaryKey"`
	SenderID       UserID         `gorm:"not null"`
	ConversationID ConversationID `gorm:"index;not null"`
	Content        string         `gorm:"type:text;not null"`
	CreatedAt      time.Time
}
