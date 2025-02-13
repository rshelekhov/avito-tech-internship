package entity

import "time"

type Purchase struct {
	ID        string
	UserID    string
	MerchID   string
	CreatedAt time.Time
}
