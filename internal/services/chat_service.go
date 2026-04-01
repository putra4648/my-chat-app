package services

import (
	"context"
	"putra4648/my-chat-app/internal/models"
	"putra4648/my-chat-app/internal/repositories"

	"github.com/gofrs/uuid/v5"
)

type ChatService struct {
	repo *repositories.ChatRepository
}

func NewChatService(repo *repositories.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) GetUserConversations(ctx context.Context, userID uuid.UUID) ([]repositories.ChatSummary, error) {
	return s.repo.GetUserConversations(ctx, userID)
}

func (s *ChatService) GetConversationMessages(ctx context.Context, conversationID uuid.UUID) ([]models.Message, error) {
	return s.repo.GetConversationMessages(ctx, conversationID)
}

func (s *ChatService) SaveMessage(ctx context.Context, msg *models.Message) error {
	return s.repo.SaveMessage(ctx, msg)
}

func (s *ChatService) GetOrCreatePrivateConversation(ctx context.Context, senderID, receiverID uuid.UUID) (uuid.UUID, error) {
	return s.repo.FindOrCreatePrivateConversation(ctx, senderID, receiverID)
}
