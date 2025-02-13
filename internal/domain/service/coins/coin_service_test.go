package coins

import (
	"context"
	"errors"
	"testing"

	"github.com/rshelekhov/avito-tech-internship/internal/domain"

	"github.com/rshelekhov/avito-tech-internship/internal/domain/service/coins/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCoinsService_UpdateUserCoins(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"

	tests := []struct {
		name          string
		mockBehavior  func(coinsStorage *mocks.Storage)
		amount        int
		expectedError error
	}{
		{
			name: "Success",
			mockBehavior: func(coinsStorage *mocks.Storage) {
				coinsStorage.EXPECT().UpdateUserCoins(ctx, userID, mock.AnythingOfType("int")).
					Once().
					Return(nil)
			},
			amount:        10,
			expectedError: nil,
		},
		{
			name: "Error – Invalid amount",
			mockBehavior: func(coinsStorage *mocks.Storage) {
			},
			amount:        -10,
			expectedError: domain.ErrAmountMustBePositive,
		},
		{
			name: "Error – Storage error",
			mockBehavior: func(coinsStorage *mocks.Storage) {
				coinsStorage.EXPECT().UpdateUserCoins(ctx, userID, mock.AnythingOfType("int")).
					Once().
					Return(errors.New("storage error"))
			},
			amount:        10,
			expectedError: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coinsStorage := mocks.NewStorage(t)
			tt.mockBehavior(coinsStorage)

			coinsService := New(coinsStorage)
			err := coinsService.UpdateUserCoins(ctx, userID, tt.amount)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
