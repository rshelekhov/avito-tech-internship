package domain

import "errors"

var (
	ErrBadRequest                       = errors.New("bad request")
	ErrUserIDNotFoundInContext          = errors.New("user id not found in context")
	ErrFailedToExtractUserIDFromContext = errors.New("failed to extract user id from context")
	ErrUserNotFound                     = errors.New("user not found")
	ErrFailedToGetUser                  = errors.New("failed to get user")
	ErrFailedToCreateUser               = errors.New("failed to create user")
	ErrFailedToGenerateToken            = errors.New("failed to generate token")
	ErrInvalidPassword                  = errors.New("invalid password")
	ErrFailedToValidatePassword         = errors.New("failed to validate password")
	ErrFailedToGeneratePasswordHash     = errors.New("failed to generate password hash")
	ErrPasswordIsNotAllowed             = errors.New("password is not allowed")
	ErrPasswordHashIsNotAllowed         = errors.New("password hash is not allowed")
	ErrFailedToGetUserInfo              = errors.New("failed to get user info")
	ErrAmountMustBePositive             = errors.New("amount must be positive")
	ErrInsufficientCoins                = errors.New("insufficient coins")
	ErrFailedToUpdateUserCoins          = errors.New("failed to update user coins")
	ErrFailedToRegisterCoinTransfer     = errors.New("failed to register coin transfer")
	ErrReceiverNotFound                 = errors.New("receiver not found")
	ErrSenderNotFound                   = errors.New("sender not found")
	ErrMerchNotFound                    = errors.New("merch not found")
	ErrFailedToGetMerch                 = errors.New("failed to get merch")
	ErrFailedToAddMerchToInventory      = errors.New("failed to add merch to inventory")
	ErrFailedToCommitTransaction        = errors.New("failed to commit transaction")
)
