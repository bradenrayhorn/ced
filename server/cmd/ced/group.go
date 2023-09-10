package main

import (
	"context"

	"github.com/bradenrayhorn/ced/server/ced"
)

type GroupCreateCmd struct {
	Name         string `required:"" help:"Name of the group."`
	MaxAttendees uint8  `required:"" help:"Max number of guests in the group."`
	SearchHints  string `help:"Comma separated list of people in the group."`
}

func (r *GroupCreateCmd) Run(ctx *CmdContext) error {
	_, err := ctx.pool.groupContract.Create(
		context.Background(),
		ced.Name(r.Name),
		r.MaxAttendees,
		r.SearchHints,
	)
	if err != nil {
		return err
	}

	return nil
}
