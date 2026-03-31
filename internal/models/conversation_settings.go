package models

import (
	"github.com/google/uuid"
)

// ConversationSettings represents the conversation_settings table.
type ConversationSettings struct {
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	ConversationID uuid.UUID `json:"conversation_id" db:"conversation_id"`
	IsMuted        bool      `json:"is_muted" db:"is_muted"`
	IsPinned       bool      `json:"is_pinned" db:"is_pinned"`
}
