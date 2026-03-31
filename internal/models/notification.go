package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

// Notification represents the notifications table.
type Notification struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	NotificationType *string    `json:"notification_type,omitempty" db:"notification_type"`
	ReferenceID      *uuid.UUID `json:"reference_id,omitempty" db:"reference_id"`
	Content          *string    `json:"content,omitempty" db:"content"`
	IsRead           bool       `json:"is_read" db:"is_read"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}
