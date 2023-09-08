package main

import (
	"context"

	"github.com/bradenrayhorn/ced/server/ced"
)

type GroupCreateCmd struct {
	Name         string `required:"" help:"Name of the group."`
	MaxAttendees uint8  `required:"" help:"Max number of guests in the group."`
}

func (r *GroupCreateCmd) Run(ctx *CmdContext) error {
	pool, err := newCmdPool(ctx)
	if err != nil {
		return err
	}
	defer pool.close(ctx)

	_, err = pool.groupContract.Create(
		context.Background(),
		ced.Name(r.Name),
		r.MaxAttendees,
	)
	if err != nil {
		return err
	}

	return nil
}
