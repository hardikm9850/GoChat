package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hardikm9850/GoChat/internal/chat/domain"
	"github.com/hardikm9850/GoChat/internal/chat/repository"
	"github.com/hardikm9850/GoChat/internal/chat/usecase"
	"net/http"
)

func NewConversationHandler(repo repository.ConversationRepository, createConversation usecase.CreateConversation) *ConversationHandler {
	return &ConversationHandler{conversationRepository: repo, conversationUseCase: createConversation}
}

type ConversationHandler struct {
	conversationRepository repository.ConversationRepository
	conversationUseCase    usecase.CreateConversation
}

type CreateConversationRequest struct {
	ParticipantID string `json:"participant_id" binding:"required" example:"user_123"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"unauthorized"`
}

// @Summary Create a conversation
// @Description Creates a conversation between the authenticated user and another participant.
// @Tags Conversations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body CreateConversationRequest true "Create Conversation Request"
// @Success 200 {object} usecase.ConversationResponse "Conversation created or already exists"
// @Failure 400 {object} ErrorResponse "Invalid request or cannot create conversation with yourself"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Conversation not found or creation failed"
// @Router /conversations [post]
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
// @Summary Retrieves a specific conversation
// @Description Retrieves a specific conversation by ID for the authenticated user.
// @Tags Conversations
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Conversation ID"
// @Success 200 {object} usecase.ConversationResponse "Conversation retrieved successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Conversation not found"
// @Router /conversations/{id} [get]
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

// @Summary Get my conversations
// @Description Retrieves all conversations for the authenticated user.
// @Tags Conversations
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string][]usecase.ConversationResponse "List of conversations"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Failed to retrieve conversations"
// @Router /conversations [get]
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
