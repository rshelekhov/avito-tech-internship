package transaction

import (
	"context"
	
	"github.com/jackc/pgx/v5"
)

type Manager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	ExecWithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

