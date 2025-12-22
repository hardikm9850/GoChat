package domain

import "errors"

var (
	ErrSameParticipant      = errors.New("participants must be distinct")
	ErrConversationExists   = errors.New("conversation already exists")
	ErrConversationNotFound = errors.New("conversation not found")
	ErrMessageNotFound = errors.New("Message not found")
)
