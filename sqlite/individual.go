package sqlite

import (
	"context"
	"strings"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/sqlite/mapper"
	"zombiezen.com/go/sqlite"
)

var _ ced.IndividualRespository = (*individualRepository)(nil)

type individualRepository struct {
	pool *Pool
}

func NewIndividualRepository(pool *Pool) *individualRepository {
	return &individualRepository{pool}
}

func (r *individualRepository) Create(ctx context.Context, individual ced.Individual) error {
	query := `INSERT INTO individuals (
		id,
		group_id,
		name,
		response,
		has_responded
	) VALUES (?, ?, ?, ?, ?);`

	return execute(ctx, r.pool, query, []any{
		individual.ID.String(),
		individual.GroupID.String(),
		string(individual.Name),
		individual.Response,
		individual.HasResponded,
	})
}

func (r *individualRepository) Get(ctx context.Context, id ced.ID) (ced.Individual, error) {
	query := `
	SELECT
		id,
		group_id,
		name,
		response,
		has_responded
	FROM individuals WHERE id = ?;`

	return mustFindResult(
		selectOne(ctx, r.pool, query,
			[]any{id.String()},
			func(stmt *sqlite.Stmt) (ced.Individual, error) {
				return mapper.Individual(stmt)
			}),
	)
}

func (r *individualRepository) Update(ctx context.Context, individual ced.Individual) error {
	query := `UPDATE individuals
		SET group_id = ?,
			name = ?,
			response = ?,
			has_responded = ?
		WHERE id = ?;`

	return execute(ctx, r.pool, query, []any{
		individual.GroupID.String(),
		string(individual.Name),
		individual.Response,
		individual.HasResponded,
		individual.ID.String(),
	})
}

func (r *individualRepository) GetByGroup(ctx context.Context, groupID ced.ID) ([]ced.Individual, error) {
	query := `
	SELECT
		id,
		group_id,
		name,
		response,
		has_responded
	FROM individuals WHERE group_id = ?
	ORDER BY id desc;`

	return selectList(ctx, r.pool, query,
		[]any{groupID.String()},
		func(stmt *sqlite.Stmt) (ced.Individual, error) {
			return mapper.Individual(stmt)
		},
	)
}

func (r *individualRepository) SearchByName(ctx context.Context, search string) (map[ced.ID][]ced.Individual, error) {
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
	individuals, err := selectList(ctx, r.pool, query,
		[]any{strings.TrimSpace(strings.ToLower(search))},
		func(stmt *sqlite.Stmt) (ced.Individual, error) {
			return mapper.Individual(stmt)
		},
	)
	if err != nil {
		return grouped, err
	}

	for _, individual := range individuals {
		key := individual.GroupID
		if _, ok := grouped[key]; !ok {
			grouped[key] = []ced.Individual{}
		}
		grouped[key] = append(grouped[key], individual)
	}
	return grouped, nil
}
