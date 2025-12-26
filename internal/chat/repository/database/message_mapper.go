package database

import (
	"github.com/hardikm9850/GoChat/internal/chat/domain"
)

func toDomainMessage(m MessageModel) domain.Message {
	return domain.Message{
		ID:             domain.MessageID(m.ID.String()),
		ConversationID: domain.ConversationID(m.ConversationID.String()),
		SenderID:       domain.UserID(m.SenderID.String()),
		Content:        m.Content,
		CreatedAt:      m.CreatedAt,
	}
}
