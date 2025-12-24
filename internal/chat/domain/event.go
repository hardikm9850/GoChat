package domain

type MessageSentEvent struct {
	Message    Message
	Recipients []string
}
