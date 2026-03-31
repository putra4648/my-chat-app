package models

import (
	"time"

	"github.com/google/uuid"
)

// MessageAttachment represents the message_attachments table.
type MessageAttachment struct {
	ID             uuid.UUID `json:"id" db:"id"`
	MessageID      uuid.UUID `json:"message_id" db:"message_id"`
	AttachmentType string    `json:"attachment_type" db:"attachment_type"` // 'text', 'image', 'file', 'link'
	URL            string    `json:"url" db:"url"`
	FileName       *string   `json:"file_name,omitempty" db:"file_name"`
	FileSize       *int64    `json:"file_size,omitempty" db:"file_size"`
	MimeType       *string   `json:"mime_type,omitempty" db:"mime_type"`
	Metadata       []byte    `json:"metadata" db:"metadata"` // JSONB
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
