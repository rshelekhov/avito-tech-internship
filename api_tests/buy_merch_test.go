package api_tests

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rshelekhov/avito-tech-internship/internal/controller/http/v1/handler"
	"github.com/stretchr/testify/require"
)

func TestBuyMerch_HappyPath(t *testing.T) {
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

	// Buy merch
	e.GET("/api/buy/{item}", "book").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK)
}

func TestBuyMerch_BadRequestInsufficientCoins(t *testing.T) {
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

	// Buy some merch to reduce coin balance

	// Spend 500 coins from 1000
	e.GET("/api/buy/{item}", "pink-hoody").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK)

	// Spend 300 coins from 500
	e.GET("/api/buy/{item}", "hoody").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK)

	// Try to buy another merch by 500 coins and expect bad request
	e.GET("/api/buy/{item}", "pink-hoody").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusBadRequest)
}

func TestBuyMerch_Unauthorized(t *testing.T) {
	e := newTestAPI(t)

	// Try to buy merch without authorization
	e.GET("/api/buy/{item}", "hoody").
		Expect().
		Status(http.StatusUnauthorized)
}
