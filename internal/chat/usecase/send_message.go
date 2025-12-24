package usecase

import (
	"errors"
	"github.com/hardikm9850/GoChat/internal/chat/domain"
	"github.com/hardikm9850/GoChat/internal/chat/repository"
	"log"
)

type SendMessageUseCase struct {
	ConversationRepo repository.ConversationRepository
	MessageRepo      repository.MessageRepository
	eventPublisher   EventPublisher // interface, not implementation
}

type SendMessageResult struct {
	Message    domain.Message
	Recipients []string
}

func NewSendMessageUseCase(
	conversationRepo repository.ConversationRepository,
	messageRepo repository.MessageRepository,
	eventPublisher EventPublisher,
) *SendMessageUseCase {
	return &SendMessageUseCase{
		ConversationRepo: conversationRepo,
		MessageRepo:      messageRepo,
		eventPublisher:   eventPublisher,
	}
}

type MessageSentEvent struct {
	Message domain.Message
}

func (usecase *SendMessageUseCase) Execute(
	senderID domain.UserID,
	conversationID domain.ConversationID,
	content string,
) (SendMessageResult, error) {
	log.Println("SendMessageUC Execute called, message content:", content)

	var emptySendMessageResult = SendMessageResult{}
	// ip validation
	if content == "" {
		return emptySendMessageResult, errors.New("empty message content")
	}
	//save message
	message, err := usecase.MessageRepo.Save(
		senderID, conversationID, content,
	)
	if err != nil {
		return SendMessageResult{}, err
	}
	// resolve recipients
	participants, err := usecase.ConversationRepo.Participants(conversationID)
	if err != nil {
		return SendMessageResult{}, err
	}

	recipients := make([]string, 0)
	for _, p := range participants {
		if p != senderID {
			recipients = append(recipients, string(p))
		}
	}
	log.Printf("SendMsg recipients count %d \n\n", len(recipients))
	event := domain.MessageSentEvent{
		Message:    message,
		Recipients: recipients,
	}

	usecase.eventPublisher.Publish(event)

	return SendMessageResult{
		Message:    message,
		Recipients: recipients,
	}, nil
}
