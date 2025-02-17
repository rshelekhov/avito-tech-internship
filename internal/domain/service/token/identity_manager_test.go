package token

import (
	"context"
	"testing"

	"github.com/rshelekhov/merch-store/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestTokenService_ExtractUserIDFromContext(t *testing.T) {
	tokenService := setup(t)

	expectedUserID := "test-user-id"

	tests := []struct {
		name           string
		setupContext   func() context.Context
		expectedUserID string
		expectedError  error
	}{
		{
			name: "Success",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), domain.UserIDKey, expectedUserID)
			},
			expectedUserID: expectedUserID,
			expectedError:  nil,
		},
		{
			name: "Error — UserID not found in context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedUserID: "",
			expectedError:  domain.ErrUserIDNotFoundInContext,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()

			userID, err := tokenService.ExtractUserIDFromContext(ctx)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
				require.Empty(t, userID)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedUserID, userID)
			}
		})
	}
}
