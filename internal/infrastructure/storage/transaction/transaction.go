package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/rshelekhov/avito-tech-internship/internal/infrastructure/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager struct {
	pool *pgxpool.Pool
}

type Executor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	ExecWithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

func NewManager(dbConn *storage.DBConnection) *Manager {
	return &Manager{
		pool: dbConn.Postgres.Pool,
	}
}

type txKey struct{}

var ErrTransactionNotFoundInCtx = errors.New("transaction not found in context")

func (m *Manager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = fn(txCtx)

	return err
}

func (m *Manager) ExecWithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	if !ok {
		return ErrTransactionNotFoundInCtx
	}

	return fn(tx)
}
