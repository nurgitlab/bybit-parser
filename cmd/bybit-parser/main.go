package main

import (
	"bybit-parser/internal/config"
	read_orderbook "bybit-parser/internal/http-server/handlers/bybit/read-orderbook"
	"bybit-parser/internal/http-server/handlers/url/delete"
	"bybit-parser/internal/http-server/handlers/url/read"
	readall "bybit-parser/internal/http-server/handlers/url/read-all"
	"bybit-parser/internal/http-server/handlers/url/save"
	"bybit-parser/internal/http-server/handlers/url/update"
	mwLogger "bybit-parser/internal/http-server/middleware/logger"
	"bybit-parser/internal/lib/logger/handlers/slogpretty"
	"bybit-parser/internal/lib/logger/sl"
	"bybit-parser/internal/storage/psql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoadConfig()
	log := setupLogger(cfg.Env)

	log.Info("starting bybit-parser", slog.String("env", cfg.Env))
	log.Debug("debug logging enabled")

	storage, err := psql.New(fmt.Sprintf(
		"postgres://%v:%v@localhost:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	))

	if err != nil {
		log.Error("Database connect Error", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID) //добавляет Request Id для трейсинга
	router.Use(middleware.RealIP)

	router.Use(middleware.Logger) //логирование всех входящих запросов
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Delete("/url", delete.Url(log, storage))
	router.Patch("/url", update.New(log, storage))
	router.Get("/url", read.New(log, storage))
	router.Get("/all-urls", readall.New(log, storage))

	router.Get("/bybit/orders", read_orderbook.New(log, storage))

	log.Info("Сервер запущен", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	log.Error("failed to start server")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
