package coins

import (
	"context"
	"fmt"

	"github.com/rshelekhov/avito-tech-internship/internal/domain"
)

type Service struct {
	storage Storage
}

type Storage interface {
	UpdateUserCoins(ctx context.Context, senderID string, amount int32) error
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) UpdateUserCoins(ctx context.Context, senderID string, amount int) error {
	const op = "service.Coins.UpdateUserCoins"

	if amount <= 0 {
		return domain.ErrAmountMustBePositive
	}

	err := s.storage.UpdateUserCoins(ctx, senderID, int32(amount))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
