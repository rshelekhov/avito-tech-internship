package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rshelekhov/merch-store/internal/domain/entity"
	"github.com/rshelekhov/merch-store/internal/infrastructure/storage"
	"github.com/rshelekhov/merch-store/internal/infrastructure/storage/user/sqlc"
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

	user, err := s.queries.GetUserByUsername(ctx, username)
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

	return s.getUserInfoByID(ctx, userID, op)
}

func (s *Storage) GetUserInfoByUsername(ctx context.Context, username string) (entity.UserInfo, error) {
	const op = "storage.user.GetUserInfoByUsername"

	userID, err := s.queries.GetUserIDByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserInfo{}, storage.ErrUserNotFound
		}
		return entity.UserInfo{}, fmt.Errorf("%s: failed to get userID by username: %w", op, err)
	}

	return s.getUserInfoByID(ctx, userID, op)
}

func (s *Storage) getUserInfoByID(ctx context.Context, userID, op string) (entity.UserInfo, error) {
	balance, err := s.queries.GetUserBalanceByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserInfo{}, storage.ErrUserNotFound
		}
		return entity.UserInfo{}, fmt.Errorf("%s: failed to get user balance: %w", op, err)
	}

	inventory, err := s.queries.GetUserInventory(ctx, userID)
	if err != nil {
		return entity.UserInfo{}, fmt.Errorf("%s: failed to get user inventory: %w", op, err)
	}

	receivedTxs, err := s.queries.GetReceivedTransactions(ctx, pgtype.Text{
		String: userID,
		Valid:  true,
	})
	if err != nil {
		return entity.UserInfo{}, fmt.Errorf("%s: failed to get received transactions: %w", op, err)
	}

	sentTxs, err := s.queries.GetSentTransactions(ctx, userID)
	if err != nil {
		return entity.UserInfo{}, fmt.Errorf("%s: failed to get sent transactions: %w", op, err)
	}

	return s.assembleUserInfo(balance, inventory, receivedTxs, sentTxs)
}

func (s *Storage) assembleUserInfo(
	balance sqlc.GetUserBalanceByIDRow,
	inventory []sqlc.GetUserInventoryRow,
	receivedTxs []sqlc.GetReceivedTransactionsRow,
	sentTxs []sqlc.GetSentTransactionsRow,
) (entity.UserInfo, error) {
	// Convert inventory to entity.Item slice
	items := make([]entity.Item, len(inventory))
	for i, item := range inventory {
		items[i] = entity.Item{
			Type:     item.Type,
			Quantity: int(item.Quantity),
		}
	}

	received := make([]entity.Transaction, len(receivedTxs))
	for i, tx := range receivedTxs {
		received[i] = entity.Transaction{
			FromUser: tx.FromUser,
			ToUser:   tx.ToUser.String,
			Amount:   int(tx.Amount),
			Date:     tx.Date,
		}
	}

	sent := make([]entity.Transaction, len(sentTxs))
	for i, tx := range sentTxs {
		sent[i] = entity.Transaction{
			FromUser: tx.FromUser,
			ToUser:   tx.ToUser,
			Amount:   int(tx.Amount),
			Date:     tx.Date,
		}
	}

	return entity.UserInfo{
		ID:        balance.ID,
		Coins:     int(balance.Coins),
		Inventory: items,
		CoinHistory: entity.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}, nil
}
