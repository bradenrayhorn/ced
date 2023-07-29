package sqlite

import (
	"context"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/sqlite/mapper"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type individualRepository struct {
	pool *Pool
}

func NewIndividualRepository(pool *Pool) *individualRepository {
	return &individualRepository{pool}
}

func (r *individualRepository) Create(ctx context.Context, individual ced.Individual) error {
	conn, done := r.pool.Conn(ctx)
	defer done()

	query := `INSERT INTO individuals (
		id,
		group_id,
		name,
		response,
		has_responded
	) VALUES (?, ?, ?, ?, ?);`

	return sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Args: []interface{}{
			individual.ID.String(),
			individual.GroupID.String(),
			string(individual.Name),
			individual.Response,
			individual.HasResponded,
		},
	})
}

func (r *individualRepository) Get(ctx context.Context, id ced.ID) (ced.Individual, error) {
	conn, done := r.pool.Conn(ctx)
	defer done()

	var individual ced.Individual
	query := `
	SELECT
		id,
		group_id,
		name,
		response,
		has_responded
	FROM individuals WHERE id = ?;`

	err := sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Args: []interface{}{id.String()},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			mapped, err := mapper.Individual(stmt)
			if err != nil {
				return err
			}
			individual = mapped

			return nil
		},
	})

	return mustFindResult(individual, err)
}
