package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/rshelekhov/merch-store/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) PasswordHash(password string) (string, error) {
	const op = "service.token.PasswordHash"

	if password == "" {
		return "", fmt.Errorf("%s: %w", op, domain.ErrPasswordIsNotAllowed)
	}

	hash, err := s.passwordHashBcrypt(password)
	if err != nil {
		return "", fmt.Errorf("%s: failed to hash password: %w", op, err)
	}

	return hash, nil
}

func (s *Service) ValidatePassword(providedPassword, passwordHash string) error {
	const op = "service.token.ValidatePassword"

	if providedPassword == "" {
		return fmt.Errorf("%s: %w", op, domain.ErrPasswordIsNotAllowed)
	}

	if passwordHash == "" {
		return fmt.Errorf("%s: %w", op, domain.ErrPasswordHashIsNotAllowed)
	}

	return s.passwordMatchBcrypt(providedPassword, passwordHash)
}

func (s *Service) passwordHashBcrypt(password string) (string, error) {
	const op = "service.token.passwordHashBcrypt"

	passwordHmac := hmac.New(sha256.New, []byte(s.passwordHash.Pepper))
	_, err := passwordHmac.Write([]byte(password))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	passwordBcrypt, err := bcrypt.GenerateFromPassword(passwordHmac.Sum(nil), s.passwordHash.BcryptCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(passwordBcrypt), nil
}

func (s *Service) passwordMatchBcrypt(providedPassword, passwordHash string) error {
	const op = "service.token.passwordMatchBcrypt"

	passwordHmac := hmac.New(sha256.New, []byte(s.passwordHash.Pepper))
	_, err := passwordHmac.Write([]byte(providedPassword))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), passwordHmac.Sum(nil))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return fmt.Errorf("%s: %w", op, domain.ErrInvalidPassword)
		}
		return fmt.Errorf("%s: failed to compare password hashes: %w", op, err)
	}

	return nil
}
