package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

// User represents the users table.
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	FirstName    *string   `json:"firstname,omitempty" db:"firstname"`
	LastName     *string   `json:"lastname,omitempty" db:"lastname"`
	PasswordHash string    `json:"-" db:"password_hash"`
	AvatarURL    *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	LastSeen     time.Time `json:"last_seen" db:"last_seen"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
