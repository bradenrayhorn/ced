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
		ID: id,
	}, nil
}
