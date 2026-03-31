package services

import (
	"context"
	"putra4648/my-chat-app/internal/models"
	"putra4648/my-chat-app/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return s.userRepo.CreateUser(ctx, user)
}
