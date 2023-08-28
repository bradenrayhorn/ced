package mapper

import (
	"github.com/bradenrayhorn/ced/ced"
	"zombiezen.com/go/sqlite"
)

func Group(stmt *sqlite.Stmt) (ced.Group, error) {
	id, err := ced.IDFromString(stmt.GetText("id"))
	if err != nil {
		return ced.Group{}, err
	}

	return ced.Group{
		ID:           id,
		Name:         ced.Name(stmt.GetText("name")),
		Attendees:    uint8(stmt.GetInt64("attendees")),
		MaxAttendees: uint8(stmt.GetInt64("max_attendees")),
		HasResponded: stmt.GetBool("has_responded"),
	}, nil
}
