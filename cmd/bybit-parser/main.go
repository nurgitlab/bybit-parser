package main

import (
	"bybit-parser/internal/config"
	"bybit-parser/internal/lig/logger/sl"
	"bybit-parser/internal/storage/psql"
	"fmt"
	"log/slog"
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

	err = storage.DeleteURL("m1")

	if err != nil {
		log.Error("Get URL Error", sl.Err(err))
	}

	//id, err := storage.SaveURL("mail1", "m11")
	//if err != nil {
	//	log.Error("Save URL Error", sl.Err(err))
	//}
	//log.Info("saved url1", id)
	//
	//_ = storage
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}

	return log
}
