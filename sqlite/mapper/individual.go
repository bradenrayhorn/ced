package mapper

import (
	"github.com/bradenrayhorn/ced/ced"
	"zombiezen.com/go/sqlite"
)

func Individual(stmt *sqlite.Stmt) (ced.Individual, error) {
	id, err := ced.IDFromString(stmt.GetText("id"))
	if err != nil {
		return ced.Individual{}, nil
	}

	groupID, err := ced.IDFromString(stmt.GetText("group_id"))
	if err != nil {
		return ced.Individual{}, nil
	}

	return ced.Individual{
		ID:           id,
		GroupID:      groupID,
		Name:         ced.Name(stmt.GetText("name")),
		Response:     stmt.GetBool("response"),
		HasResponded: stmt.GetBool("has_responded"),
	}, nil
}
