package update

import (
	"bybit-parser/internal/lib/api/response"
	"bybit-parser/internal/lib/logger/sl"
	"bybit-parser/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	NewUrl string `json:"newurl" validate:"required,url"`
	Alias  string `json:"alias,omitempty" validate:"required"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLUpdate interface {
	UpdateURL(alias string, newUrl string) error
}

func New(log *slog.Logger, urlUpdate URLUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.update.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("failed to validate request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		err = urlUpdate.UpdateURL(req.Alias, req.NewUrl)

		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.NewUrl))
			render.JSON(w, r, response.Error("url already exists"))
			return
		}

		if err != nil {
			log.Error("failed to update url", sl.Err(err))
			render.JSON(w, r, response.Error("failed to update url"))
			return
		}

		responseOK(w, r, req.Alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    alias,
	})
}
