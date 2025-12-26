package handler

import (
	"github.com/hardikm9850/GoChat/internal/chat/handler/http/dto"
	"github.com/hardikm9850/GoChat/internal/chat/repository"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hardikm9850/GoChat/internal/chat/domain"
	"github.com/hardikm9850/GoChat/internal/chat/usecase"
	"github.com/hardikm9850/GoChat/internal/hub"
)

type MessagesHandler struct {
	SendMessageUC     *usecase.SendMessageUseCase
	GetMessageUseCase *usecase.GetMessagesUseCase
	Hub               *hub.Hub
}

func NewMessageHandler(
	sendMessageUseCase *usecase.SendMessageUseCase,
	getMessageUseCase *usecase.GetMessagesUseCase,
	h *hub.Hub,
) *MessagesHandler {
	return &MessagesHandler{
		SendMessageUC:     sendMessageUseCase,
		GetMessageUseCase: getMessageUseCase,
		Hub:               h,
	}
}

func (h *MessagesHandler) SendMessage(c *gin.Context) {
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

func (h *MessagesHandler) GetMessages(c *gin.Context) {
	conversationID := domain.ConversationID(c.Param("id"))

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	afterStr := c.Query("after")
	var after *time.Time
	if afterStr != "" {
		t, err := time.Parse(time.RFC3339, afterStr)
		if err == nil {
			after = &t
		}
	}

	orderStr := c.DefaultQuery("order", "desc")
	order := repository.OrderDesc
	if orderStr == "asc" {
		order = repository.OrderAsc
	}

	messages, err := h.GetMessageUseCase.Execute(conversationID, limit, after, order)
	if err != nil {
		if err == domain.ErrConversationNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}
