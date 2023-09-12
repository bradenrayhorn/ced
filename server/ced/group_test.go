package ced

import (
	"testing"

	"github.com/matryer/is"
)

func TestNewGroup(t *testing.T) {
	t.Run("can create group", func(t *testing.T) {
		is := is.New(t)
		group, err := NewGroup(Name("Bob"), 2, "Bob Lob")
		is.NoErr(err)

		is.True(!group.ID.Empty())
		is.Equal(group.Name, Name("Bob"))
		is.Equal(group.Attendees, uint8(0))
		is.Equal(group.MaxAttendees, uint8(2))
		is.Equal(group.HasResponded, false)
		is.Equal(group.SearchHints, "Bob Lob")
	})

	t.Run("validates name", func(t *testing.T) {
		is := is.New(t)
		_, err := NewGroup(Name(""), 2, "Bob Lob")
		is.True(err != nil)
		is.Equal(err.Error(), "Invalid data provided")
	})

	t.Run("trims name and search string", func(t *testing.T) {
		is := is.New(t)
		group, err := NewGroup(Name(" Bob "), 2, " Bob Lob ")
		is.NoErr(err)

		is.Equal(group.Name, Name("Bob"))
		is.Equal(group.SearchHints, "Bob Lob")
	})
}
