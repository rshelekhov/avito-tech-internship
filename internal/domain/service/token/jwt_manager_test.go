package token

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenService_GenerateToken(t *testing.T) {
	tokenService := setup(t)

	userID := "test-user-id"

	token, err := tokenService.GenerateToken(userID)

	require.NoError(t, err)
	require.NotEmpty(t, token)
}
