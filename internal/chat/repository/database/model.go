package database

import "time"

type Conversation struct {
	ID        string `gorm:"type:char(36);primaryKey"`
	Title     string `gorm:"varchar(100)"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Participants []ConversationParticipant `gorm:"foreignKey:ConversationID"`
}

type ConversationParticipant struct {
	ConversationID string    `gorm:"type:char(36);primaryKey"`
	UserID         string    `gorm:"type:char(36);primaryKey"`
	JoinedAt       time.Time `gorm:"autoCreateTime"`
}
