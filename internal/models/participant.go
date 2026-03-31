package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

// Participant represents the participants table.
type Participant struct {
	ID             uuid.UUID `json:"id" db:"id"`
	ConversationID uuid.UUID `json:"conversation_id" db:"conversation_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	Role           string    `json:"role" db:"role"` // 'admin', 'member'
	JoinedAt       time.Time `json:"joined_at" db:"joined_at"`
}
