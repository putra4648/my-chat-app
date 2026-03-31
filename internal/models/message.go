package models

import (
	"time"

	"github.com/google/uuid"
)

// Message represents the messages table.
type Message struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	ConversationID uuid.UUID  `json:"conversation_id" db:"conversation_id"`
	SenderID       *uuid.UUID `json:"sender_id,omitempty" db:"sender_id"`
	Content        string     `json:"content" db:"content"`
	ClientMsgID    *string    `json:"client_msg_id,omitempty" db:"client_msg_id"`
	HasAttachments bool       `json:"has_attachments" db:"has_attachments"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}
