package sqlite

import (
	"context"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/sqlite/mapper"
	"zombiezen.com/go/sqlite"
)

var _ ced.GroupRespository = (*groupRepository)(nil)

type groupRepository struct {
	pool *Pool
}

func NewGroupRepository(pool *Pool) *groupRepository {
	return &groupRepository{pool}
}

func (r *groupRepository) Create(ctx context.Context, group ced.Group) error {
	query := `INSERT INTO groups (id) VALUES (?);`

	return execute(ctx, r.pool, query, []any{group.ID.String()})
}

func (r *groupRepository) Get(ctx context.Context, id ced.ID) (ced.Group, error) {
	query := `SELECT id FROM groups WHERE id = ?;`

	return mustFindResult(selectOne(ctx, r.pool, query,
		[]any{id.String()},
		func(stmt *sqlite.Stmt) (ced.Group, error) {
			return mapper.Group(stmt)
		},
	))
}
