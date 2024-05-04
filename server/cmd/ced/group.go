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

type GroupUpdateCmd struct {
	CurrentName string `arg:"" help:"Current name of the group."`

	Name         *string `optional:"" help:"New name of the group."`
	MaxAttendees *uint8  `optional:"" help:"Max number of guests in the group."`
	Attendees    *uint8  `optional:"" help:"Number of attendees in the group. Changing this will mark the group as reserved."`
	SearchHints  *string `optional:"" help:"Comma separated list of people in the group."`
}

func (r *GroupUpdateCmd) Run(ctx *CmdContext) error {
	group, err := ctx.pool.groupContract.FindOne(context.Background(), r.CurrentName)
	if err != nil {
		return err
	}

	proceed, err := askYesNo(fmt.Sprintf(
		"Are you sure you wish to update '%s' (%d/%d) (responded: %t)? [y/n]",
		group.Name,
		group.Attendees,
		group.MaxAttendees,
		group.HasResponded,
	), ctx.in, ctx.out)
	if err != nil {
		return err
	}

	if proceed {
		update := ced.GroupUpdate{
			ID:           group.ID,
			MaxAttendees: r.MaxAttendees,
			Attendees:    r.Attendees,
			SearchHints:  r.SearchHints,
		}
		if r.Name != nil {
			name := ced.Name(*r.Name)
			update.Name = &name
		}

		if err := ctx.pool.groupContract.Update(context.Background(), update); err != nil {
			return err
		}

		if _, err := fmt.Fprint(ctx.out, "Group updated.\n"); err != nil {
			return err
		}
	}

	return nil
}

type GroupDeleteCmd struct {
	Name string `arg:"" help:"Name of the group."`
}

func (r *GroupDeleteCmd) Run(ctx *CmdContext) error {
	group, err := ctx.pool.groupContract.FindOne(context.Background(), r.Name)
	if err != nil {
		return err
	}

	proceed, err := askYesNo(fmt.Sprintf(
		"Are you sure you wish to delete '%s' (%d/%d) (responded: %t)? [y/n]",
		group.Name,
		group.Attendees,
		group.MaxAttendees,
		group.HasResponded,
	), ctx.in, ctx.out)
	if err != nil {
		return err
	}

	if proceed {
		if err := ctx.pool.groupContract.Delete(context.Background(), group.ID); err != nil {
			return err
		}

		if _, err := fmt.Fprint(ctx.out, "Group deleted.\n"); err != nil {
			return err
		}
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
