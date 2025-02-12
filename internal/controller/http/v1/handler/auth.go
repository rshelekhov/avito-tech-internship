package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/rshelekhov/avito-tech-internship/internal/domain/entity"
)

type AuthHandler struct {
	log      *slog.Logger
	validate *validator.Validate
	usecase  AuthUsecase
}

type AuthUsecase interface {
	Authenticate(ctx context.Context, credentials entity.UserCredentials) (string, error)
}

func NewAuthHandler(log *slog.Logger, validate *validator.Validate, usecase AuthUsecase) *AuthHandler {
	return &AuthHandler{
		log:      log,
		validate: validate,
		usecase:  usecase,
	}
}

type AuthRequest struct {
	username string `json:"username" validate:"required"`
	password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	token string `json:"token"`
}

func (h *AuthHandler) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.Auth"

		log := h.log.With(slog.String("op", op))

		request := &AuthRequest{}
		if err := render.Decode(r, request); err != nil {
			err = fmt.Errorf("failed to decode request: %w", err)
			h.log.Error(err.Error())

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}

		if err := h.validate.Struct(request); err != nil {
			handleValidationErrors(w, r, err, log)
			return
		}

		ctx := r.Context()
		user := toUserCredentials(request)

		token, err := h.usecase.Authenticate(ctx, user)
		if err != nil {
			err = fmt.Errorf("failed to authenticate user: %w", err)
			handleInternalError(w, r, err, log)
			return
		}

		log.Info("user authenticated", slog.String("username", user.Username))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, AuthResponse{token: token})
	}
}
