package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/env"
	"github.com/bradenrayhorn/ced/sqlite"
)

func main() {
	config := env.LoadConfig()

	setupLogger(config)

	slog.Info("starting cedd...")

	pool, err := sqlite.CreatePool(context.Background(), config)
	if err != nil {
		slog.Error("failed to create sqlite pool", "error", err)
		os.Exit(1)
		return
	}
	defer func() { _ = pool.Close(context.Background()) }()
}

func setupLogger(config ced.Config) {
	if !config.PrettyLog {
		h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})
		slog.SetDefault(slog.New(h))
	}
}
