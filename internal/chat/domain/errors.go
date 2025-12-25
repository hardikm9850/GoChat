package domain

import "errors"

var (
	ErrSameParticipant              = errors.New("participants must be distinct")
	ErrConversationExists           = errors.New("conversation already exists")
	ErrConversationNotFound         = errors.New("conversation not found")
	ErrMessageNotFound              = errors.New("Message not found")
	ErrUnauthorized                 = errors.New("Unauthorized")
	ErrParticipantIdNotFound        = errors.New("participant_id is required")
	ErrInvalidParticipant           = errors.New("cannot create conversation with yourself")
	ErrConversationNotfound         = errors.New("conversation not found")
	ErrNotParticipantInConversation = errors.New("you are not a participant in this conversation")
)
