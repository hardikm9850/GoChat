package database

import (
	"time"

	"github.com/google/uuid"
)

type MessageModel struct {
	ID             uuid.UUID `gorm:"type:char(36);primaryKey"`
	ConversationID uuid.UUID `gorm:"type:char(36);index;not null"`
	SenderID       uuid.UUID `gorm:"type:char(36);not null"`
	Content        string    `gorm:"type:text;not null"`
	CreatedAt      time.Time `gorm:"index"`
}
