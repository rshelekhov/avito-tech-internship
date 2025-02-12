package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func handleValidationErrors(w http.ResponseWriter, r *http.Request, err error, log *slog.Logger) {
	var validationErrors validator.ValidationErrors
	ok := errors.As(err, &validationErrors)
	if !ok {
		// If this is not a validation error, return a general error
		err = fmt.Errorf("validation error: %w", err)
		log.Error(err.Error())

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: err.Error()})
		return
	}

	err = fmt.Errorf("validation error: %s", processValidationErrors(validationErrors))
	log.Error(err.Error())

	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, ErrorResponse{Error: err.Error()})
}

func processValidationErrors(errors validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errors {
		switch err.Tag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' failed on %s validation", err.Field(), err.Tag()))
		}
	}

	return strings.Join(errMsgs, "; ")
}
