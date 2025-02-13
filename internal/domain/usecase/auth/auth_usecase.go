package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rshelekhov/avito-tech-internship/internal/domain"
	"github.com/rshelekhov/avito-tech-internship/internal/domain/entity"
	"github.com/rshelekhov/avito-tech-internship/internal/lib/e"
)

type Usecase struct {
	log         *slog.Logger
	userMgr     UserManager
	tokenMgr    TokenManager
	passwordMgr PasswordManager
}

type (
	UserManager interface {
		GetUserByName(ctx context.Context, username string) (entity.User, error)
		CreateUser(ctx context.Context, user entity.User) error
	}

	TokenManager interface {
		GenerateToken(userID string) (string, error)
	}

	PasswordManager interface {
		PasswordHash(password string) (string, error)
		ValidatePassword(providedPassword, passwordHash string) error
	}
)

func NewUsecase(
	log *slog.Logger,
	userMgr UserManager,
	tokenMgr TokenManager,
	passwordMgr PasswordManager,
) *Usecase {
	return &Usecase{
		log:         log,
		userMgr:     userMgr,
		tokenMgr:    tokenMgr,
		passwordMgr: passwordMgr,
	}
}

func (u *Usecase) Authenticate(ctx context.Context, credentials entity.UserCredentials) (string, error) {
	const op = "usecase.Auth.Authenticate"

	log := u.log.With(slog.String("op", op))

	existingUser, err := u.userMgr.GetUserByName(ctx, credentials.Username)
	if errors.Is(err, domain.ErrUserNotFound) {
		// If user not found in storage, register new user and generate token
		return u.registerNewUser(ctx, credentials)
	}
	if err != nil {
		e.LogError(ctx, log, domain.ErrFailedToGetUser, err)
		return "", domain.ErrFailedToGetUser
	}

	// If user found in storage, authenticate and generate token
	return u.authenticateExistingUser(ctx, existingUser, credentials.Password)
}

func (u *Usecase) registerNewUser(ctx context.Context, credentials entity.UserCredentials) (string, error) {
	const op = "usecase.Auth.registerNewUser"

	log := u.log.With(slog.String("op", op))

	passwordHash, err := u.passwordMgr.PasswordHash(credentials.Password)
	if err != nil {
		e.LogError(ctx, log, domain.ErrFailedToGeneratePasswordHash, err)
		return "", domain.ErrFailedToGeneratePasswordHash
	}

	newUser := entity.NewUser(credentials, passwordHash)

	if err = u.userMgr.CreateUser(ctx, newUser); err != nil {
		e.LogError(ctx, log, domain.ErrFailedToCreateUser, err)
		return "", domain.ErrFailedToCreateUser
	}

	token, err := u.tokenMgr.GenerateToken(newUser.ID)
	if err != nil {
		e.LogError(ctx, log, domain.ErrFailedToGenerateToken, err)
		return "", domain.ErrFailedToGenerateToken
	}

	return token, nil
}

func (u *Usecase) authenticateExistingUser(ctx context.Context, existingUser entity.User, providedPassword string) (string, error) {
	const op = "usecase.Auth.authenticateExistingUser"

	log := u.log.With(slog.String("op", op))

	if err := u.passwordMgr.ValidatePassword(providedPassword, existingUser.PasswordHash); err != nil {
		if errors.Is(err, domain.ErrInvalidPassword) {
			e.LogError(ctx, log, domain.ErrInvalidPassword, err)
			return "", domain.ErrBadRequest
		}

		e.LogError(ctx, log, domain.ErrFailedToValidatePassword, err)
		return "", domain.ErrFailedToValidatePassword
	}

	token, err := u.tokenMgr.GenerateToken(existingUser.ID)
	if err != nil {
		e.LogError(ctx, log, domain.ErrFailedToGenerateToken, err)
		return "", domain.ErrFailedToGenerateToken
	}

	return token, nil
}
