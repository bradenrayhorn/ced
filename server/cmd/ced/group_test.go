package main

import (
	"context"
	"os"
	"testing"

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

	res, err := pool.groupContract.Search(context.Background(), "Max Hoover")
	is.NoErr(err)
	is.Equal(1, len(res))

	group := res[0]
	is.Equal(string(group.Name), "Max Hoover & family")
	is.Equal(group.MaxAttendees, uint8(2))
	is.Equal(group.SearchHints, "Max Hoover, Tod Frog")
}
