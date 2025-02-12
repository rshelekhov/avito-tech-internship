package settings

import (
	"time"

	"github.com/rshelekhov/avito-tech-internship/internal/infrastructure/storage"
	"github.com/rshelekhov/avito-tech-internship/pkg/storage/postgres"
)

type Postgres struct {
	ConnURL      string        `mapstructure:"DB_POSTGRES_CONN_URL"`
	ConnPoolSize int           `mapstructure:"DB_POSTGRES_CONN_POOL_SIZE" envDefault:"10"`
	ReadTimeout  time.Duration `mapstructure:"DB_POSTGRES_READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout time.Duration `mapstructure:"DB_POSTGRES_WRITE_TIMEOUT" envDefault:"5s"`
	IdleTimeout  time.Duration `mapstructure:"DB_POSTGRES_IDLE_TIMEOUT" envDefault:"60s"`
	DialTimeout  time.Duration `mapstructure:"DB_POSTGRES_DIAL_TIMEOUT" envDefault:"10s"`
}

func ToStorageConfig(params Postgres) storage.Config {
	return storage.Config{
		Postgres: &postgres.Config{
			ConnURL:      params.ConnURL,
			ConnPoolSize: params.ConnPoolSize,
			ReadTimeout:  params.ReadTimeout,
			WriteTimeout: params.WriteTimeout,
			IdleTimeout:  params.IdleTimeout,
			DialTimeout:  params.DialTimeout,
		},
	}
}
