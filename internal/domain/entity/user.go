package entity

import (
	"time"

	"github.com/segmentio/ksuid"
)

type (
	UserCredentials struct {
		Username string
		Password string
	}

	User struct {
		ID           string
		Username     string
		PasswordHash string
		Balance      int
		CreatedAt    time.Time
		UpdatedAt    time.Time
		DeletedAt    time.Time
	}

	UserInfo struct {
		ID          string
		Coins       int
		Inventory   []Item
		CoinHistory CoinHistory
	}
)

// TODO: perhaps need to move it to env file
const defaultBalance = 1000

func NewUser(credentials UserCredentials, passwordHash string) User {
	return User{
		ID:           ksuid.New().String(),
		Username:     credentials.Username,
		PasswordHash: passwordHash,
		Balance:      defaultBalance,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
