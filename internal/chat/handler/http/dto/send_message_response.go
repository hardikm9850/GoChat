package dto

type SendMessageResponse struct {
	MessageID  string   `json:"message_id"`
	Recipients []string `json:"recipients"`
}
