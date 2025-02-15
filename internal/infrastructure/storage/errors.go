package storage

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrMerchNotFound = errors.New("merch not found")
)
