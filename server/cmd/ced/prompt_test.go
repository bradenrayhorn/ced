package main

import (
	"bytes"
	"testing"

	"github.com/matryer/is"
)

func TestPromptYesNo(t *testing.T) {

	t.Run("simple yes", func(t *testing.T) {
		is := is.New(t)
		in := bytes.NewBufferString("y")
		out := bytes.NewBufferString("")

		res, err := askYesNo("are you sure?", in, out)
		is.NoErr(err)
		is.Equal(true, res)

		is.Equal("are you sure?\n", out.String())
	})

	t.Run("simple no", func(t *testing.T) {
		is := is.New(t)
		in := bytes.NewBufferString("n")
		out := bytes.NewBufferString("")

		res, err := askYesNo("are you sure?", in, out)
		is.NoErr(err)
		is.Equal(false, res)

		is.Equal("are you sure?\n", out.String())
	})

	t.Run("repeats question", func(t *testing.T) {
		is := is.New(t)
		in := bytes.NewBufferString("xn")
		out := bytes.NewBufferString("")

		res, err := askYesNo("are you sure?", in, out)
		is.NoErr(err)
		is.Equal(false, res)

		is.Equal("are you sure?\nare you sure?\n", out.String())
	})
}
