package entity

import (
	"time"

	"github.com/segmentio/ksuid"
)

type Transaction struct {
	FromUser string    `json:"fromUser,omitempty"`
	ToUser   string    `json:"toUser,omitempty"`
	Amount   int       `json:"amount"`
	Date     time.Time `json:"date"`
}
type CoinHistory struct {
	Received []Transaction `json:"received"`
	Sent     []Transaction `json:"sent"`
}

type TransactionType string

const (
	TransactionTypeTransferCoins TransactionType = "transfer_coins"
	TransactionTypePurchaseMerch TransactionType = "purchase_merch"
)

func (t TransactionType) String() string {
	return string(t)
}

type CoinTransfer struct {
	ID              string
	SenderID        string
	ReceiverID      string
	TransactionType TransactionType
	Amount          int32
	Date            time.Time
}

func NewCoinTransfer(senderID, receiverID string, tt TransactionType, amount int, date time.Time) CoinTransfer {
	return CoinTransfer{
		ID:              ksuid.New().String(),
		SenderID:        senderID,
		ReceiverID:      receiverID,
		TransactionType: tt,
		Amount:          int32(amount),
		Date:            date,
	}
}
