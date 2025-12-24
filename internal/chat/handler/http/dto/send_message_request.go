package dto

type SendMessageRequest struct {
    ConversationID string `json:"conversation_id" binding:"required"`
    Content        string `json:"content" binding:"required"`
}
