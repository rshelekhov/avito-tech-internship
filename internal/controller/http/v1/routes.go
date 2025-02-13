package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	mwlogger "github.com/rshelekhov/avito-tech-internship/internal/lib/middleware/logger"
)

func (ar *Router) initRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// By default, middleware.logger uses its own internal logger,
	// which should be overridden to use ours. Otherwise, problems
	// may arise - for example, with log collection. We can use
	// our own middleware to log requests:
	r.Use(mwlogger.New(ar.log))

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/health", HealthCheck())

	r.Post("/api/auth", ar.authHandler.Auth())

	r.Group(func(r chi.Router) {
		r.Use(ar.jwtMgr.HTTPMiddleware)

		r.Route("/api", func(r chi.Router) {
			r.Get("/user", ar.coinsHandler.GetInfo())
			r.Post("sendCoin", ar.coinsHandler.SendCoin())
			r.Get("buy/{item}", ar.coinsHandler.BuyMerch())
		})
	})

	return r
}
