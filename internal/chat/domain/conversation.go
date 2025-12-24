package domain

import (
	"errors"
	"time"
)

// 1 conversation can have * messages [1 to many ]

type ConversationType string
type ConversationID string
type UserID string

var ErrInvalidUserID = errors.New("invalid user id")

func NewUserID(value string) (UserID, error) {
	if value == "" {
		return "", ErrInvalidUserID
	}
	return UserID(value), nil
}

func NewConversationID(value string) ConversationID {
	return ConversationID(value)
}

const (
	ConversationTypeDirect ConversationType = "direct"
)

type Conversation struct {
	ID           ConversationID
	Type         ConversationType
	Participants []UserID
	CreatedAt    time.Time
}

func (c Conversation) HasParticipant(userID UserID) bool {
	for _, p := range c.Participants {
		if p == userID {
			return true
		}
	}
	return false
}
