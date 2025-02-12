package user

import (
	"context"

	"github.com/rshelekhov/avito-tech-internship/internal/domain/entity"
)

type Service struct {
	storage Storage
}

type Storage interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByName(ctx context.Context, username string) (entity.User, error)
	GetUserInfoByID(ctx context.Context, userID string) (entity.UserInfo, error)
	GetUserInfoByUsername(ctx context.Context, toUsername string) (entity.UserInfo, error)
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) CreateUser(ctx context.Context, user entity.User) error {
	return s.storage.CreateUser(ctx, user)
}

func (s *Service) GetUserByName(ctx context.Context, username string) (entity.User, error) {
	return s.storage.GetUserByName(ctx, username)
}

func (s *Service) GetUserInfoByID(ctx context.Context, userID string) (entity.UserInfo, error) {
	return s.storage.GetUserInfoByID(ctx, userID)
}

func (s *Service) GetUserInfoByUsername(ctx context.Context, toUsername string) (entity.UserInfo, error) {
	return s.storage.GetUserInfoByUsername(ctx, toUsername)
}
