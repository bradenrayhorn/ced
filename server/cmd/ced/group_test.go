package main

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/bradenrayhorn/ced/server/ced"
	"github.com/bradenrayhorn/ced/server/internal/testutils"
	"github.com/matryer/is"
)

func TestGroupCreate(t *testing.T) {
	is := is.New(t)
	pool, err := newCmdPool(testutils.InMemoryPoolPath)
	is.NoErr(err)
	defer pool.close(os.Stdout)

	cmd := GroupCreateCmd{Name: "Max Hoover & family", MaxAttendees: 2, SearchHints: "Max Hoover, Tod Frog"}
	err = cmd.Run(&CmdContext{pool: pool})
	is.NoErr(err)

	res, err := pool.groupContract.Search(context.Background(), ced.ReqContext{}, "Max Hoover")
	is.NoErr(err)
	is.Equal(1, len(res))

	group := res[0]
	is.Equal(string(group.Name), "Max Hoover & family")
	is.Equal(group.MaxAttendees, uint8(2))
	is.Equal(group.SearchHints, "Max Hoover, Tod Frog")
}

func TestGroupImport(t *testing.T) {
	is := is.New(t)
	pool, err := newCmdPool(testutils.InMemoryPoolPath)
	is.NoErr(err)
	defer pool.close(os.Stdout)

	in := bytes.NewBufferString(strings.TrimSpace(`
	Bob,5,"Bob Lob"
	`))
	out := bytes.NewBufferString("")

	cmd := GroupImportCmd{}
	err = cmd.Run(&CmdContext{in, out, pool})
	is.NoErr(err)

	output := out.String()
	is.Equal(output, "importing 1 groups...")

	res, err := pool.groupContract.Search(context.Background(), ced.ReqContext{}, "Bob Lob")
	is.NoErr(err)
	is.Equal(1, len(res))

	group := res[0]
	is.Equal(string(group.Name), "Bob")
	is.Equal(group.MaxAttendees, uint8(5))
	is.Equal(group.SearchHints, "Bob Lob")
}

func TestGroupExport(t *testing.T) {
	is := is.New(t)
	pool, err := newCmdPool(testutils.InMemoryPoolPath)
	is.NoErr(err)
	defer pool.close(os.Stdout)

	_, err = pool.groupContract.Create(context.Background(), ced.Name("Bob Lob"), 2, "")
	is.NoErr(err)

	in := bytes.NewBufferString("")
	out := bytes.NewBufferString("")

	cmd := GroupExportCmd{}
	err = cmd.Run(&CmdContext{in, out, pool})
	is.NoErr(err)

	is.Equal(strings.TrimSpace(out.String()), strings.TrimSpace(`
Name,Max Attendees,Attendees,Has Responded
Bob Lob,2,0,false
	`))
}
