package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rshelekhov/avito-tech-internship/internal/domain"
)

func (s *Service) GenerateToken(userID string) (string, error) {
	const op = "service.token.GenerateToken"

	claims := jwt.MapClaims{
		domain.UserIDKey:     userID,
		domain.ExpirationKey: time.Now().Add(s.jwt.TTL).Unix(),
		domain.IssuedAtKey:   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwt.Secret))
	if err != nil {
		return "", fmt.Errorf("%s: failed to sign token: %w", op, err)
	}

	return tokenString, nil
}
