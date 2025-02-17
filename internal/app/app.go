package app

import (
	"fmt"
	"log/slog"

	"github.com/rshelekhov/merch-store/internal/infrastructure/storage/transaction"

	validator "github.com/go-playground/validator/v10"
	"github.com/rshelekhov/merch-store/internal/app/http"
	"github.com/rshelekhov/merch-store/internal/config"
	"github.com/rshelekhov/merch-store/internal/config/settings"
	v1 "github.com/rshelekhov/merch-store/internal/controller/http/v1"
	"github.com/rshelekhov/merch-store/internal/controller/http/v1/handler"
	coinsService "github.com/rshelekhov/merch-store/internal/domain/service/coins"
	merchService "github.com/rshelekhov/merch-store/internal/domain/service/merch"
	"github.com/rshelekhov/merch-store/internal/domain/service/token"
	userService "github.com/rshelekhov/merch-store/internal/domain/service/user"
	"github.com/rshelekhov/merch-store/internal/domain/usecase/auth"
	"github.com/rshelekhov/merch-store/internal/domain/usecase/coins"
	"github.com/rshelekhov/merch-store/internal/infrastructure/storage"
	coinsDB "github.com/rshelekhov/merch-store/internal/infrastructure/storage/coins"
	merchDB "github.com/rshelekhov/merch-store/internal/infrastructure/storage/merch"
	userDB "github.com/rshelekhov/merch-store/internal/infrastructure/storage/user"
	"github.com/rshelekhov/merch-store/internal/lib/middleware/jwt"
)

type App struct {
	HTTPServer *http.App
	dbConn     *storage.DBConnection
}

func New(cfg *config.ServerSettings, log *slog.Logger) (*App, error) {
	dbConn, err := newDBConnection(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to init database connection: %w", err)
	}

	// Init transaction manager
	txMgr := transaction.NewManager(dbConn)

	// Init storages
	coinsStorage := coinsDB.NewStorage(dbConn.Postgres.Pool, txMgr)
	merchStorage := merchDB.NewStorage(dbConn.Postgres.Pool, txMgr)
	userStorage := userDB.NewStorage(dbConn.Postgres.Pool)

	// Init managers
	coinsMgr := coinsService.New(coinsStorage)
	merchMgr := merchService.New(merchStorage)
	userMgr := userService.New(userStorage)
	tokenService := newTokenService(cfg.JWT, cfg.PasswordHash)

	// Init usecases
	authUsecase := auth.NewUsecase(log, userMgr, tokenService, tokenService)
	coinsUsecase := coins.NewUsecase(log, tokenService, userMgr, coinsMgr, merchMgr, txMgr)

	validate := validator.New()

	// Init handlers
	authHandler := handler.NewAuthHandler(log, validate, authUsecase)
	coinsHandler := handler.NewCoinsHandler(log, validate, coinsUsecase)

	// Init managers
	jwtMgr := jwt.NewManager(cfg.JWT.Secret)

	// Init HTTP server
	router := v1.NewRouter(log, jwtMgr, authHandler, coinsHandler)
	httpServer := http.New(cfg.HTTPServer, log, router)

	return &App{
		HTTPServer: httpServer,
		dbConn:     dbConn,
	}, nil
}

func (a *App) Stop() error {
	const method = "app.Stop"

	// Shutdown HTTP server
	if err := a.HTTPServer.Stop(); err != nil {
		return fmt.Errorf("%s:failed to stop http server: %w", method, err)
	}

	// Close database connection
	a.dbConn.Close()

	return nil
}

func newDBConnection(cfg settings.Postgres) (*storage.DBConnection, error) {
	storageConfig := settings.ToStorageConfig(cfg)

	return storage.NewDBConnection(storageConfig)
}

func newTokenService(jwt settings.JWT, passwordHash settings.PasswordHash) *token.Service {
	jwtConfig := settings.ToJWTConfig(jwt)

	passwordHashConfig := settings.ToPasswordHashConfig(passwordHash)

	return token.NewService(token.Config{
		JWT:          jwtConfig,
		PasswordHash: passwordHashConfig,
	})
}
