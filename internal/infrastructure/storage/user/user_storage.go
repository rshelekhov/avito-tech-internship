package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rshelekhov/avito-tech-internship/internal/domain/entity"
	"github.com/rshelekhov/avito-tech-internship/internal/infrastructure/storage"
	"github.com/rshelekhov/avito-tech-internship/internal/infrastructure/storage/user/sqlc"
)

type Storage struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (s *Storage) CreateUser(ctx context.Context, user entity.User) error {
	const op = "storage.user.CreateUser"

	params := sqlc.CreateUserParams{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Balance:      int32(user.Balance),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	if err := s.queries.CreateUser(ctx, params); err != nil {
		return fmt.Errorf("%s: failed to create user: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUserByName(ctx context.Context, username string) (entity.User, error) {
	const op = "storage.user.GetUserByName"

	user, err := s.queries.GetUserByName(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, storage.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("%s: failed to get user: %w", op, err)
	}

	return entity.User{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Balance:      int(user.Balance),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (s *Storage) GetUserInfoByID(ctx context.Context, userID string) (entity.UserInfo, error) {
	const op = "storage.user.GetUserInfoByID"

	userInfo, err := s.queries.GetUserInfoByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserInfo{}, storage.ErrUserNotFound
		}
		return entity.UserInfo{}, fmt.Errorf("%s: failed to get user info: %w", op, err)
	}

	// Decode inventory
	var inventory []entity.Item
	if err = json.Unmarshal(userInfo.Inventory, &inventory); err != nil {
		return entity.UserInfo{}, fmt.Errorf("%s: failed to unmarshal inventory: %w", op, err)
	}

	// Decode coin history
	var coinHistory entity.CoinHistory
	if err = json.Unmarshal(userInfo.CoinHistory, &coinHistory); err != nil {
		return entity.UserInfo{}, fmt.Errorf("%s: failed to unmarshal coin history: %w", op, err)
	}

	userInfoTemp := entity.UserInfo{
		ID:          userInfo.ID,
		Coins:       int(userInfo.Coins),
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}

	return userInfoTemp, nil
}

func (s *Storage) GetUserInfoByUsername(ctx context.Context, username string) (entity.UserInfo, error) {
	const op = "storage.user.GetUserInfoByUsername"

	userInfo, err := s.queries.GetUserInfoByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserInfo{}, storage.ErrUserNotFound
		}
		return entity.UserInfo{}, fmt.Errorf("%s: failed to get user info: %w", op, err)
	}

	// Decode inventory
	var inventory []entity.Item
	if err = json.Unmarshal(userInfo.Inventory, &inventory); err != nil {
		return entity.UserInfo{}, fmt.Errorf("%s: failed to unmarshal inventory: %w", op, err)
	}

	// Decode coin history
	var coinHistory entity.CoinHistory
	if err = json.Unmarshal(userInfo.CoinHistory, &coinHistory); err != nil {
		return entity.UserInfo{}, fmt.Errorf("%s: failed to unmarshal coin history: %w", op, err)
	}

	return entity.UserInfo{
		ID:          userInfo.ID,
		Coins:       int(userInfo.Coins),
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}, nil
}
