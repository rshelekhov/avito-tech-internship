package coins

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rshelekhov/merch-store/internal/domain"
	"github.com/rshelekhov/merch-store/internal/domain/entity"
	"github.com/rshelekhov/merch-store/internal/domain/usecase/coins/mocks"
	"github.com/rshelekhov/merch-store/internal/lib/logger/handler/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsecase_GetUserInfo(t *testing.T) {
	ctx := context.Background()
	logger := slogdiscard.NewDiscardLogger()

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
		name         string
		mockBehavior func(
			identityMgr *mocks.IdentityManager,
			userMgr *mocks.UserManager,
		)
		expectedInfo  entity.UserInfo
		expectedError error
	}{
		{
			name: "Success",
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(expectedUserInfo.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, expectedUserInfo.ID).
					Once().
					Return(expectedUserInfo, nil)
			},
			expectedInfo:  expectedUserInfo,
			expectedError: nil,
		},
		{
			name: "Error - Failed to extract userID",
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return("", errors.New("identity manager error"))
			},
			expectedInfo:  entity.UserInfo{},
			expectedError: domain.ErrFailedToExtractUserIDFromContext,
		},
		{
			name: "Error - User not found",
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(expectedUserInfo.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, expectedUserInfo.ID).
					Once().
					Return(entity.UserInfo{}, domain.ErrUserNotFound)
			},
			expectedInfo:  entity.UserInfo{},
			expectedError: domain.ErrBadRequest,
		},
		{
			name: "Error - Failed to get user info",
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(expectedUserInfo.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, expectedUserInfo.ID).
					Once().
					Return(entity.UserInfo{}, errors.New("user manager error"))
			},
			expectedInfo:  entity.UserInfo{},
			expectedError: domain.ErrFailedToGetUserInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identityMgr := mocks.NewIdentityManager(t)
			userMgr := mocks.NewUserManager(t)
			coinsMgr := mocks.NewCoinManager(t)
			merchMgr := mocks.NewMerchManager(t)
			txMgr := mocks.NewTransactionManager(t)

			tt.mockBehavior(identityMgr, userMgr)

			usecase := NewUsecase(logger, identityMgr, userMgr, coinsMgr, merchMgr, txMgr)
			info, err := usecase.GetUserInfo(ctx)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
				require.Empty(t, info)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedInfo, info)
			}
		})
	}
}

func TestUsecase_SendCoin(t *testing.T) {
	ctx := context.Background()
	logger := slogdiscard.NewDiscardLogger()

	senderUsername := "test-sender-username"
	receiverUsername := "test-receiver-username"

	sender := entity.UserInfo{
		ID:    "test-sender-id",
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
					ToUser:   senderUsername,
					Amount:   100,
					Date:     time.Now().Add(-time.Hour),
				},
			},
			Sent: []entity.Transaction{
				{
					FromUser: senderUsername,
					ToUser:   "test-to-username",
					Amount:   100,
					Date:     time.Now().Add(-time.Hour),
				},
			},
		},
	}

	receiver := entity.UserInfo{
		ID:    "test-receiver-id",
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
					ToUser:   receiverUsername,
					Amount:   100,
					Date:     time.Now().Add(-time.Hour),
				},
			},
			Sent: []entity.Transaction{
				{
					FromUser: receiverUsername,
					ToUser:   "test-to-username",
					Amount:   100,
					Date:     time.Now().Add(-time.Hour),
				},
			},
		},
	}

	tests := []struct {
		name         string
		toUsername   string
		amount       int
		mockBehavior func(
			identityMgr *mocks.IdentityManager,
			userMgr *mocks.UserManager,
			coinsMgr *mocks.CoinManager,
			txMgr *mocks.TransactionManager,
		)
		expectedError error
	}{
		{
			name:       "Success",
			toUsername: receiverUsername,
			amount:     50,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(sender.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, sender.ID).
					Once().
					Return(sender, nil)

				userMgr.EXPECT().GetUserInfoByUsername(ctx, receiverUsername).
					Once().
					Return(receiver, nil)

				txMgr.EXPECT().WithinTransaction(ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				coinsMgr.EXPECT().UpdateUserCoins(ctx, sender.ID, sender.Coins-50).
					Once().
					Return(nil)

				coinsMgr.EXPECT().UpdateUserCoins(ctx, receiver.ID, receiver.Coins+50).
					Once().
					Return(nil)

				coinsMgr.EXPECT().RegisterCoinTransfer(ctx, mock.AnythingOfType("entity.CoinTransfer")).
					Once().
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:       "Error - Invalid amount",
			toUsername: receiverUsername,
			amount:     0,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				txMgr *mocks.TransactionManager,
			) {
			},
			expectedError: domain.ErrBadRequest,
		},
		{
			name:       "Error — Failed to extract userID from context",
			toUsername: receiverUsername,
			amount:     50,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return("", errors.New("identity manager error"))
			},
			expectedError: domain.ErrFailedToExtractUserIDFromContext,
		},
		{
			name:       "Error — Sender not found",
			toUsername: receiverUsername,
			amount:     50,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(sender.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, sender.ID).
					Once().
					Return(entity.UserInfo{}, domain.ErrUserNotFound)
			},
			expectedError: domain.ErrBadRequest,
		},
		{
			name:       "Error — Failed to get sender info by ID",
			toUsername: receiverUsername,
			amount:     50,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(sender.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, sender.ID).
					Once().
					Return(entity.UserInfo{}, errors.New("user manager error"))
			},
			expectedError: domain.ErrFailedToGetUserInfo,
		},
		{
			name:       "Error - Insufficient coins",
			toUsername: receiverUsername,
			amount:     1500,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(sender.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, sender.ID).
					Once().
					Return(sender, nil)
			},
			expectedError: domain.ErrBadRequest,
		},
		{
			name:       "Error — Receiver not found",
			toUsername: "nonexistent",
			amount:     50,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(sender.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, sender.ID).
					Once().
					Return(sender, nil)

				userMgr.EXPECT().GetUserInfoByUsername(ctx, "nonexistent").
					Once().
					Return(entity.UserInfo{}, domain.ErrUserNotFound)
			},
			expectedError: domain.ErrBadRequest,
		},
		{
			name:       "Error — Failed to get receiver info by ID",
			toUsername: "nonexistent",
			amount:     50,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(sender.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, sender.ID).
					Once().
					Return(sender, nil)

				userMgr.EXPECT().GetUserInfoByUsername(ctx, "nonexistent").
					Once().
					Return(entity.UserInfo{}, errors.New("user manager error"))
			},
			expectedError: domain.ErrFailedToGetUserInfo,
		},
		{
			name:       "Error — Failed to update sender coins",
			toUsername: receiverUsername,
			amount:     50,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(sender.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, sender.ID).
					Once().
					Return(sender, nil)

				userMgr.EXPECT().GetUserInfoByUsername(ctx, receiverUsername).
					Once().
					Return(receiver, nil)

				txMgr.EXPECT().WithinTransaction(ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				coinsMgr.EXPECT().UpdateUserCoins(ctx, sender.ID, sender.Coins-50).
					Once().
					Return(errors.New("coins manager error"))
			},
			expectedError: domain.ErrFailedToUpdateUserCoins,
		},
		{
			name:       "Error — Failed to update receiver coins",
			toUsername: receiverUsername,
			amount:     50,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(sender.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, sender.ID).
					Once().
					Return(sender, nil)

				userMgr.EXPECT().GetUserInfoByUsername(ctx, receiverUsername).
					Once().
					Return(receiver, nil)

				txMgr.EXPECT().WithinTransaction(ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				coinsMgr.EXPECT().UpdateUserCoins(ctx, sender.ID, sender.Coins-50).
					Once().
					Return(nil)

				coinsMgr.EXPECT().UpdateUserCoins(ctx, receiver.ID, receiver.Coins+50).
					Once().
					Return(errors.New("coins manager error"))
			},
			expectedError: domain.ErrFailedToUpdateUserCoins,
		},
		{
			name:       "Error — Failed to register coin transfer",
			toUsername: receiverUsername,
			amount:     50,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(sender.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, sender.ID).
					Once().
					Return(sender, nil)

				userMgr.EXPECT().GetUserInfoByUsername(ctx, receiverUsername).
					Once().
					Return(receiver, nil)

				txMgr.EXPECT().WithinTransaction(ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				coinsMgr.EXPECT().UpdateUserCoins(ctx, sender.ID, sender.Coins-50).
					Once().
					Return(nil)

				coinsMgr.EXPECT().UpdateUserCoins(ctx, receiver.ID, receiver.Coins+50).
					Once().
					Return(nil)

				coinsMgr.EXPECT().RegisterCoinTransfer(ctx, mock.AnythingOfType("entity.CoinTransfer")).
					Once().
					Return(errors.New("coins manager error"))
			},
			expectedError: domain.ErrFailedToRegisterCoinTransfer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identityMgr := mocks.NewIdentityManager(t)
			userMgr := mocks.NewUserManager(t)
			coinsMgr := mocks.NewCoinManager(t)
			merchMgr := mocks.NewMerchManager(t)
			txMgr := mocks.NewTransactionManager(t)

			tt.mockBehavior(identityMgr, userMgr, coinsMgr, txMgr)

			usecase := NewUsecase(logger, identityMgr, userMgr, coinsMgr, merchMgr, txMgr)
			err := usecase.SendCoin(ctx, tt.toUsername, tt.amount)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUsecase_BuyMerch(t *testing.T) {
	ctx := context.Background()
	logger := slogdiscard.NewDiscardLogger()

	testUserInfo := entity.UserInfo{
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

	testMerch := entity.Merch{
		ID:    "merch-id",
		Name:  "Test Merch",
		Price: 50,
	}

	tests := []struct {
		name         string
		itemName     string
		mockBehavior func(
			identityMgr *mocks.IdentityManager,
			userMgr *mocks.UserManager,
			coinsMgr *mocks.CoinManager,
			merchMgr *mocks.MerchManager,
			txMgr *mocks.TransactionManager,
		)
		expectedError error
	}{
		{
			name:     "Success",
			itemName: testMerch.Name,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				merchMgr *mocks.MerchManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(testUserInfo.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, testUserInfo.ID).
					Once().
					Return(testUserInfo, nil)

				merchMgr.EXPECT().GetMerchByName(ctx, testMerch.Name).
					Once().
					Return(testMerch, nil)

				txMgr.EXPECT().WithinTransaction(ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				coinsMgr.EXPECT().UpdateUserCoins(ctx, testUserInfo.ID, testUserInfo.Coins-testMerch.Price).
					Once().
					Return(nil)

				merchMgr.EXPECT().AddToInventory(ctx, testUserInfo.ID, testMerch.ID).
					Once().
					Return(nil)

				coinsMgr.EXPECT().RegisterCoinTransfer(ctx, mock.AnythingOfType("entity.CoinTransfer")).
					Once().
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "Error — Failed to extract userID from context",
			itemName: testMerch.Name,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				merchMgr *mocks.MerchManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return("", errors.New("identity manager error"))
			},
			expectedError: domain.ErrFailedToExtractUserIDFromContext,
		},
		{
			name:     "Error — User not found",
			itemName: testMerch.Name,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				merchMgr *mocks.MerchManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(testUserInfo.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, testUserInfo.ID).
					Once().
					Return(entity.UserInfo{}, domain.ErrUserNotFound)
			},
			expectedError: domain.ErrBadRequest,
		},
		{
			name:     "Error — Failed to get user info",
			itemName: testMerch.Name,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				merchMgr *mocks.MerchManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(testUserInfo.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, testUserInfo.ID).
					Once().
					Return(entity.UserInfo{}, errors.New("user manager error"))
			},
			expectedError: domain.ErrFailedToGetUserInfo,
		},
		{
			name:     "Error - Failed to get merch",
			itemName: "nonexistent",
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				merchMgr *mocks.MerchManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(testUserInfo.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, testUserInfo.ID).
					Once().
					Return(testUserInfo, nil)

				merchMgr.EXPECT().GetMerchByName(ctx, "nonexistent").
					Once().
					Return(entity.Merch{}, domain.ErrFailedToGetMerch)
			},
			expectedError: domain.ErrFailedToGetMerch,
		},
		{
			name:     "Error - Insufficient coins",
			itemName: testMerch.Name,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				merchMgr *mocks.MerchManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(testUserInfo.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, testUserInfo.ID).
					Once().
					Return(entity.UserInfo{ID: testUserInfo.ID, Coins: 20}, nil)

				merchMgr.EXPECT().GetMerchByName(ctx, testMerch.Name).
					Once().
					Return(testMerch, nil)
			},
			expectedError: domain.ErrBadRequest,
		},
		{
			name:     "Error — Failed to update user coins",
			itemName: testMerch.Name,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				merchMgr *mocks.MerchManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(testUserInfo.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, testUserInfo.ID).
					Once().
					Return(testUserInfo, nil)

				merchMgr.EXPECT().GetMerchByName(ctx, testMerch.Name).
					Once().
					Return(testMerch, nil)

				txMgr.EXPECT().WithinTransaction(ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				coinsMgr.EXPECT().UpdateUserCoins(ctx, testUserInfo.ID, testUserInfo.Coins-testMerch.Price).
					Once().
					Return(errors.New("coins manager error"))
			},
			expectedError: domain.ErrFailedToUpdateUserCoins,
		},
		{
			name:     "Error - Failed to add to inventory",
			itemName: testMerch.Name,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				merchMgr *mocks.MerchManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(testUserInfo.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, testUserInfo.ID).
					Once().
					Return(testUserInfo, nil)

				merchMgr.EXPECT().GetMerchByName(ctx, testMerch.Name).
					Once().
					Return(testMerch, nil)

				txMgr.EXPECT().WithinTransaction(ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				coinsMgr.EXPECT().UpdateUserCoins(ctx, testUserInfo.ID, testUserInfo.Coins-testMerch.Price).
					Once().
					Return(nil)

				merchMgr.EXPECT().AddToInventory(ctx, testUserInfo.ID, testMerch.ID).
					Once().
					Return(domain.ErrFailedToAddMerchToInventory)
			},
			expectedError: domain.ErrFailedToAddMerchToInventory,
		},
		{
			name:     "Error - Failed to register coin transfer",
			itemName: testMerch.Name,
			mockBehavior: func(
				identityMgr *mocks.IdentityManager,
				userMgr *mocks.UserManager,
				coinsMgr *mocks.CoinManager,
				merchMgr *mocks.MerchManager,
				txMgr *mocks.TransactionManager,
			) {
				identityMgr.EXPECT().ExtractUserIDFromContext(ctx).
					Once().
					Return(testUserInfo.ID, nil)

				userMgr.EXPECT().GetUserInfoByID(ctx, testUserInfo.ID).
					Once().
					Return(testUserInfo, nil)

				merchMgr.EXPECT().GetMerchByName(ctx, testMerch.Name).
					Once().
					Return(testMerch, nil)

				txMgr.EXPECT().WithinTransaction(ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				coinsMgr.EXPECT().UpdateUserCoins(ctx, testUserInfo.ID, testUserInfo.Coins-testMerch.Price).
					Once().
					Return(nil)

				merchMgr.EXPECT().AddToInventory(ctx, testUserInfo.ID, testMerch.ID).
					Once().
					Return(nil)

				coinsMgr.EXPECT().RegisterCoinTransfer(ctx, mock.AnythingOfType("entity.CoinTransfer")).
					Once().
					Return(errors.New("coins manager error"))
			},
			expectedError: domain.ErrFailedToRegisterCoinTransfer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identityMgr := mocks.NewIdentityManager(t)
			userMgr := mocks.NewUserManager(t)
			coinsMgr := mocks.NewCoinManager(t)
			merchMgr := mocks.NewMerchManager(t)
			txMgr := mocks.NewTransactionManager(t)

			tt.mockBehavior(identityMgr, userMgr, coinsMgr, merchMgr, txMgr)

			usecase := NewUsecase(logger, identityMgr, userMgr, coinsMgr, merchMgr, txMgr)
			err := usecase.BuyMerch(ctx, tt.itemName)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
