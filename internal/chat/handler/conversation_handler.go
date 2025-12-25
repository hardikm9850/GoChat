package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hardikm9850/GoChat/internal/chat/domain"
	"github.com/hardikm9850/GoChat/internal/chat/repository"
	"github.com/hardikm9850/GoChat/internal/chat/usecase"
	"net/http"
)

type ConversationHandler struct {
	conversationRepository repository.ConversationRepository
	conversationUseCase    usecase.CreateConversation
}

func NewConversationHandler(repo repository.ConversationRepository, createConversation usecase.CreateConversation) *ConversationHandler {
	return &ConversationHandler{conversationRepository: repo, conversationUseCase: createConversation}
}

type CreateConversationRequest struct {
	ParticipantID string `json:"participant_id" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *ConversationHandler) CreateConversation(c *gin.Context) {
	currentUserID := c.GetString("userID")

	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "unauthorized",
		})
		return
	}

	// Parse request body
	var req CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "participant_id is required",
		})
		return
	}

	// Validate: can't create conversation with yourself
	if req.ParticipantID == currentUserID {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "cannot create conversation with yourself",
		})
		return
	}

	// Check if conversation already exists between these two users
	creatorId, _ := domain.NewUserID(currentUserID)
	participantId, _ := domain.NewUserID(req.ParticipantID)

	conversation, err := h.conversationUseCase.CreateConversation(creatorId, participantId)

	if err == nil {
		// Conversation already exists, return it
		c.JSON(http.StatusOK, usecase.ConversationResponse{
			ConversationID: string(conversation.ID),
			Participants:   userIDsToStrings(conversation.Participants),
			CreatedAt:      conversation.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
		return
	}
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: err.Error(),
	})
}

// GetConversation retrieves a specific conversation
func (h *ConversationHandler) GetConversation(c *gin.Context) {
	currentUserID := c.GetString("userID")
	conversationID := c.Param("id")
	userId, _ := domain.NewUserID(currentUserID)
	conversation, err := h.conversationUseCase.FindConversation(domain.ConversationID(conversationID), userId)

	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "conversation not found",
		})
		return
	}

	c.JSON(http.StatusOK, usecase.ConversationResponse{
		ConversationID: string(conversation.ID),
		Participants:   userIDsToStrings(conversation.Participants),
		CreatedAt:      conversation.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *ConversationHandler) GetMyConversations(c *gin.Context) {
	currentUserID := c.GetString("userID")
	conversation, err := h.conversationUseCase.FindAllConversation(currentUserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to retrieve conversations",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"conversations": conversation,
	})
}

func userIDsToStrings(userIDs []domain.UserID) []string {
	result := make([]string, len(userIDs))
	for i, id := range userIDs {
		result[i] = string(id)
	}
	return result
}
