package main

import (
	"context"
	"fmt"

	"github.com/bradenrayhorn/ced/server/ced"
	"github.com/bradenrayhorn/ced/server/csv"
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

type GroupImportCmd struct {
}

func (r *GroupImportCmd) Run(ctx *CmdContext) error {
	records, err := csv.ParseGroupImport(ctx.in)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(ctx.out, "importing %d groups...", len(records)); err != nil {
		return err
	}

	return ctx.pool.groupContract.Import(context.Background(), records)
}

type GroupExportCmd struct {
}

func (r *GroupExportCmd) Run(ctx *CmdContext) error {
	groups, err := ctx.pool.groupContract.Export(context.Background())
	if err != nil {
		return err
	}

	return csv.GroupsExport(ctx.out, groups)
}
