package auth

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func (m *manager) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			handleResponseError(w, http.StatusUnauthorized, "authorization header not found in http request")
			return
		}
		if len(tokenStr) > 7 && strings.ToUpper(tokenStr[0:6]) == "BEARER" {
			tokenStr = tokenStr[7:]
		} else {
			handleResponseError(w, http.StatusUnauthorized, "invalid authorization header in http request")
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte("secret"), nil
		})
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func handleResponseError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	render.JSON(w, nil, ErrorResponse{Error: message})
}