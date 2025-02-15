package merch

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/rshelekhov/avito-tech-internship/internal/domain"
	"github.com/rshelekhov/avito-tech-internship/internal/domain/entity"
	"github.com/rshelekhov/avito-tech-internship/internal/domain/service/merch/mocks"
	"github.com/rshelekhov/avito-tech-internship/internal/infrastructure/storage"
	"github.com/stretchr/testify/require"
)

func TestMerchService_GetMerchByName(t *testing.T) {
	ctx := context.Background()
	itemName := "test-item-name"

	expectedMerch := entity.Merch{
		ID:    "test-merch-id",
		Name:  itemName,
		Price: 10,
	}

	tests := []struct {
		name          string
		mockBehavior  func(coinsStorage *mocks.Storage)
		expectedMerch entity.Merch
		expectedError error
	}{
		{
			name: "Success",
			mockBehavior: func(coinsStorage *mocks.Storage) {
				coinsStorage.EXPECT().GetMerchByName(ctx, itemName).
					Once().
					Return(expectedMerch, nil)
			},
			expectedMerch: expectedMerch,
			expectedError: nil,
		},
		{
			name: "Error – Merch not found",
			mockBehavior: func(coinsStorage *mocks.Storage) {
				coinsStorage.EXPECT().GetMerchByName(ctx, itemName).
					Once().
					Return(entity.Merch{}, storage.ErrMerchNotFound)
			},
			expectedMerch: entity.Merch{},
			expectedError: domain.ErrMerchNotFound,
		},
		{
			name: "Error – Storage error",
			mockBehavior: func(coinsStorage *mocks.Storage) {
				coinsStorage.EXPECT().GetMerchByName(ctx, itemName).
					Once().
					Return(entity.Merch{}, errors.New("storage error"))
			},
			expectedMerch: entity.Merch{},
			expectedError: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merchStorage := mocks.NewStorage(t)
			tt.mockBehavior(merchStorage)

			merchService := New(merchStorage)
			merch, err := merchService.GetMerchByName(ctx, itemName)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError.Error())
				require.Empty(t, merch)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, merch)
			}
		})
	}
}

func TestMerchService_AddToInventory(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"
	merchID := "test-merch-id"

	tests := []struct {
		name          string
		mockBehavior  func(coinsStorage *mocks.Storage)
		expectedError error
	}{
		{
			name: "Success",
			mockBehavior: func(coinsStorage *mocks.Storage) {
				coinsStorage.EXPECT().AddToInventory(ctx, mock.AnythingOfType("entity.Purchase")).
					Once().
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error – Storage error",
			mockBehavior: func(coinsStorage *mocks.Storage) {
				coinsStorage.EXPECT().AddToInventory(ctx, mock.AnythingOfType("entity.Purchase")).
					Once().
					Return(errors.New("storage error"))
			},
			expectedError: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merchStorage := mocks.NewStorage(t)
			tt.mockBehavior(merchStorage)

			merchService := New(merchStorage)
			err := merchService.AddToInventory(ctx, userID, merchID)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
