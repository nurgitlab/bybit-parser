package read

import (
	"bybit-parser/internal/lib/api/response"
	"bybit-parser/internal/lib/logger/sl"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	response.Response
	Alias []string `json:"alias"`
}

type URLRead interface {
	GetAlias(url string) ([]string, error)
}

func New(log *slog.Logger, urlSaver URLRead) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.read.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		urlQuery := r.URL.Query().Get("url")

		if len(urlQuery) == 0 {
			log.Error("failed to decode request")
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		fmt.Println("urlQuery", urlQuery)
		urls, err := urlSaver.GetAlias(urlQuery)

		if err != nil {
			log.Error("failed to get alias", sl.Err(err))
			render.JSON(w, r, response.Error("failed to get alias"))
			return
		}

		log.Info("urls founded", slog.Any("urls", urls))

		responseOK(w, r, urls)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, urls []string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    urls,
	})
}
