package memory

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/hardikm9850/GoChat/internal/chat/domain"
	"github.com/hardikm9850/GoChat/internal/chat/repository"
)

var ErrConversationNotFound = errors.New("conversation not found")

type InMemoryMessageRepository struct {
	mu       sync.RWMutex
	messages map[domain.ConversationID][]domain.Message
}

func NewInMemoryMessageRepository() *InMemoryMessageRepository {
	return &InMemoryMessageRepository{
		messages: make(map[domain.ConversationID][]domain.Message),
	}
}

func (r *InMemoryMessageRepository) Save(
	senderID domain.UserID,
	conversationID domain.ConversationID,
	content string,
) (domain.Message, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	// For in-memory, assume conversation exists
	// In real DB, this would validate via FK
	msg := domain.Message{
		ID:             domain.NewMessageID(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		CreatedAt:      time.Now().UTC(),
	}

	r.messages[conversationID] = append(
		r.messages[conversationID],
		msg,
	)

	return msg, nil
}

func (r *InMemoryMessageRepository) Find(
	conversationID domain.ConversationID,
	limit int,
	after *time.Time,
	order repository.MessageOrder,
) ([]domain.Message, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	msgs, ok := r.messages[conversationID]
	if !ok {
		return nil, ErrConversationNotFound
	}

	// Defensive copy
	filtered := make([]domain.Message, 0, len(msgs))

	for _, msg := range msgs {
		if after == nil {
			filtered = append(filtered, msg)
			continue
		}

		if order == repository.OrderAsc && msg.CreatedAt.After(*after) {
			filtered = append(filtered, msg)
		}

		if order == repository.OrderDesc && msg.CreatedAt.Before(*after) {
			filtered = append(filtered, msg)
		}
	}

	// Sort
	sort.Slice(filtered, func(i, j int) bool {
		if order == repository.OrderAsc {
			return filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
		}
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	// Apply limit
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}

	return filtered, nil
}
