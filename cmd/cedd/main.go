package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/env"
)

func main() {
	config := env.LoadConfig()

	setupLogger(config)

	application, err := NewApplication(config)
	if err != nil {
		slog.Error("failed to construct", "error", err.Error())
		os.Exit(1)
		return
	}
	defer func() {
		if err := application.Stop(); err != nil {
			slog.Error("failed to shutdown", "error", err.Error())
		}
	}()

	if err := application.Start(); err != nil {
		slog.Error("failed to start cedd", "error", err.Error())
		os.Exit(1)
		return
	}

	slog.Info(fmt.Sprintf("listening on port %s", config.HttpPort))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	slog.Info("shutting down cedd...")
}

func setupLogger(config ced.Config) {
	if !config.PrettyLog {
		h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})
		slog.SetDefault(slog.New(h))
	}
}
