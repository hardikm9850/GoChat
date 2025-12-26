package database

import (
	"errors"
	"github.com/google/uuid"
	"github.com/hardikm9850/GoChat/internal/chat/domain"
	"github.com/hardikm9850/GoChat/internal/chat/repository"
	"gorm.io/gorm"
	"time"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepo(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Save(
	senderID domain.UserID,
	conversationID domain.ConversationID,
	content string,
) (domain.Message, error) {

	convID, err := uuid.Parse(string(conversationID))
	if err != nil {
		return domain.Message{}, err // or panic if impossible
	}

	senderUUID, err := uuid.Parse(string(senderID))
	if err != nil {
		return domain.Message{}, err
	}

	msg := MessageModel{
		ID:             uuid.New(),
		ConversationID: convID,
		SenderID:       senderUUID,
		Content:        content,
		CreatedAt:      time.Now().UTC(),
	}

	if err := r.db.Create(&msg).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return domain.Message{}, repository.ErrConversationNotFound
		}
		return domain.Message{}, err
	}
	return toDomainMessage(msg), nil
}

func (r *MessageRepository) Find(
	conversationID domain.ConversationID,
	limit int,
	after *time.Time,
	order repository.MessageOrder,
) ([]domain.Message, error) {

	var models []MessageModel

	q := r.db.
		Where("conversation_id = ?", conversationID)

	if after != nil {
		if order == repository.OrderAsc {
			q = q.Where("created_at > ?", *after)
		} else {
			q = q.Where("created_at < ?", *after)
		}
	}

	if order == repository.OrderAsc {
		q = q.Order("created_at ASC")
	} else {
		q = q.Order("created_at DESC")
	}

	if err := q.Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}

	messages := make([]domain.Message, 0, len(models))
	for _, m := range models {
		message := toDomainMessage(m)
		messages = append(messages, message)
	}

	return messages, nil
}
