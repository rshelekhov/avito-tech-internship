package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/rshelekhov/avito-tech-internship/internal/domain"
	"github.com/rshelekhov/avito-tech-internship/internal/domain/entity"
)

type CoinsHandler struct {
	log      *slog.Logger
	validate *validator.Validate
	usecase  CoinsUsecase
}

type CoinsUsecase interface {
	GetUserInfo(ctx context.Context) (entity.UserInfo, error)
	SendCoin(ctx context.Context, toUser string, amount int) error
	BuyMerch(ctx context.Context, itemName string) error
}

func NewCoinsHandler(log *slog.Logger, validate *validator.Validate, usecase CoinsUsecase) *CoinsHandler {
	return &CoinsHandler{
		log:      log,
		validate: validate,
		usecase:  usecase,
	}
}

type InfoResponse struct {
	Coins       int                `json:"coins"`
	Inventory   []entity.Item      `json:"inventory"`
	CoinHistory entity.CoinHistory `json:"coinHistory"`
}

func (h *CoinsHandler) GetInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.GetInfo"

		ctx := r.Context()

		userInfo, err := h.usecase.GetUserInfo(ctx)

		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			err = fmt.Errorf("%s: failed to get user info: %w", op, err)
			handleNotFoundError(w, r, err, h.log)
			return
		case err != nil:
			err = fmt.Errorf("%s: failed to get user info: %w", op, err)

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, InfoResponse{
			Coins:       userInfo.Coins,
			Inventory:   userInfo.Inventory,
			CoinHistory: userInfo.CoinHistory,
		})
	}
}

type SendCoinRequest struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int    `json:"amount" validate:"required"`
}

func (h *CoinsHandler) SendCoin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.SendCoin"

		log := h.log.With(slog.String("op", op))

		request := &SendCoinRequest{}
		if err := render.Decode(r, request); err != nil {
			err = fmt.Errorf("%s: failed to decode request: %w", op, err)
			log.Error(err.Error())

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}

		if err := h.validate.Struct(request); err != nil {
			handleValidationErrors(w, r, err, log)
			return
		}

		ctx := r.Context()

		err := h.usecase.SendCoin(ctx, request.ToUser, request.Amount)

		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			err = fmt.Errorf("%s: failed to send coin: %w", op, err)
			handleNotFoundError(w, r, err, log)
			return
		case errors.Is(err, domain.ErrBadRequest):
			err = fmt.Errorf("%s: failed to send coin: %w", op, err)
			handleBadRequestError(w, r, err, log)
			return
		case err != nil:
			err = fmt.Errorf("%s: failed to send coin: %w", op, err)
			handleInternalError(w, r, err, log)
			return
		}

		log.Info("coin sent", "toUser", request.ToUser, "amount", request.Amount)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, nil)
	}
}

func (h *CoinsHandler) BuyMerch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.BuyMerch"

		log := h.log.With(slog.String("op", op))

		itemName := chi.URLParam(r, "item")
		if itemName == "" {
			err := fmt.Errorf("%s: item name is empty in request", op)
			handleBadRequestError(w, r, err, log)
			return
		}

		ctx := r.Context()

		err := h.usecase.BuyMerch(ctx, itemName)

		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			err = fmt.Errorf("%s: failed to buy merch: %w", op, err)
			handleNotFoundError(w, r, err, log)
			return
		case errors.Is(err, domain.ErrBadRequest):
			err = fmt.Errorf("%s: failed to buy merch: %w", op, err)
			handleBadRequestError(w, r, err, log)
			return
		case err != nil:
			err = fmt.Errorf("%s: failed to buy merch: %w", op, err)
			handleInternalError(w, r, err, log)
			return
		}

		log.Info("merch bought", "item", itemName)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, nil)
	}
}
