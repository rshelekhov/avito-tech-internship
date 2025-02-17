package merch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/rshelekhov/merch-store/internal/domain"
	"github.com/rshelekhov/merch-store/internal/domain/entity"
	"github.com/rshelekhov/merch-store/internal/infrastructure/storage"
)

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

type Storage interface {
	GetMerchByName(ctx context.Context, itemName string) (entity.Merch, error)
	AddToInventory(ctx context.Context, purchase entity.Purchase) error
}

func (s *Service) GetMerchByName(ctx context.Context, itemName string) (entity.Merch, error) {
	const op = "service.merch.GetMerchByName"

	merch, err := s.storage.GetMerchByName(ctx, itemName)
	if err != nil {
		if errors.Is(err, storage.ErrMerchNotFound) {
			return entity.Merch{}, domain.ErrMerchNotFound
		}
		return entity.Merch{}, fmt.Errorf("%s: %w", op, err)
	}

	return merch, nil
}

func (s *Service) AddToInventory(ctx context.Context, userID, merchID string) error {
	const op = "service.merch.AddToInventory"

	purchase := entity.Purchase{
		ID:        ksuid.New().String(),
		UserID:    userID,
		MerchID:   merchID,
		CreatedAt: time.Now(),
	}

	err := s.storage.AddToInventory(ctx, purchase)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
