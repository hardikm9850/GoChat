package usecase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/hardikm9850/GoChat/internal/chat/domain"
	"github.com/hardikm9850/GoChat/internal/chat/repository"
	"time"
)

type ConversationResponse struct {
	ConversationID string   `json:"conversation_id"`
	Participants   []string `json:"participants"`
	CreatedAt      string   `json:"created_at"`
}

type CreateConversation struct {
	clock func() time.Time
	repo  repository.ConversationRepository
}

func New(repo repository.ConversationRepository) CreateConversation {
	return CreateConversation{repo: repo, clock: time.Now}
}

func (uc *CreateConversation) CreateConversation(
	initiatorID, participantID domain.UserID) (domain.Conversation, error) {

	if initiatorID == participantID {
		return domain.Conversation{}, domain.ErrSameParticipant
	}
	initiatorID, participantID = NormalizeParticipants(initiatorID, participantID)

	existingConversation, err := uc.repo.Find(initiatorID, participantID)
	if err == nil {
		// conversation already exist
		return existingConversation, nil
	}
	if !errors.Is(err, domain.ErrConversationNotFound) {
		return domain.Conversation{}, err
	}

	// conversation doesn't exist between this 2 users so we create a one.

	conversation := domain.Conversation{
		ID:           domain.NewConversationID(uuid.NewString()),
		Participants: []domain.UserID{initiatorID, participantID},
		Type:         domain.ConversationTypeDirect,
		CreatedAt:    uc.clock(),
	}

	err = uc.repo.Create(conversation)

	if errors.Is(err, domain.ErrConversationExists) {
		// Race condition protection
		return uc.repo.Find(initiatorID, participantID)
	}
	if err != nil {
		return domain.Conversation{}, err
	}

	return conversation, nil
}

func (uc *CreateConversation) FindConversation(
	conversationId domain.ConversationID,
	userId domain.UserID) (domain.Conversation, error) {
	conversation, err := uc.repo.FindByID(conversationId, userId)
	if err != nil {
		return domain.Conversation{}, err
	}

	return conversation, nil
}

func (uc *CreateConversation) FindAllConversation(id string) ([]ConversationResponse, error) {
	userId, _ := domain.NewUserID(id)
	conversations, e := uc.repo.FindMyConversation(userId)
	if e != nil {
		return nil, e
	}
	res := make([]ConversationResponse, 0, len(conversations))
	for _, conv := range conversations {
		res = append(res, ConversationResponse{
			ConversationID: string(conv.ID),
			Participants:   userIDsToStrings(conv.Participants),
			CreatedAt:      conv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return res, e
}

// NormalizeParticipants Define a single canonical order.
func NormalizeParticipants(a, b domain.UserID) (domain.UserID, domain.UserID) {
	if a < b {
		return a, b
	}
	return b, a
}

func userIDsToStrings(userIDs []domain.UserID) []string {
	result := make([]string, len(userIDs))
	for i, id := range userIDs {
		result[i] = string(id)
	}
	return result
}
