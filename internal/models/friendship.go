package models

import (
	"time"

	"github.com/google/uuid"
)

// Friendship represents the friendships table.
type Friendship struct {
	ID          uuid.UUID `json:"id" db:"id"`
	RequesterID uuid.UUID `json:"requester_id" db:"requester_id"`
	AddresseeID uuid.UUID `json:"addressee_id" db:"addressee_id"`
	Status      string    `json:"status" db:"status"` // 'pending', 'accepted', 'blocked'
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
