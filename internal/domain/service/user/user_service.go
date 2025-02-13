package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/rshelekhov/avito-tech-internship/internal/domain"
	"github.com/rshelekhov/avito-tech-internship/internal/domain/entity"
	"github.com/rshelekhov/avito-tech-internship/internal/infrastructure/storage"
)

type Service struct {
	storage Storage
}

type Storage interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByName(ctx context.Context, username string) (entity.User, error)
	GetUserInfoByID(ctx context.Context, userID string) (entity.UserInfo, error)
	GetUserInfoByUsername(ctx context.Context, username string) (entity.UserInfo, error)
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) CreateUser(ctx context.Context, user entity.User) error {
	const op = "service.user.CreateUser"

	if err := s.storage.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) GetUserByName(ctx context.Context, username string) (entity.User, error) {
	const op = "service.user.GetUserByName"

	user, err := s.storage.GetUserByName(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return entity.User{}, domain.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Service) GetUserInfoByID(ctx context.Context, userID string) (entity.UserInfo, error) {
	const op = "service.user.GetUserInfoByID"

	userInfo, err := s.storage.GetUserInfoByID(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return entity.UserInfo{}, domain.ErrUserNotFound
		}
		return entity.UserInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	return userInfo, nil
}

func (s *Service) GetUserInfoByUsername(ctx context.Context, toUsername string) (entity.UserInfo, error) {
	const op = "service.user.GetUserInfoByUsername"

	userInfo, err := s.storage.GetUserInfoByUsername(ctx, toUsername)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return entity.UserInfo{}, domain.ErrUserNotFound
		}
		return entity.UserInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	return userInfo, nil
}
