package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/contract"
	"github.com/bradenrayhorn/ced/env"
	"github.com/bradenrayhorn/ced/http"
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

	httpServer := http.NewServer(
		contract.NewIndividualContract(sqlite.NewIndividualRepository(pool)),
	)

	if err := httpServer.Open(":" + config.HttpPort); err != nil {
		slog.Error("failed to start http server", "error", err)
		os.Exit(1)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	slog.Info("shutting down cedd...")

	if err := httpServer.Close(); err != nil {
		slog.Error("failed to stop http server", "error", err)
	}
}

func setupLogger(config ced.Config) {
	if !config.PrettyLog {
		h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})
		slog.SetDefault(slog.New(h))
	}
}
