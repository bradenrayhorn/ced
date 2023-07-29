package ced

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestValidate(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("is required", func(t *testing.T) {
		name := Name("")
		err := name.Validate()
		is.Equal(err.Error(), ":field is required")
	})

	t.Run("cannot be more than 255 characters", func(t *testing.T) {
		name := Name(strings.Repeat("a", 256))
		err := name.Validate()
		is.Equal(err.Error(), ":field must be at most 255 characters")
	})

	t.Run("can be 255 characters", func(t *testing.T) {
		name := Name(strings.Repeat("a", 255))
		err := name.Validate()
		is.NoErr(err)
	})
}
