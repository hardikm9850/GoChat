package usecase

import (
	"github.com/hardikm9850/GoChat/internal/chat/repository"
	"time"

	"github.com/hardikm9850/GoChat/internal/chat/domain"
)

type MessageRepository interface {
	Find(
		conversationID domain.ConversationID,
		limit int,
		after *time.Time,
		order repository.MessageOrder,
	) ([]domain.Message, error)
}

// GetMessagesUseCase handles fetching chat history
type GetMessagesUseCase struct {
	repo MessageRepository
}

func NewGetMessagesUseCase(repo MessageRepository) *GetMessagesUseCase {
	return &GetMessagesUseCase{repo: repo}
}

// Execute fetches messages for a conversation
func (uc *GetMessagesUseCase) Execute(
	conversationID domain.ConversationID,
	limit int,
	after *time.Time,
	order repository.MessageOrder,
) ([]domain.Message, error) {
	return uc.repo.Find(conversationID, limit, after, order)
}
