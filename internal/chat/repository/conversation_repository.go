package repository

import (
	"github.com/hardikm9850/GoChat/internal/chat/domain"
)

type ConversationRepository interface {
	Create(conversation domain.Conversation) error

	FindByID(
		conversationID domain.ConversationID,
		userID domain.UserID,
	) (domain.Conversation, error)

	Find(
		userA, userB domain.UserID,
	) (domain.Conversation, error)

	FindMyConversation(
		userA domain.UserID,
	) ([]domain.Conversation, error)

	Participants(conversationID domain.ConversationID) ([]domain.UserID, error)
}
