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
