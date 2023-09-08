package main

import (
	"context"
	"fmt"

	"github.com/bradenrayhorn/ced/server/ced"
	"github.com/bradenrayhorn/ced/server/contract"
	"github.com/bradenrayhorn/ced/server/sqlite"
)

type cmdPool struct {
	pool          *sqlite.Pool
	groupContract ced.GroupContract
}

func (c *cmdPool) close(ctx *CmdContext) {
	err := c.pool.Close(context.Background())
	if err != nil {
		_, _ = fmt.Fprintf(ctx.out, "failed to close pool: %v", err)
	}
}

func newCmdPool(ctx *CmdContext) (*cmdPool, error) {
	pool, err := sqlite.CreatePool(
		context.Background(),
		fmt.Sprintf("file:%s", ctx.dbPath),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlite pool: %w", err)
	}

	groupRepository := sqlite.NewGroupRepository(pool)

	return &cmdPool{
		pool:          pool,
		groupContract: contract.NewGroupContract(groupRepository),
	}, nil
}
