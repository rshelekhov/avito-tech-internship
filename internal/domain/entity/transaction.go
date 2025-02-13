package entity

import "time"

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
