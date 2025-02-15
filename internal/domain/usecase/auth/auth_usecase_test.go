package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/rshelekhov/avito-tech-internship/internal/domain"
	"github.com/rshelekhov/avito-tech-internship/internal/domain/entity"
	"github.com/rshelekhov/avito-tech-internship/internal/domain/usecase/auth/mocks"
	"github.com/rshelekhov/avito-tech-internship/internal/lib/logger/handler/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsecase_Authenticate(t *testing.T) {
	ctx := context.Background()
	logger := slogdiscard.NewDiscardLogger()

	testUser := entity.User{
		ID:           "test-user-id",
		Username:     "testuser",
		PasswordHash: "hashed_password",
	}

	testCreds := entity.UserCredentials{
		Username: "testuser",
		Password: "password123",
	}

	tests := []struct {
		name         string
		credentials  entity.UserCredentials
		mockBehavior func(
			userMgr *mocks.UserManager,
			tokenMgr *mocks.TokenManager,
			passwordMgr *mocks.PasswordManager,
		)
		expectedToken string
		expectedError error
	}{
		{
			name:        "Success - Existing user authentication",
			credentials: testCreds,
			mockBehavior: func(
				userMgr *mocks.UserManager,
				tokenMgr *mocks.TokenManager,
				passwordMgr *mocks.PasswordManager,
			) {
				userMgr.EXPECT().GetUserByName(ctx, testCreds.Username).
					Once().
					Return(testUser, nil)

				passwordMgr.EXPECT().ValidatePassword(testCreds.Password, testUser.PasswordHash).
					Once().
					Return(nil)

				tokenMgr.EXPECT().GenerateToken(testUser.ID).
					Once().
					Return("valid_token", nil)
			},
			expectedToken: "valid_token",
			expectedError: nil,
		},
		{
			name:        "Success - New user registration",
			credentials: testCreds,
			mockBehavior: func(
				userMgr *mocks.UserManager,
				tokenMgr *mocks.TokenManager,
				passwordMgr *mocks.PasswordManager,
			) {
				userMgr.EXPECT().GetUserByName(ctx, testCreds.Username).
					Once().
					Return(entity.User{}, domain.ErrUserNotFound)

				passwordMgr.EXPECT().PasswordHash(testCreds.Password).
					Once().
					Return("hashed_password", nil)

				userMgr.EXPECT().CreateUser(ctx, mock.AnythingOfType("entity.User")).
					Once().
					Return(nil)

				tokenMgr.EXPECT().GenerateToken(mock.AnythingOfType("string")).
					Once().
					Return("new_user_token", nil)
			},
			expectedToken: "new_user_token",
			expectedError: nil,
		},
		{
			name:        "Error - Failed to get user",
			credentials: testCreds,
			mockBehavior: func(
				userMgr *mocks.UserManager,
				tokenMgr *mocks.TokenManager,
				passwordMgr *mocks.PasswordManager,
			) {
				userMgr.EXPECT().GetUserByName(ctx, testCreds.Username).
					Once().
					Return(entity.User{}, errors.New("user manager error"))
			},
			expectedToken: "",
			expectedError: domain.ErrFailedToGetUser,
		},
		{
			name:        "Error - Invalid password",
			credentials: testCreds,
			mockBehavior: func(
				userMgr *mocks.UserManager,
				tokenMgr *mocks.TokenManager,
				passwordMgr *mocks.PasswordManager,
			) {
				userMgr.EXPECT().GetUserByName(ctx, testCreds.Username).
					Once().
					Return(testUser, nil)

				passwordMgr.EXPECT().ValidatePassword(testCreds.Password, testUser.PasswordHash).
					Once().
					Return(domain.ErrInvalidPassword)
			},
			expectedToken: "",
			expectedError: domain.ErrBadRequest,
		},
		{
			name:        "Error - Failed to validate password",
			credentials: testCreds,
			mockBehavior: func(
				userMgr *mocks.UserManager,
				tokenMgr *mocks.TokenManager,
				passwordMgr *mocks.PasswordManager,
			) {
				userMgr.EXPECT().GetUserByName(ctx, testCreds.Username).
					Once().
					Return(testUser, nil)

				passwordMgr.EXPECT().ValidatePassword(testCreds.Password, testUser.PasswordHash).
					Once().
					Return(errors.New("password manager error"))
			},
			expectedToken: "",
			expectedError: domain.ErrFailedToValidatePassword,
		},
		{
			name:        "Error - Failed to generate token for a new user",
			credentials: testCreds,
			mockBehavior: func(
				userMgr *mocks.UserManager,
				tokenMgr *mocks.TokenManager,
				passwordMgr *mocks.PasswordManager,
			) {
				userMgr.EXPECT().GetUserByName(ctx, testCreds.Username).
					Once().
					Return(entity.User{}, domain.ErrUserNotFound)

				passwordMgr.EXPECT().PasswordHash(testCreds.Password).
					Once().
					Return("hashed_password", nil)

				userMgr.EXPECT().CreateUser(ctx, mock.AnythingOfType("entity.User")).
					Once().
					Return(nil)

				tokenMgr.EXPECT().GenerateToken(mock.AnythingOfType("string")).
					Once().
					Return("", errors.New("token manager error"))
			},
			expectedToken: "",
			expectedError: domain.ErrFailedToGenerateToken,
		},
		{
			name:        "Error - Failed to generate token for an existing user",
			credentials: testCreds,
			mockBehavior: func(
				userMgr *mocks.UserManager,
				tokenMgr *mocks.TokenManager,
				passwordMgr *mocks.PasswordManager,
			) {
				userMgr.EXPECT().GetUserByName(ctx, testCreds.Username).
					Once().
					Return(testUser, nil)

				passwordMgr.EXPECT().ValidatePassword(testCreds.Password, testUser.PasswordHash).
					Once().
					Return(nil)

				tokenMgr.EXPECT().GenerateToken(testUser.ID).
					Once().
					Return("", errors.New("token manager error"))
			},
			expectedToken: "",
			expectedError: domain.ErrFailedToGenerateToken,
		},
		{
			name:        "Error - Failed to create new user",
			credentials: testCreds,
			mockBehavior: func(
				userMgr *mocks.UserManager,
				tokenMgr *mocks.TokenManager,
				passwordMgr *mocks.PasswordManager,
			) {
				userMgr.EXPECT().GetUserByName(ctx, testCreds.Username).
					Once().
					Return(entity.User{}, domain.ErrUserNotFound)

				passwordMgr.EXPECT().PasswordHash(testCreds.Password).
					Once().
					Return("hashed_password", nil)

				userMgr.EXPECT().CreateUser(ctx, mock.AnythingOfType("entity.User")).
					Once().
					Return(errors.New("user manager error"))
			},
			expectedToken: "",
			expectedError: domain.ErrFailedToCreateUser,
		},
		{
			name:        "Error - Failed to generate password hash",
			credentials: testCreds,
			mockBehavior: func(
				userMgr *mocks.UserManager,
				tokenMgr *mocks.TokenManager,
				passwordMgr *mocks.PasswordManager,
			) {
				userMgr.EXPECT().GetUserByName(ctx, testCreds.Username).
					Once().
					Return(entity.User{}, domain.ErrUserNotFound)

				passwordMgr.EXPECT().PasswordHash(testCreds.Password).
					Once().
					Return("", errors.New("password manager error"))
			},
			expectedToken: "",
			expectedError: domain.ErrFailedToGeneratePasswordHash,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userMgr := mocks.NewUserManager(t)
			tokenMgr := mocks.NewTokenManager(t)
			passwordMgr := mocks.NewPasswordManager(t)

			tt.mockBehavior(userMgr, tokenMgr, passwordMgr)

			usecase := NewUsecase(logger, userMgr, tokenMgr, passwordMgr)
			token, err := usecase.Authenticate(ctx, tt.credentials)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
				require.Empty(t, token)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedToken, token)
			}
		})
	}
}
