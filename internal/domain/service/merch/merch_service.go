package merch

import (
	"context"

	"github.com/rshelekhov/avito-tech-internship/internal/domain/entity"
)

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

type Storage interface {
	GetMerchByName(ctx context.Context, itemName string) (entity.Merch, error)
	AddToInventory(ctx context.Context, userID, merchID string) error
}

func (s *Service) GetMerchByName(ctx context.Context, itemName string) (entity.Merch, error) {
	return s.storage.GetMerchByName(ctx, itemName)
}

func (s *Service) AddToInventory(ctx context.Context, userID, merchID string) error {
	return s.storage.AddToInventory(ctx, userID, merchID)
}
