package token

import (
	"testing"

	"github.com/rshelekhov/avito-tech-internship/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestTokenService_PasswordHash(t *testing.T) {
	tokenService := setup(t)

	password := "test-password"

	tests := []struct {
		name          string
		password      string
		expectedError error
	}{
		{
			name:          "Success",
			password:      password,
			expectedError: nil,
		},
		{
			name:          "Error — Password is not allowed",
			password:      "",
			expectedError: domain.ErrPasswordIsNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := tokenService.PasswordHash(tt.password)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
				require.Empty(t, hashedPassword)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, hashedPassword)
				require.Greater(t, len(hashedPassword), 0)
			}
		})
	}
}

func TestTokenService_ValidatePassword(t *testing.T) {
	tokenService := setup(t)

	password := "test-password"

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		setupHash      bool
		expectedError  error
	}{
		{
			name:          "Success",
			password:      password,
			setupHash:     true,
			expectedError: nil,
		},
		{
			name:          "Error – Empty password",
			password:      "",
			setupHash:     true,
			expectedError: domain.ErrPasswordIsNotAllowed,
		},
		{
			name:           "Error — Empty hash",
			hashedPassword: "",
			password:       password,
			expectedError:  domain.ErrPasswordHashIsNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var hashedPassword string
			if tt.setupHash {
				var err error
				hashedPassword, err = tokenService.PasswordHash(password)
				require.NoError(t, err)
				require.NotEmpty(t, hashedPassword)
			} else {
				hashedPassword = tt.hashedPassword
			}

			err := tokenService.ValidatePassword(tt.password, hashedPassword)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
