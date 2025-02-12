package coins

import (
	"context"
	"errors"
	"fmt"
	"github.com/rshelekhov/avito-tech-internship/internal/domain"
	"github.com/rshelekhov/avito-tech-internship/internal/domain/entity"
	"github.com/rshelekhov/avito-tech-internship/internal/infrastructure/storage"
	"github.com/rshelekhov/avito-tech-internship/internal/lib/e"
	"log/slog"
)

type Usecase struct {
	log         *slog.Logger
	identityMgr IdentityManager
	userMgr     UserManager
	coinsMgr    CoinManager
	merchMgr    MerchManager
}

type (
	IdentityManager interface {
		ExtractUserIDFromContext(ctx context.Context) (string, error)
	}

	UserManager interface {
		GetUserInfoByID(ctx context.Context, userID string) (entity.UserInfo, error)
		GetUserInfoByUsername(ctx context.Context, toUsername string) (entity.UserInfo, error)
	}

	CoinManager interface {
		UpdateUserCoins(ctx context.Context, senderID string, amount int) error
	}

	MerchManager interface {
		GetMerchByName(ctx context.Context, itemName string) (entity.Merch, error)
	}
)

func NewUsecase(
	log *slog.Logger,
	identityMgr IdentityManager,
	userSrv UserManager,
	coinsSrv CoinManager,
	merchSrv MerchManager,
) *Usecase {
	return &Usecase{
		log:         log,
		identityMgr: identityMgr,
		userMgr:     userSrv,
		coinsMgr:    coinsSrv,
		merchMgr:    merchSrv,
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
		if errors.Is(err, storage.ErrUserNotFound) {
			e.LogError(ctx, log, domain.ErrSenderNotFound, err)
			return domain.ErrSenderNotFound
		}
		e.LogError(ctx, log, domain.ErrFailedToGetUserInfo, err)
		return domain.ErrFailedToGetUserInfo
	}

	if senderInfo.Coins < amount {
		err := fmt.Errorf("%s: %w", op, domain.ErrInsufficientCoins)
		e.LogError(ctx, log, domain.ErrBadRequest, err)
		return domain.ErrBadRequest
	}

	receiverUser, err := u.userMgr.GetUserInfoByUsername(ctx, toUsername)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			e.LogError(ctx, log, domain.ErrReceiverNotFound, err)
			return domain.ErrReceiverNotFound
		}
		e.LogError(ctx, log, domain.ErrFailedToGetUserInfo, err)
		return domain.ErrFailedToGetUserInfo
	}

	// TODO: add transaction here

	// Update sender coins
	if err = u.coinsMgr.UpdateUserCoins(ctx, senderID, senderInfo.Coins-amount); err != nil {
		e.LogError(ctx, log, domain.ErrFailedToUpdateUserCoins, err)
		return domain.ErrFailedToUpdateUserCoins
	}

	// Update receiver coins
	if err = u.coinsMgr.UpdateUserCoins(ctx, receiverUser.ID, receiverUser.Coins+amount); err != nil {
		e.LogError(ctx, log, domain.ErrFailedToUpdateUserCoins, err)
		return domain.ErrFailedToUpdateUserCoins
	}

	return nil
}

func (u *Usecase) BuyCoin(ctx context.Context, itemName string) error {
	const op = "usecase.Coins.BuyCoin"

	log := u.log.With(slog.String("op", op))

	userID, err := u.identityMgr.ExtractUserIDFromContext(ctx)
	if err != nil {
		e.LogError(ctx, log, domain.ErrFailedToExtractUserIDFromContext, err)
		return domain.ErrFailedToExtractUserIDFromContext
	}

	userInfo, err := u.userMgr.GetUserInfoByID(ctx, userID)
	if err != nil {
		e.LogError(ctx, log, domain.ErrFailedToGetUserInfo, err)
		return domain.ErrFailedToGetUserInfo
	}

	merch, err := u.merchMgr.GetMerchByName(ctx, itemName)
	if err != nil {
		if errors.Is(err, storage.ErrMerchNotFound) {
			e.LogError(ctx, log, domain.ErrMerchNotFound, err)
			return domain.ErrMerchNotFound
		}
		e.LogError(ctx, log, domain.ErrFailedToGetMerch, err)
		return domain.ErrFailedToGetMerch
	}

	return nil
}
