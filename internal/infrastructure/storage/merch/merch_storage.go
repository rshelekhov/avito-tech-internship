package merch

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rshelekhov/merch-store/internal/domain/entity"
	"github.com/rshelekhov/merch-store/internal/infrastructure/storage"
	"github.com/rshelekhov/merch-store/internal/infrastructure/storage/merch/sqlc"
)

type Storage struct {
	pool    *pgxpool.Pool
	txMgr   TransactionManager
	queries *sqlc.Queries
}

type TransactionManager interface {
	ExecWithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

func NewStorage(pool *pgxpool.Pool, txMgr TransactionManager) *Storage {
	return &Storage{
		pool:    pool,
		txMgr:   txMgr,
		queries: sqlc.New(pool),
	}
}

func (s *Storage) GetMerchByName(ctx context.Context, itemName string) (entity.Merch, error) {
	const op = "storage.merch.GetMerchByName"

	merch, err := s.queries.GetMerchByName(ctx, itemName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Merch{}, storage.ErrMerchNotFound
		}
		return entity.Merch{}, fmt.Errorf("%s: failed to get merch: %w", op, err)
	}

	return entity.Merch{
		ID:    merch.ID,
		Name:  merch.Name,
		Price: int(merch.Price),
	}, nil
}

func (s *Storage) AddToInventory(ctx context.Context, purchase entity.Purchase) error {
	const op = "storage.merch.AddToInventory"

	params := sqlc.AddToInventoryParams{
		ID:        purchase.ID,
		UserID:    purchase.UserID,
		MerchID:   purchase.MerchID,
		CreatedAt: purchase.CreatedAt,
	}

	if err := s.txMgr.ExecWithinTx(ctx, func(tx pgx.Tx) error {
		return s.queries.WithTx(tx).AddToInventory(ctx, params)
	}); err != nil {
		return fmt.Errorf("%s: failed to add to inventory: %w", op, err)
	}

	return nil
}
