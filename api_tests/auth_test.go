package api_tests

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rshelekhov/merch-store/internal/controller/http/v1/handler"
	"github.com/stretchr/testify/require"
)

func TestRegisterNewUser_HappyPath(t *testing.T) {
	e := newTestAPI(t)

	// Auth user
	resp := e.POST("/api/auth").
		WithJSON(handler.AuthRequest{
			Username: gofakeit.Username(),
			Password: randomFakePassword(),
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// Check token from response
	token := resp.Value("token").String().Raw()
	require.NotEmpty(t, token)
}

func TestAuthenticateExistingUser_HappyPath(t *testing.T) {
	e := newTestAPI(t)

	username := gofakeit.Username()
	password := randomFakePassword()

	// Register new user
	e.POST("/api/auth").
		WithJSON(handler.AuthRequest{
			Username: username,
			Password: password,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().NotEmpty()

	// Auth existing user
	resp := e.POST("/api/auth").
		WithJSON(handler.AuthRequest{
			Username: username,
			Password: password,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// Check token from response
	token := resp.Value("token").String().Raw()
	require.NotEmpty(t, token)
}

func TestAuthUser_BadRequest(t *testing.T) {
	e := newTestAPI(t)

	tests := []struct {
		name     string
		username string
		password string
	}{
		{
			name:     "Empty username",
			username: "",
			password: randomFakePassword(),
		},
		{
			name:     "Empty password",
			username: gofakeit.Username(),
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e.POST("/api/auth").
				WithJSON(handler.AuthRequest{
					Username: tt.username,
					Password: tt.password,
				}).
				Expect().
				Status(http.StatusBadRequest)
		})
	}
}
