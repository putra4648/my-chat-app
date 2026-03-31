package models

import (
	"time"

	"github.com/google/uuid"
)

// UserSettings represents the user_settings table.
type UserSettings struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Settings  []byte    `json:"settings" db:"settings"` // JSONB
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
