package sqlite

import (
	"context"
	"strings"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/sqlite/mapper"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var _ ced.IndividualRespository = (*individualRepository)(nil)

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

func (r *individualRepository) Update(ctx context.Context, individual ced.Individual) error {
	conn, done := r.pool.Conn(ctx)
	defer done()

	query := `UPDATE individuals
		SET group_id = ?,
			name = ?,
			response = ?,
			has_responded = ?
		WHERE id = ?;`

	return sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Args: []interface{}{
			individual.GroupID.String(),
			string(individual.Name),
			individual.Response,
			individual.HasResponded,
			individual.ID.String(),
		},
	})
}

func (r *individualRepository) GetByGroup(ctx context.Context, groupID ced.ID) ([]ced.Individual, error) {
	conn, done := r.pool.Conn(ctx)
	defer done()

	query := `
	SELECT
		id,
		group_id,
		name,
		response,
		has_responded
	FROM individuals WHERE group_id = ?
	ORDER BY id desc;`

	individuals := []ced.Individual{}
	err := sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Args: []interface{}{groupID.String()},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			mapped, err := mapper.Individual(stmt)
			if err != nil {
				return err
			}
			individuals = append(individuals, mapped)

			return nil
		},
	})

	return individuals, err
}

func (r *individualRepository) SearchByName(ctx context.Context, search string) (map[ced.ID][]ced.Individual, error) {
	conn, done := r.pool.Conn(ctx)
	defer done()

	query := `
	SELECT
		individuals.id,
		individuals.group_id,
		individuals.name,
		individuals.response,
		individuals.has_responded
	FROM individuals
	JOIN individuals search ON search.group_id = individuals.group_id
		AND TRIM(LOWER(search.name)) = ?
	ORDER BY individuals.id desc;`

	grouped := map[ced.ID][]ced.Individual{}
	err := sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Args: []interface{}{
			strings.TrimSpace(strings.ToLower(search)),
		},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			mapped, err := mapper.Individual(stmt)
			if err != nil {
				return err
			}

			key := mapped.GroupID
			if _, ok := grouped[key]; !ok {
				grouped[key] = []ced.Individual{}
			}
			grouped[key] = append(grouped[key], mapped)

			return nil
		},
	})

	return grouped, err
}
