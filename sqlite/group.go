package sqlite

import (
	"context"
	"strings"

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
	query := `INSERT INTO groups
				(id,name,attendees,max_attendees,has_responded)
				VALUES (?,?,?,?,?)
	;`

	return execute(ctx, r.pool, query, []any{
		group.ID.String(),
		string(group.Name),
		group.Attendees,
		group.MaxAttendees,
		group.HasResponded,
	})
}

func (r *groupRepository) Update(ctx context.Context, group ced.Group) error {
	query := `
	UPDATE groups
	SET
		name = ?,
		attendees = ?,
		max_attendees = ?,
		has_responded = ?
	WHERE id = ?
	;`

	return execute(ctx, r.pool, query, []any{
		string(group.Name),
		group.Attendees,
		group.MaxAttendees,
		group.HasResponded,
		group.ID.String(),
	})
}

func (r *groupRepository) Get(ctx context.Context, id ced.ID) (ced.Group, error) {
	query := `SELECT * FROM groups WHERE id = ?;`

	return mustFindResult(selectOne(ctx, r.pool, query,
		[]any{id.String()},
		func(stmt *sqlite.Stmt) (ced.Group, error) {
			return mapper.Group(stmt)
		},
	))
}

func (r *groupRepository) SearchByName(ctx context.Context, name string) ([]ced.Group, error) {
	query := `SELECT * FROM groups WHERE TRIM(LOWER(name)) like ?;`

	return selectList(ctx, r.pool, query,
		[]any{"%" + strings.TrimSpace(strings.ToLower(name)) + "%"},

		func(stmt *sqlite.Stmt) (ced.Group, error) {
			return mapper.Group(stmt)
		},
	)
}
