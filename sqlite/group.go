package sqlite

import (
	"context"

	"github.com/bradenrayhorn/ced/ced"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type groupRepository struct {
	pool *Pool
}

func NewGroupRepository(pool *Pool) *groupRepository {
	return &groupRepository{pool}
}

func (r *groupRepository) Create(ctx context.Context, group ced.Group) error {
	conn, done := r.pool.Conn(ctx)
	defer done()

	query := `INSERT INTO groups (id) VALUES (?);`

	return sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Args: []interface{}{group.ID.String()},
	})
}

func (r *groupRepository) Get(ctx context.Context, id ced.ID) (ced.Group, error) {
	conn, done := r.pool.Conn(ctx)
	defer done()

	var group ced.Group
	query := `SELECT id FROM groups WHERE id = ?;`

	err := sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Args: []interface{}{id.String()},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			group.ID = id
			return nil
		},
	})

	return mustFindResult(group, err)
}
