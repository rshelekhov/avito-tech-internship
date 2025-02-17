package coins

import (
	"context"
	"fmt"

	"github.com/rshelekhov/merch-store/internal/domain/entity"

	"github.com/rshelekhov/merch-store/internal/domain"
)

type Service struct {
	storage Storage
}

type Storage interface {
	UpdateUserCoins(ctx context.Context, senderID string, amount int32) error
	RegisterCoinTransfer(ctx context.Context, ct entity.CoinTransfer) error
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
		return fmt.Errorf("%s: failed to update user coins %w", op, err)
	}

	return nil
}

func (s *Service) RegisterCoinTransfer(ctx context.Context, ct entity.CoinTransfer) error {
	const op = "service.Coins.RegisterCoinTransfer"

	err := s.storage.RegisterCoinTransfer(ctx, ct)
	if err != nil {
		return fmt.Errorf("%s: failed to register coin transfer %w", op, err)
	}

	return nil
}
