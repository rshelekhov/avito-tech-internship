package token

import "time"

type Service struct {
	jwt          JWT
	passwordHash PasswordHash
}

func NewService(cfg Config) *Service {
	return &Service{
		jwt:          cfg.JWT,
		passwordHash: cfg.PasswordHash,
	}
}

type (
	Config struct {
		JWT          JWT
		PasswordHash PasswordHash
	}

	JWT struct {
		Secret string
		TTL    time.Duration
	}

	PasswordHash struct {
		Pepper     string
		BcryptCost int
	}
)
