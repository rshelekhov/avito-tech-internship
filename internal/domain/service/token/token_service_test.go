package token

import "testing"

func setup(t *testing.T) *Service {
	return NewService(Config{
		JWT: JWT{
			Secret: "secret",
			TTL:    3600,
		},
		PasswordHash: PasswordHash{
			Pepper:     "red-hot-chili-peppers",
			BcryptCost: 10,
		},
	})
}
