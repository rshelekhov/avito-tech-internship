package api_tests

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rshelekhov/avito-tech-internship/internal/controller/http/v1/handler"
	"github.com/stretchr/testify/require"
)

func TestGetInfo_HappyPath(t *testing.T) {
	e := newTestAPI(t)

	// Register user
	resp := e.POST("/api/auth").
		WithJSON(handler.AuthRequest{
			Username: gofakeit.Username(),
			Password: randomFakePassword(),
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// Get token from response
	token := resp.Value("token").String().Raw()
	require.NotEmpty(t, token)

	// Buy some merch by 500 coins
	e.GET("/api/buy/{item}", "pink-hoody").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK)

	// Register receiver for sending coins
	receiverUsername := gofakeit.Username()
	amount := 100

	e.POST("/api/auth").
		WithJSON(handler.AuthRequest{
			Username: receiverUsername,
			Password: randomFakePassword(),
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().NotEmpty()

	// Send some coins
	e.POST("/api/sendCoin").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(handler.SendCoinRequest{
			ToUser: receiverUsername,
			Amount: amount,
		}).
		Expect().
		Status(http.StatusOK)

	// Get info about user with details about his inventory and coin history
	resp = e.GET("/api/user").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// Check inventory
	inventory := resp.Value("inventory").Array().Raw()
	require.Len(t, inventory, 1)

	// Check coin history
	coinHistory := resp.Value("coinHistory").Object()
	sentTransactions := coinHistory.Value("sent").Array()
	sentTransactions.Length().IsEqual(1)

	sentTransaction := sentTransactions.Value(0).Object()
	sentTransaction.Value("toUser").String().IsEqual(receiverUsername)
	sentTransaction.Value("amount").Number().IsEqual(float64(amount))
	sentTransaction.Value("date").String().NotEmpty()
}

func TestGetInfo_Unauthorized(t *testing.T) {
	e := newTestAPI(t)

	// Try to get info without authorization
	e.GET("/api/user").
		Expect().
		Status(http.StatusUnauthorized)
}
