package coins

import "context"

type Service struct {
	storage Storage
}

type Storage interface {
	UpdateUserCoins(ctx context.Context, senderID string, amount int) error
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) UpdateUserCoins(ctx context.Context, senderID string, amount int) error {
	return s.storage.UpdateUserCoins(ctx, senderID, amount)
}
