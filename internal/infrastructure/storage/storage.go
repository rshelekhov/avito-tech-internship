package storage

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	pgStorage "github.com/rshelekhov/avito-tech-internship/pkg/storage/postgres"
)

type DBConnection struct {
	Postgres *Postgres
}

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewDBConnection(cfg Config) (*DBConnection, error) {
	return newPostgresStorage(cfg)
}

func newPostgresStorage(cfg Config) (*DBConnection, error) {
	const method = "storage.newPostgresStorage"

	pool, err := pgStorage.New(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create new postgres storage: %w", method, err)
	}

	return &DBConnection{
		Postgres: &Postgres{
			Pool: pool,
		},
	}, nil
}

type Config struct {
	Postgres *pgStorage.Config
}

func (d *DBConnection) Close() {
	pgStorage.Close(d.Postgres.Pool)
}
