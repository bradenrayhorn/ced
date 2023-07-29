package main

import (
	"context"
	"os"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/env"
	"github.com/bradenrayhorn/ced/sqlite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	config := env.LoadConfig()

	setupLogger(config)

	log.Info().Msg("starting cedd...")

	pool, err := sqlite.CreatePool(context.Background(), config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create sqlite pool")
		return
	}
	defer func() { _ = pool.Close(context.Background()) }()
}

func setupLogger(config ced.Config) {
	if config.PrettyLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
