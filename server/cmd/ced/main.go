package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/bradenrayhorn/ced/server/env"
)

type CmdContext struct {
	in   io.Reader
	out  io.Writer
	pool *cmdPool
}

var cli struct {
	Group struct {
		Create GroupCreateCmd `cmd:"" help:"Creates a group invitation."`
		Import GroupImportCmd `cmd:"" help:"Imports groups from a csv. Pass csv as stdin."`
		Export GroupExportCmd `cmd:"" help:"Exports all groups to a csv."`
	} `cmd:""`
}

func main() {
	ctx := kong.Parse(&cli, kong.Bind(os.Stdin))

	dbPath := env.GetDbPath()
	if strings.TrimSpace(dbPath) == "" {
		ctx.FatalIfErrorf(fmt.Errorf("DB_PATH env is required"))
		return
	}

	pool, err := newCmdPool(fmt.Sprintf("file:%s", dbPath))
	if err != nil {
		ctx.FatalIfErrorf(err)
		return
	}
	defer pool.close(os.Stdout)

	err = ctx.Run(&CmdContext{in: os.Stdin, out: os.Stdout, pool: pool})
	ctx.FatalIfErrorf(err)
}
