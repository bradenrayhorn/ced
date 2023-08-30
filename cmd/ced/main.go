package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/bradenrayhorn/ced/env"
)

type CmdContext struct {
	in     io.Reader
	out    io.Writer
	dbPath string
}

var cli struct {
	Group struct {
		Create GroupCreateCmd `cmd:"" help:"Creates a group invitation."`
	} `cmd:""`
}

func main() {
	ctx := kong.Parse(&cli, kong.Bind(os.Stdin))

	dbPath := env.GetDbPath()
	if strings.TrimSpace(dbPath) == "" {
		ctx.FatalIfErrorf(fmt.Errorf("DB_PATH env is required"))
		return
	}

	err := ctx.Run(&CmdContext{in: os.Stdin, out: os.Stdout, dbPath: dbPath})
	ctx.FatalIfErrorf(err)
}
