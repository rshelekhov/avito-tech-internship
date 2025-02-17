package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rshelekhov/merch-store/internal/domain"
	"github.com/rshelekhov/merch-store/internal/domain/entity"
	"github.com/rshelekhov/merch-store/internal/domain/service/user/mocks"
	"github.com/rshelekhov/merch-store/internal/infrastructure/storage"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()
	user := entity.User{
		ID:           "test-user-id",
		Username:     "test-username",
		PasswordHash: "test-password-hash",
		Balance:      entity.DefaultBalance,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tests := []struct {
		name          string
		mockBehavior  func(userStorage *mocks.Storage)
		expectedError error
	}{
		{
			name: "Success",
			mockBehavior: func(userStorage *mocks.Storage) {
				userStorage.EXPECT().CreateUser(ctx, user).
					Once().
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error – Storage error",
			mockBehavior: func(userStorage *mocks.Storage) {
				userStorage.EXPECT().CreateUser(ctx, user).
					Once().
					Return(errors.New("storage error"))
			},
			expectedError: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userStorage := mocks.NewStorage(t)
			tt.mockBehavior(userStorage)

			userService := New(userStorage)
			err := userService.CreateUser(ctx, user)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserService_GetUserByName(t *testing.T) {
	ctx := context.Background()
	username := "test-username"
	expectedUser := entity.User{
		ID:           "test-user-id",
		Username:     username,
		PasswordHash: "test-password-hash",
		Balance:      entity.DefaultBalance,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tests := []struct {
		name          string
		mockBehavior  func(userStorage *mocks.Storage)
		expectedUser  entity.User
		expectedError error
	}{
		{
			name: "Success",
			mockBehavior: func(userStorage *mocks.Storage) {
				userStorage.EXPECT().GetUserByName(ctx, username).
					Once().
					Return(expectedUser, nil)
			},
			expectedUser:  expectedUser,
			expectedError: nil,
		},
		{
			name: "Error – User not found",
			mockBehavior: func(userStorage *mocks.Storage) {
				userStorage.EXPECT().GetUserByName(ctx, username).
					Once().
					Return(entity.User{}, storage.ErrUserNotFound)
			},
			expectedUser:  entity.User{},
			expectedError: domain.ErrUserNotFound,
		},
		{
			name: "Error – Storage error",
			mockBehavior: func(userStorage *mocks.Storage) {
				userStorage.EXPECT().GetUserByName(ctx, username).
					Once().
					Return(entity.User{}, errors.New("storage error"))
			},
			expectedUser:  entity.User{},
			expectedError: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userStorage := mocks.NewStorage(t)
			tt.mockBehavior(userStorage)

			userService := New(userStorage)
			user, err := userService.GetUserByName(ctx, username)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError.Error())
				require.Empty(t, user)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, user)
			}
		})
	}
}

func TestUserService_GetUserInfoByID(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"
	expectedUserInfo := entity.UserInfo{
		ID:    userID,
		Coins: 1000,
		Inventory: []entity.Item{
			{
				Type:     "test-item-type",
				Quantity: 1,
			},
		},
		CoinHistory: entity.CoinHistory{
			Received: []entity.Transaction{
				{
					FromUser: "test-from-username",
					ToUser:   "test-username",
					Amount:   100,
					Date:     time.Now().Add(-time.Hour),
				},
			},
			Sent: []entity.Transaction{
				{
					FromUser: "test-username",
					ToUser:   "test-to-username",
					Amount:   100,
					Date:     time.Now().Add(-time.Hour),
				},
			},
		},
	}

	tests := []struct {
		name             string
		mockBehavior     func(userStorage *mocks.Storage)
		expectedUserInfo entity.UserInfo
		expectedError    error
	}{
		{
			name: "Success",
			mockBehavior: func(userStorage *mocks.Storage) {
				userStorage.EXPECT().GetUserInfoByID(ctx, userID).
					Once().
					Return(expectedUserInfo, nil)
			},
			expectedUserInfo: expectedUserInfo,
			expectedError:    nil,
		},
		{
			name: "Error – User not found",
			mockBehavior: func(userStorage *mocks.Storage) {
				userStorage.EXPECT().GetUserInfoByID(ctx, userID).
					Once().
					Return(entity.UserInfo{}, storage.ErrUserNotFound)
			},
			expectedUserInfo: entity.UserInfo{},
			expectedError:    domain.ErrUserNotFound,
		},
		{
			name: "Error – Storage error",
			mockBehavior: func(userStorage *mocks.Storage) {
				userStorage.EXPECT().GetUserInfoByID(ctx, userID).
					Once().
					Return(entity.UserInfo{}, errors.New("storage error"))
			},
			expectedUserInfo: entity.UserInfo{},
			expectedError:    errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userStorage := mocks.NewStorage(t)
			tt.mockBehavior(userStorage)

			userService := New(userStorage)
			user, err := userService.GetUserInfoByID(ctx, userID)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError.Error())
				require.Empty(t, user)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, user)
			}
		})
	}
}

func TestUserService_GetUserInfoByUsername(t *testing.T) {
	ctx := context.Background()
	username := "test-username"
	expectedUserInfo := entity.UserInfo{
		ID:    "test-user-id",
		Coins: 1000,
		Inventory: []entity.Item{
			{
				Type:     "test-item-type",
				Quantity: 1,
			},
		},
		CoinHistory: entity.CoinHistory{
			Received: []entity.Transaction{
				{
					FromUser: "test-from-username",
					ToUser:   username,
					Amount:   100,
					Date:     time.Now().Add(-time.Hour),
				},
			},
			Sent: []entity.Transaction{
				{
					FromUser: username,
					ToUser:   "test-to-username",
					Amount:   100,
					Date:     time.Now().Add(-time.Hour),
				},
			},
		},
	}

	tests := []struct {
		name             string
		mockBehavior     func(userStorage *mocks.Storage)
		expectedUserInfo entity.UserInfo
		expectedError    error
	}{
		{
			name: "Success",
			mockBehavior: func(userStorage *mocks.Storage) {
				userStorage.EXPECT().GetUserInfoByUsername(ctx, username).
					Once().
					Return(expectedUserInfo, nil)
			},
			expectedUserInfo: expectedUserInfo,
			expectedError:    nil,
		},
		{
			name: "Error – User not found",
			mockBehavior: func(userStorage *mocks.Storage) {
				userStorage.EXPECT().GetUserInfoByUsername(ctx, username).
					Once().
					Return(entity.UserInfo{}, storage.ErrUserNotFound)
			},
			expectedUserInfo: entity.UserInfo{},
			expectedError:    domain.ErrUserNotFound,
		},
		{
			name: "Error – Storage error",
			mockBehavior: func(userStorage *mocks.Storage) {
				userStorage.EXPECT().GetUserInfoByUsername(ctx, username).
					Once().
					Return(entity.UserInfo{}, errors.New("storage error"))
			},
			expectedUserInfo: entity.UserInfo{},
			expectedError:    errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userStorage := mocks.NewStorage(t)
			tt.mockBehavior(userStorage)

			userService := New(userStorage)
			user, err := userService.GetUserInfoByUsername(ctx, username)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError.Error())
				require.Empty(t, user)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, user)
			}
		})
	}
}
