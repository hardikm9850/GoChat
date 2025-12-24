package handler

import (
	"github.com/hardikm9850/GoChat/internal/chat/handler/http/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hardikm9850/GoChat/internal/chat/domain"
	"github.com/hardikm9850/GoChat/internal/chat/usecase"
	"github.com/hardikm9850/GoChat/internal/hub"
)

type MessageHandler struct {
	SendMessageUC *usecase.SendMessageUseCase
	Hub           *hub.Hub
}

func NewMessageHandler(
	uc *usecase.SendMessageUseCase,
	h *hub.Hub,
) *MessageHandler {
	return &MessageHandler{
		SendMessageUC: uc,
		Hub:           h,
	}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	var req dto.SendMessageRequest
	log.Println("Message handler SendMessage invoked")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Later from JWT middleware
	userID := c.GetString("user_id")

	result, err := h.SendMessageUC.Execute(
		domain.UserID(userID),
		domain.ConversationID(req.ConversationID),
		req.Content,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.Hub.Broadcast <- hub.MessageEvent{
		Message:    result.Message,
		Recipients: result.Recipients,
	}

	c.JSON(http.StatusCreated, dto.SendMessageResponse{
		MessageID:  string(result.Message.ID),
		Recipients: result.Recipients,
	})
}
