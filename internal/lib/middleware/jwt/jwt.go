package jwt

import "context"

type manager struct {
	secret string
}

type (
	ContextManager interface {
		FromContext(ctx context.Context) (string, bool)
		ToContext(ctx context.Context, value string) context.Context
	}
)
