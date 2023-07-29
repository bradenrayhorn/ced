package ced

import (
	"testing"

	"github.com/matryer/is"
)

func TestFromString(t *testing.T) {
	is := is.New(t)

	t.Run("parses id", func(t *testing.T) {
		id, err := IDFromString("0ujsswThIGTUYm2K8FjOOfXtY1K")
		is.NoErr(err)
		is.Equal(id.String(), "0ujsswThIGTUYm2K8FjOOfXtY1K")
	})

	t.Run("handles parse error", func(t *testing.T) {
		_, err := IDFromString("jibberish")
		is.True(err != nil)
	})

	t.Run("handles empty string", func(t *testing.T) {
		_, err := IDFromString("jibberish")
		is.True(err != nil)
	})
}
