package entity

type Transaction struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int    `json:"amount"`
}
type CoinHistory struct {
	Received []Transaction `json:"received"`
	Sent     []Transaction `json:"sent"`
}
