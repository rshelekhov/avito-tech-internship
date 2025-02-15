package api_tests

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rshelekhov/avito-tech-internship/internal/controller/http/v1/handler"
	"github.com/stretchr/testify/require"
)

func TestSendCoin_HappyPath(t *testing.T) {
	e := newTestAPI(t)

	// Register sender
	senderUsername := gofakeit.Username()
	resp := e.POST("/api/auth").
		WithJSON(handler.AuthRequest{
			Username: senderUsername,
			Password: randomFakePassword(),
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// Get token from response
	token := resp.Value("token").String().Raw()
	require.NotEmpty(t, token)

	// Register receiver
	receiverUsername := gofakeit.Username()
	e.POST("/api/auth").
		WithJSON(handler.AuthRequest{
			Username: receiverUsername,
			Password: randomFakePassword(),
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().NotEmpty()

	// Send coins from sender to receiver
	e.POST("/api/sendCoin").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(handler.SendCoinRequest{
			ToUser: receiverUsername,
			Amount: 100,
		}).
		Expect().
		Status(http.StatusOK)
}

func TestSendCoin_BadRequest(t *testing.T) {
	e := newTestAPI(t)

	// Register sender
	senderUsername := gofakeit.Username()
	resp := e.POST("/api/auth").
		WithJSON(handler.AuthRequest{
			Username: senderUsername,
			Password: randomFakePassword(),
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// Get token from response
	token := resp.Value("token").String().Raw()
	require.NotEmpty(t, token)

	// Register receiver
	receiverUsername := gofakeit.Username()
	e.POST("/api/auth").
		WithJSON(handler.AuthRequest{
			Username: receiverUsername,
			Password: randomFakePassword(),
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().NotEmpty()

	tests := []struct {
		name             string
		receiverUsername string
		amount           int
	}{
		{
			name:             "Error – Receiver not found",
			receiverUsername: gofakeit.Username(),
			amount:           100,
		},
		{
			name:             "Error — Invalid amount",
			receiverUsername: receiverUsername,
			amount:           -100,
		},
		{
			name:             "Error — Insufficient coins",
			receiverUsername: receiverUsername,
			amount:           100000,
		},
		{
			name:             "Error — Empty username",
			receiverUsername: "",
			amount:           100,
		},
		{
			name:             "Error — Empty amount",
			receiverUsername: receiverUsername,
			amount:           0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e.POST("/api/sendCoin").
				WithHeader("Authorization", "Bearer "+token).
				WithJSON(handler.SendCoinRequest{
					ToUser: tt.receiverUsername,
					Amount: tt.amount,
				}).
				Expect().
				Status(http.StatusBadRequest)
		})
	}
}

func TestSendCoin_Unauthorized(t *testing.T) {
	e := newTestAPI(t)

	// Try to send coins without authorization
	e.POST("/api/sendCoin").
		Expect().
		Status(http.StatusUnauthorized)
}
