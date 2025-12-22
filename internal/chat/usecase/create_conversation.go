package usecase

import (
	"github.com/google/uuid"
	"github.com/hardikm9850/GoChat/internal/chat/domain"
	"github.com/hardikm9850/GoChat/internal/chat/repository"
	"time"
)

type CreateConversation struct {
	clock func() time.Time
	repo  repository.ConversationRepository
}

func New(repo repository.ConversationRepository) *CreateConversation {
	return &CreateConversation{repo: repo, clock: time.Now}
}

func (uc *CreateConversation) Execute(
	initiatorID, participantID domain.UserID) (*domain.Conversation, error) {

	if initiatorID == participantID {
		return nil, domain.ErrSameParticipant
	}

	existingConversation, err := uc.repo.Find(initiatorID, participantID)
	if err == nil {
		// conversation already exist
		return existingConversation, nil
	}

	// conversation doesn't exist between this 2 users so we create a one.

	conversation := &domain.Conversation{
		ID:           domain.NewConversationID(uuid.NewString()),
		Participants: []domain.UserID{initiatorID, participantID},
		Type:         domain.ConversationTypeDirect,
		CreatedAt:    uc.clock(),
	}

	err = uc.repo.Create(*conversation)

	if err == domain.ErrConversationExists {
		// Race condition protection
		return uc.repo.Find(initiatorID, participantID)
	}
	if err != nil {
		return nil, err
	}

	return conversation, nil
}
