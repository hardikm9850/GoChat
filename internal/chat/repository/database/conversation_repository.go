package database

import (
	"errors"

	domain "github.com/hardikm9850/GoChat/internal/chat/domain"
	"github.com/hardikm9850/GoChat/internal/chat/repository"
	"gorm.io/gorm"
)

type ConversationRepositoryMySQL struct {
	db *gorm.DB
}

func New(db *gorm.DB) repository.ConversationRepository {
	return &ConversationRepositoryMySQL{db: db}
}

func (r *ConversationRepositoryMySQL) Create(
	conversation domain.Conversation,
) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		conv := Conversation{
			ID:    string(conversation.ID),
			Title: "",
		}

		if err := tx.Create(&conv).Error; err != nil {
			return err
		}

		for _, userID := range conversation.Participants {
			cp := ConversationParticipant{
				ConversationID: conv.ID,
				UserID:         string(userID),
			}
			if err := tx.Create(&cp).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *ConversationRepositoryMySQL) FindByID(
	conversationID domain.ConversationID,
	userID domain.UserID,
) (domain.Conversation, error) {

	var count int64
	err := r.db.Table("conversation_participants").
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Count(&count).Error

	if err != nil {
		return domain.Conversation{}, err
	}

	if count == 0 {
		return domain.Conversation{}, repository.ErrConversationNotFound
	}

	var conv Conversation
	if err := r.db.
		Preload("Participants").
		First(&conv, "id = ?", conversationID).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Conversation{}, repository.ErrConversationNotFound
		}
		return domain.Conversation{}, err
	}

	return mapToDomain(conv), nil
}

func (r *ConversationRepositoryMySQL) Find(
	userA, userB domain.UserID,
) (domain.Conversation, error) {

	var conv Conversation

	err := r.db.Raw(`
		SELECT c.*
		FROM conversations c
		JOIN conversation_participants p1 ON p1.conversation_id = c.id
		JOIN conversation_participants p2 ON p2.conversation_id = c.id
		WHERE p1.user_id = ? AND p2.user_id = ?
		LIMIT 1
	`, userA, userB).Scan(&conv).Error

	if err != nil || conv.ID == "" {
		return domain.Conversation{}, repository.ErrConversationNotFound
	}

	var participants []ConversationParticipant
	r.db.Where("conversation_id = ?", conv.ID).Find(&participants)
	conv.Participants = participants

	return mapToDomain(conv), nil
}

func (r *ConversationRepositoryMySQL) FindMyConversation(
	userID domain.UserID,
) ([]domain.Conversation, error) {

	var convs []Conversation

	err := r.db.Raw(`
		SELECT DISTINCT c.*
		FROM conversations c
		JOIN conversation_participants cp ON cp.conversation_id = c.id
		WHERE cp.user_id = ?
		ORDER BY c.updated_at DESC
	`, userID).Scan(&convs).Error

	if err != nil {
		return nil, err
	}

	var result []domain.Conversation
	for _, conv := range convs {
		var participants []ConversationParticipant
		r.db.Where("conversation_id = ?", conv.ID).Find(&participants)
		conv.Participants = participants
		result = append(result, mapToDomain(conv))
	}

	return result, nil
}

func (r *ConversationRepositoryMySQL) Participants(
	conversationID domain.ConversationID,
) ([]domain.UserID, error) {

	var rows []ConversationParticipant

	if err := r.db.
		Where("conversation_id = ?", conversationID).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, repository.ErrConversationNotFound
	}

	var users []domain.UserID
	for _, r := range rows {
		users = append(users, domain.UserID(r.UserID))
	}

	return users, nil
}


func mapToDomain(c Conversation) domain.Conversation {
	var participants []domain.UserID
	for _, p := range c.Participants {
		participants = append(participants, domain.UserID(p.UserID))
	}

	return domain.Conversation{
		ID:           domain.ConversationID(c.ID),
		Participants: participants,
		CreatedAt:    c.CreatedAt,
	}
}
