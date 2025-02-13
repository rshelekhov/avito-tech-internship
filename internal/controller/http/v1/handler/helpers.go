package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

func handleInternalError(w http.ResponseWriter, r *http.Request, err error, log *slog.Logger) {
	log.Error(err.Error())

	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, ErrorResponse{Error: err.Error()})
}

func handleBadRequestError(w http.ResponseWriter, r *http.Request, err error, log *slog.Logger) {
	log.Error(err.Error())

	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, ErrorResponse{Error: err.Error()})
}

func handleNotFoundError(w http.ResponseWriter, r *http.Request, err error, log *slog.Logger) {
	log.Error(err.Error())

	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, ErrorResponse{Error: err.Error()})
}
