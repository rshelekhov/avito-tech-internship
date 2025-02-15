package token

import (
	"context"
	"fmt"

	"github.com/rshelekhov/avito-tech-internship/internal/domain"
)

func (s *Service) ExtractUserIDFromContext(ctx context.Context) (string, error) {
	const op = "service.token.ExtractUserIDFromContext"

	userID, ok := ctx.Value(domain.UserIDKey).(string)
	if !ok {
		return "", fmt.Errorf("%s: %w", op, domain.ErrUserIDNotFoundInContext)
	}

	return userID, nil
}
