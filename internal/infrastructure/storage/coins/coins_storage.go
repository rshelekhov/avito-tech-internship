package coins

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rshelekhov/avito-tech-internship/internal/infrastructure/storage/coins/sqlc"
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

func (s *Storage) UpdateUserCoins(ctx context.Context, senderID string, amount int32) error {
	const op = "storage.coins.UpdateUserCoins"

	params := sqlc.UpdateUserCoinsParams{
		ID:      senderID,
		Balance: amount,
	}

	if err := s.txMgr.ExecWithinTx(ctx, func(tx pgx.Tx) error {
		return s.queries.WithTx(tx).UpdateUserCoins(ctx, params)
	}); err != nil {
		return fmt.Errorf("%s: failed to update user coins: %w", op, err)
	}

	return nil
}
