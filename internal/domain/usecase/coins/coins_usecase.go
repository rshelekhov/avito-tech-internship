package coins

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rshelekhov/merch-store/internal/domain"
	"github.com/rshelekhov/merch-store/internal/domain/entity"
	"github.com/rshelekhov/merch-store/internal/lib/e"
)

type Usecase struct {
	log         *slog.Logger
	identityMgr IdentityManager
	userMgr     UserManager
	coinsMgr    CoinManager
	merchMgr    MerchManager
	txMgr       TransactionManager
}

type (
	IdentityManager interface {
		ExtractUserIDFromContext(ctx context.Context) (string, error)
	}

	UserManager interface {
		GetUserInfoByID(ctx context.Context, userID string) (entity.UserInfo, error)
		GetUserInfoByUsername(ctx context.Context, username string) (entity.UserInfo, error)
	}

	CoinManager interface {
		UpdateUserCoins(ctx context.Context, senderID string, amount int) error
		RegisterCoinTransfer(ctx context.Context, ct entity.CoinTransfer) error
	}

	MerchManager interface {
		GetMerchByName(ctx context.Context, itemName string) (entity.Merch, error)
		AddToInventory(ctx context.Context, userID, merchID string) error
	}

	TransactionManager interface {
		WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	}
)

func NewUsecase(
	log *slog.Logger,
	identityMgr IdentityManager,
	userSrv UserManager,
	coinsSrv CoinManager,
	merchSrv MerchManager,
	txMgr TransactionManager,
) *Usecase {
	return &Usecase{
		log:         log,
		identityMgr: identityMgr,
		userMgr:     userSrv,
		coinsMgr:    coinsSrv,
		merchMgr:    merchSrv,
		txMgr:       txMgr,
	}
}

func (u *Usecase) GetUserInfo(ctx context.Context) (entity.UserInfo, error) {
	const op = "usecase.Coins.GetUserInfo"

	log := u.log.With(slog.String("op", op))

	userID, err := u.identityMgr.ExtractUserIDFromContext(ctx)
	if err != nil {
		e.LogError(ctx, log, domain.ErrFailedToExtractUserIDFromContext, err)
		return entity.UserInfo{}, domain.ErrFailedToExtractUserIDFromContext
	}

	userInfo, err := u.userMgr.GetUserInfoByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			e.LogError(ctx, log, domain.ErrUserNotFound, err)
			return entity.UserInfo{}, domain.ErrBadRequest
		}

		e.LogError(ctx, log, domain.ErrFailedToGetUserInfo, err)
		return entity.UserInfo{}, domain.ErrFailedToGetUserInfo
	}

	return userInfo, nil
}

func (u *Usecase) SendCoin(ctx context.Context, toUsername string, amount int) error {
	const op = "usecase.Coins.SendCoin"

	log := u.log.With(slog.String("op", op))

	if amount <= 0 {
		err := fmt.Errorf("%s: %w", op, domain.ErrAmountMustBePositive)
		e.LogError(ctx, log, domain.ErrBadRequest, err)
		return domain.ErrBadRequest
	}

	senderID, err := u.identityMgr.ExtractUserIDFromContext(ctx)
	if err != nil {
		e.LogError(ctx, log, domain.ErrFailedToExtractUserIDFromContext, err)
		return domain.ErrFailedToExtractUserIDFromContext
	}

	senderInfo, err := u.userMgr.GetUserInfoByID(ctx, senderID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			e.LogError(ctx, log, domain.ErrSenderNotFound, err)
			return domain.ErrBadRequest
		}

		e.LogError(ctx, log, domain.ErrFailedToGetUserInfo, err)
		return domain.ErrFailedToGetUserInfo
	}

	if senderInfo.Coins < amount {
		err = fmt.Errorf("%s: %w", op, domain.ErrInsufficientCoins)
		e.LogError(ctx, log, domain.ErrBadRequest, err)
		return domain.ErrBadRequest
	}

	receiverUser, err := u.userMgr.GetUserInfoByUsername(ctx, toUsername)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			e.LogError(ctx, log, domain.ErrReceiverNotFound, err)
			return domain.ErrBadRequest
		}

		e.LogError(ctx, log, domain.ErrFailedToGetUserInfo, err)
		return domain.ErrFailedToGetUserInfo
	}

	if err = u.txMgr.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Update sender coins
		if err = u.coinsMgr.UpdateUserCoins(txCtx, senderID, senderInfo.Coins-amount); err != nil {
			e.LogError(txCtx, log, domain.ErrFailedToUpdateUserCoins, err)
			return domain.ErrFailedToUpdateUserCoins
		}

		// Update receiver coins
		if err = u.coinsMgr.UpdateUserCoins(txCtx, receiverUser.ID, receiverUser.Coins+amount); err != nil {
			e.LogError(txCtx, log, domain.ErrFailedToUpdateUserCoins, err)
			return domain.ErrFailedToUpdateUserCoins
		}

		// Register coin transfer
		ct := entity.NewCoinTransfer(senderID, receiverUser.ID, entity.TransactionTypeTransferCoins, amount, time.Now())

		if err = u.coinsMgr.RegisterCoinTransfer(txCtx, ct); err != nil {
			e.LogError(txCtx, log, domain.ErrFailedToRegisterCoinTransfer, err)
			return domain.ErrFailedToRegisterCoinTransfer
		}

		return nil
	}); err != nil {
		e.LogError(ctx, log, domain.ErrFailedToCommitTransaction, err,
			slog.Any("senderID", senderID),
			slog.Any("receiverID", receiverUser.ID),
		)
		return err
	}

	return nil
}

func (u *Usecase) BuyMerch(ctx context.Context, itemName string) error {
	const op = "usecase.Coins.BuyMerch"

	log := u.log.With(slog.String("op", op))

	userID, err := u.identityMgr.ExtractUserIDFromContext(ctx)
	if err != nil {
		e.LogError(ctx, log, domain.ErrFailedToExtractUserIDFromContext, err)
		return domain.ErrFailedToExtractUserIDFromContext
	}

	userInfo, err := u.userMgr.GetUserInfoByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			e.LogError(ctx, log, domain.ErrUserNotFound, err)
			return domain.ErrBadRequest
		}

		e.LogError(ctx, log, domain.ErrFailedToGetUserInfo, err)
		return domain.ErrFailedToGetUserInfo
	}

	merch, err := u.merchMgr.GetMerchByName(ctx, itemName)
	if err != nil {
		e.LogError(ctx, log, domain.ErrFailedToGetMerch, err)
		return domain.ErrFailedToGetMerch
	}

	if userInfo.Coins < merch.Price {
		err = fmt.Errorf("%s: %w", op, domain.ErrInsufficientCoins)
		e.LogError(ctx, log, domain.ErrBadRequest, err)
		return domain.ErrBadRequest
	}

	if err = u.txMgr.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err = u.coinsMgr.UpdateUserCoins(txCtx, userID, userInfo.Coins-merch.Price); err != nil {
			e.LogError(txCtx, log, domain.ErrFailedToUpdateUserCoins, err)
			return domain.ErrFailedToUpdateUserCoins
		}

		if err = u.merchMgr.AddToInventory(txCtx, userID, merch.ID); err != nil {
			e.LogError(txCtx, log, domain.ErrFailedToAddMerchToInventory, err)
			return domain.ErrFailedToAddMerchToInventory
		}

		// Register coin transfer
		ct := entity.NewCoinTransfer(userID, "", entity.TransactionTypePurchaseMerch, merch.Price, time.Now())

		if err = u.coinsMgr.RegisterCoinTransfer(txCtx, ct); err != nil {
			e.LogError(txCtx, log, domain.ErrFailedToRegisterCoinTransfer, err)
			return domain.ErrFailedToRegisterCoinTransfer
		}

		return nil
	}); err != nil {
		e.LogError(ctx, log, domain.ErrFailedToCommitTransaction, err,
			slog.Any("userID", userID),
		)
		return err
	}
	return nil
}
