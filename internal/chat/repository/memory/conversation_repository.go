package memory

import (
	"sync"

	"github.com/hardikm9850/GoChat/internal/chat/domain"
	"github.com/hardikm9850/GoChat/internal/chat/repository"
)

type ConversationRepositoryInMemory struct {
	mu            sync.Mutex
	conversations map[domain.ConversationID]domain.Conversation
}

// Create saves a new conversation
func (c *ConversationRepositoryInMemory) Create(conversation domain.Conversation) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conversations == nil {
		c.conversations = make(map[domain.ConversationID]domain.Conversation)
	}
	c.conversations[conversation.ID] = conversation
	return nil
}

// FindByID returns a conversation by ID
func (c *ConversationRepositoryInMemory) FindByID(conversationID domain.ConversationID) (*domain.Conversation, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	conv, ok := c.conversations[conversationID]
	if !ok {
		return nil, ErrConversationNotFound
	}
	return &conv, nil
}

// Find looks for a conversation between two users
func (c *ConversationRepositoryInMemory) Find(userA, userB domain.UserID) (*domain.Conversation, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, conv := range c.conversations {
		if len(conv.Participants) == 2 &&
			(conv.Participants[0] == userA && conv.Participants[1] == userB ||
				conv.Participants[0] == userB && conv.Participants[1] == userA) {
			return &conv, nil
		}
	}
	return nil, ErrConversationNotFound
}

// Participants returns user IDs of a conversation
func (c *ConversationRepositoryInMemory) Participants(id domain.ConversationID) ([]domain.UserID, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	conv, ok := c.conversations[id]
	if !ok {
		return nil, ErrConversationNotFound
	}
	return conv.Participants, nil
}

// New returns a new in-memory repository
func New() repository.ConversationRepository {
	return &ConversationRepositoryInMemory{
		conversations: make(map[domain.ConversationID]domain.Conversation),
	}
}
