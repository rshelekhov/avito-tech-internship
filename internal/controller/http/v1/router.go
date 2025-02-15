package v1

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rshelekhov/avito-tech-internship/internal/lib/middleware/jwt"
)

type Router struct {
	log          *slog.Logger
	jwtMgr       jwt.Manager
	authHandler  AuthHandler
	coinsHandler CoinsHandler
}

type (
	AuthHandler interface {
		Auth() http.HandlerFunc
	}

	CoinsHandler interface {
		GetInfo() http.HandlerFunc
		SendCoin() http.HandlerFunc
		BuyMerch() http.HandlerFunc
	}
)

func NewRouter(
	log *slog.Logger,
	jwtMgr jwt.Manager,
	authHandler AuthHandler,
	coinsHandler CoinsHandler,
) *chi.Mux {
	ar := &Router{
		log:          log,
		jwtMgr:       jwtMgr,
		authHandler:  authHandler,
		coinsHandler: coinsHandler,
	}

	return ar.initRoutes()
}
