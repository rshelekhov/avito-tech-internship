package coins

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool  *pgxpool.Pool
	txMgr TransactionManager
}

type TransactionManager interface {
	ExecWithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

func NewStorage(pool *pgxpool.Pool, txMgr TransactionManager) *Storage {
	return &Storage{
		pool:  pool,
		txMgr: txMgr,
	}
}

func (s *Storage) UpdateUserCoins(ctx context.Context, senderID string, amount int) error {}
