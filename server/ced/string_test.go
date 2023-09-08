package ced

import (
	"testing"

	"github.com/matryer/is"
)

func TestEmpty(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("empty", func(t *testing.T) {
		s := ValidatableString("")
		is.Equal(s.Empty(), true)
	})

	t.Run("ignores whitespace", func(t *testing.T) {
		s := ValidatableString("  ")
		is.Equal(s.Empty(), true)
	})

	t.Run("not empty", func(t *testing.T) {
		s := ValidatableString(" x ")
		is.Equal(s.Empty(), false)
	})
}

func TestLength(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("empty", func(t *testing.T) {
		s := ValidatableString("")
		is.Equal(s.Length(), 0)
	})

	t.Run("includes whitespace", func(t *testing.T) {
		s := ValidatableString("  ")
		is.Equal(s.Length(), 2)
	})

	t.Run("counts alpha characters", func(t *testing.T) {
		s := ValidatableString("a")
		is.Equal(s.Length(), 1)
	})
}
