package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hardikm9850/GoChat/internal/chat/domain"
	"github.com/hardikm9850/GoChat/internal/chat/repository"
)

type SendMessageUseCase struct {
	conversationRepo repository.ConversationRepository
	messageRepo      repository.MessageRepository
	eventPublisher   EventPublisher // interface, not implementation
}

func NewSendMessageUseCase(
	conversationRepo repository.ConversationRepository,
	messageRepo repository.MessageRepository,
	eventPublisher EventPublisher,
) *SendMessageUseCase {
	return &SendMessageUseCase{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		eventPublisher:   eventPublisher,
	}
}

type MessageSentEvent struct {
	Message domain.Message
}

func (uc *SendMessageUseCase) Execute(
	senderID domain.UserID,
	conversationID domain.ConversationID,
	content string,
) (domain.Message, error) {

	if content == "" {
		return domain.Message{}, errors.New("empty message content")
	}

	conversation, err := uc.conversationRepo.FindByID(conversationID)
	if err != nil {
		return domain.Message{}, err
	}

	if !conversation.HasParticipant(senderID) {
		return domain.Message{}, errors.New("sender not part of conversation")
	}

	message := domain.Message{
		ID:             domain.MessageID(uuid.NewString()),
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		CreatedAt:      time.Now().UTC(),
	}

	if err := uc.messageRepo.Save(message); err != nil {
		return domain.Message{}, err
	}

	uc.eventPublisher.Publish(MessageSentEvent{
		Message: message,
	})

	return message, nil
}
