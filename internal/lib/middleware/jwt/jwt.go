package jwt

import (
	"context"
	"net/http"

	"github.com/rshelekhov/merch-store/internal/domain"
)

type manager struct {
	secret string
}

type (
	Manager interface {
		HTTPMiddleware(next http.Handler) http.Handler
	}
)

func NewManager(secret string) Manager {
	return &manager{
		secret: secret,
	}
}

func (m *manager) toContext(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, domain.UserIDKey, value)
}
