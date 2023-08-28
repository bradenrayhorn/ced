package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/contract"
	"github.com/bradenrayhorn/ced/http"
	"github.com/bradenrayhorn/ced/sqlite"
)

type Application struct {
	config ced.Config

	httpServer *http.Server
	pool       *sqlite.Pool
}

func NewApplication(config ced.Config) (*Application, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Application{config: config}, nil
}

func (a *Application) Start() error {
	pool, err := sqlite.CreatePool(
		context.Background(),
		fmt.Sprintf("file:%s", a.config.DbPath),
	)
	if err != nil {
		return fmt.Errorf("failed to create sqlite pool: %w", err)
	}
	a.pool = pool

	groupRepository := sqlite.NewGroupRepository(pool)

	httpServer := http.NewServer(
		contract.NewGroupContract(groupRepository),
	)

	if err := httpServer.Open(":" + a.config.HttpPort); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}
	a.httpServer = httpServer

	return nil
}

func (a *Application) Stop() error {
	var httpError error
	var poolError error
	if a.httpServer != nil {
		httpError = a.httpServer.Close()
	}

	if a.pool != nil {
		poolError = a.pool.Close(context.Background())
	}

	return errors.Join(httpError, poolError)
}
