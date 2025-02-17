package jwt

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rshelekhov/merch-store/internal/domain"
)

func (m *manager) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			handleResponseError(w, "authorization header not found in http request")
			return
		}
		if len(tokenStr) > 7 && strings.ToUpper(tokenStr[0:6]) == "BEARER" {
			tokenStr = tokenStr[7:]
		} else {
			handleResponseError(w, "invalid authorization header in http request")
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(m.secret), nil
		})
		if err != nil {
			handleResponseError(w, err.Error())
			return
		}

		if !token.Valid {
			handleResponseError(w, "invalid token")
			return
		}

		userID, ok := claims[domain.UserIDKey].(string)
		if !ok {
			handleResponseError(w, "invalid token claims")
			return
		}

		ctx := m.toContext(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func handleResponseError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	render.JSON(w, nil, ErrorResponse{Error: message})
}
