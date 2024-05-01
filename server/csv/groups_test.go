package csv

import (
	"strings"
	"testing"

	"github.com/bradenrayhorn/ced/server/ced"
	"github.com/matryer/is"
)

func TestParseGroupImport(t *testing.T) {

	t.Run("can parse", func(t *testing.T) {
		is := is.New(t)

		in := `
		Bob,5,"Bob,Robert"
		Greg,2,"Greg Robert"
		`

		res, err := ParseGroupImport(strings.NewReader(strings.TrimSpace(in)))
		is.NoErr(err)
		is.Equal(res, []ced.GroupImport{
			{Name: "Bob", MaxAttendees: 5, SearchHints: "Bob,Robert"},
			{Name: "Greg", MaxAttendees: 2, SearchHints: "Greg Robert"},
		})
	})

	t.Run("trims leading whitespace", func(t *testing.T) {
		is := is.New(t)

		in := `
		 Bob,5,"Bob,Robert"
		Greg , 2, "Greg Robert"
		`

		res, err := ParseGroupImport(strings.NewReader(strings.TrimSpace(in)))
		is.NoErr(err)
		is.Equal(res, []ced.GroupImport{
			{Name: "Bob", MaxAttendees: 5, SearchHints: "Bob,Robert"},
			{Name: "Greg ", MaxAttendees: 2, SearchHints: "Greg Robert"},
		})
	})

	t.Run("error if csv has too many fields", func(t *testing.T) {
		is := is.New(t)

		in := `
		Bob,5,"Bob,Robert"
		Bob,5,"Bob,Robert",57
		`

		_, err := ParseGroupImport(strings.NewReader(strings.TrimSpace(in)))
		is.Equal(err.Error(), "failed to parse csv: record on line 2: wrong number of fields")
	})

	t.Run("error if cannot parse max attendees", func(t *testing.T) {
		is := is.New(t)

		in := `
		Bob,"a string","Bob,Robert"
		`

		_, err := ParseGroupImport(strings.NewReader(strings.TrimSpace(in)))
		is.Equal(err.Error(), "failed to parse record on line 1: invalid max_attendees: a string")
	})

	t.Run("error if cannot parse because not uint8", func(t *testing.T) {
		is := is.New(t)

		in := `
		Bob,-1,"Bob,Robert"
		`

		_, err := ParseGroupImport(strings.NewReader(strings.TrimSpace(in)))
		is.Equal(err.Error(), "failed to parse record on line 1: invalid max_attendees: -1")
	})
}

func TestGroupExport(t *testing.T) {

	t.Run("can export groups", func(t *testing.T) {
		is := is.New(t)

		group1 := ced.Group{
			Name:         ced.Name("Charlie"),
			MaxAttendees: 5,
			Attendees:    3,
			HasResponded: true,
		}
		group2 := ced.Group{
			Name:         ced.Name("Evelyn"),
			MaxAttendees: 1,
			Attendees:    0,
			HasResponded: false,
		}

		var out strings.Builder
		err := GroupsExport(&out, []ced.Group{group1, group2})
		is.NoErr(err)

		res := out.String()
		expected := `
Name,Max Attendees,Attendees,Has Responded
Charlie,5,3,true
Evelyn,1,0,false
		`
		is.Equal(strings.TrimSpace(expected), strings.TrimSpace(res))
	})
}
