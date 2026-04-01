package repositories

import (
	"context"
	"errors"
	"putra4648/my-chat-app/internal/models"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository struct {
	pool *pgxpool.Pool
}

func NewChatRepository(pool *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{pool: pool}
}

// FindOrCreatePrivateConversation finds or creates a private conversation between two users
func (r *ChatRepository) FindOrCreatePrivateConversation(ctx context.Context, senderID, receiverID uuid.UUID) (uuid.UUID, error) {
	var conversationID uuid.UUID

	// Check if a private conversation already exists between these two
	query := `
		SELECT p1.conversation_id 
		FROM participants p1
		JOIN participants p2 ON p1.conversation_id = p2.conversation_id
		JOIN conversations c ON p1.conversation_id = c.id
		WHERE c.type = 'private' 
		  AND p1.user_id = $1 
		  AND p2.user_id = $2
	`
	err := r.pool.QueryRow(ctx, query, senderID, receiverID).Scan(&conversationID)
	if err == nil {
		return conversationID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, err
	}

	// Create new conversation
	newID, _ := uuid.NewV4()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO conversations (id, type, created_at) VALUES ($1, 'private', $2)`, newID, time.Now())
	if err != nil {
		return uuid.Nil, err
	}

	// Add participants
	p1ID, _ := uuid.NewV4()
	p2ID, _ := uuid.NewV4()
	_, err = tx.Exec(ctx, `INSERT INTO participants (id, conversation_id, user_id, joined_at) VALUES ($1, $2, $3, $4)`, p1ID, newID, senderID, time.Now())
	if err != nil {
		return uuid.Nil, err
	}
	_, err = tx.Exec(ctx, `INSERT INTO participants (id, conversation_id, user_id, joined_at) VALUES ($1, $2, $3, $4)`, p2ID, newID, receiverID, time.Now())
	if err != nil {
		return uuid.Nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return newID, nil
}

// ChatSummary represents a summary of a conversation for the user interface
type ChatSummary struct {
	ID              uuid.UUID `json:"id"`
	Type            string    `json:"type"`
	OtherUserID     uuid.UUID `json:"other_user_id"`
	OtherUserName   string    `json:"other_user_name"`
	OtherUserEmail  string    `json:"other_user_email"`
	LastMessage     string    `json:"last_message"`
	LastMessageTime time.Time `json:"last_message_time"`
	UnreadCount     int       `json:"unread_count"`
}

func (r *ChatRepository) GetUserConversations(ctx context.Context, userID uuid.UUID) ([]ChatSummary, error) {
	query := `
		SELECT 
			c.id, 
			c.type, 
			u.id as other_user_id,
			u.username as other_user_name, 
			u.email as other_user_email,
			COALESCE(m.content, 'No messages yet') as last_message,
			COALESCE(m.created_at, c.created_at) as last_message_time,
			0 as unread_count
		FROM conversations c
		JOIN participants p1 ON c.id = p1.conversation_id AND p1.user_id = $1
		JOIN participants p2 ON c.id = p2.conversation_id AND p2.user_id != $1
		JOIN users u ON p2.user_id = u.id
		LEFT JOIN LATERAL (
			SELECT content, created_at FROM messages 
			WHERE conversation_id = c.id 
			ORDER BY created_at DESC LIMIT 1
		) m ON true
		ORDER BY last_message_time DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []ChatSummary
	for rows.Next() {
		var s ChatSummary
		err := rows.Scan(&s.ID, &s.Type, &s.OtherUserID, &s.OtherUserName, &s.OtherUserEmail, &s.LastMessage, &s.LastMessageTime, &s.UnreadCount)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}

func (r *ChatRepository) SaveMessage(ctx context.Context, msg *models.Message) error {
	if msg.ID == uuid.Nil {
		msg.ID, _ = uuid.NewV4()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}
	query := `INSERT INTO messages (id, conversation_id, sender_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, query, msg.ID, msg.ConversationID, msg.SenderID, msg.Content, msg.CreatedAt)
	return err
}

func (r *ChatRepository) GetConversationMessages(ctx context.Context, conversationID uuid.UUID) ([]models.Message, error) {
	var messages []models.Message
	query := `SELECT id, conversation_id, sender_id, content, created_at FROM messages WHERE conversation_id = $1 ORDER BY created_at ASC`
	rows, err := r.pool.Query(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Message
		err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
