package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

// Conversation represents the conversations table.
type Conversation struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      *string   `json:"name,omitempty" db:"name"`
	Type      string    `json:"type" db:"type"` // 'private' or 'group'
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
