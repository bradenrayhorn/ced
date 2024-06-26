package sqlite

import (
	"context"
	"slices"
	"sort"
	"strings"
	"unicode"

	"github.com/bradenrayhorn/ced/server/ced"
	"github.com/bradenrayhorn/ced/server/sqlite/mapper"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var _ ced.GroupRespository = (*groupRepository)(nil)

type groupRepository struct {
	pool *Pool
}

func NewGroupRepository(pool *Pool) *groupRepository {
	return &groupRepository{pool}
}

func (r *groupRepository) Create(ctx context.Context, group ced.Group) error {
	return r.CreateMany(ctx, []ced.Group{group})
}

func (r *groupRepository) CreateMany(ctx context.Context, groups []ced.Group) error {
	query := `INSERT INTO groups
				(id,name,attendees,max_attendees,has_responded,search_hints)
				VALUES (?,?,?,?,?,?)
	;`

	conn, done, err := r.pool.Conn(ctx)
	if err != nil {
		return err
	}
	defer done()

	insertGroups := func(conn *sqlite.Conn) (err error) {
		defer sqlitex.Save(conn)(&err)

		for _, group := range groups {
			err = executeWithConn(conn, query, []any{
				group.ID.String(),
				string(group.Name),
				group.Attendees,
				group.MaxAttendees,
				group.HasResponded,
				group.SearchHints,
			})

			if err != nil {
				break
			}
		}

		return
	}

	return insertGroups(conn)
}

func (r *groupRepository) Update(ctx context.Context, group ced.Group) error {
	query := `
	UPDATE groups
	SET
		name = ?,
		attendees = ?,
		max_attendees = ?,
		has_responded = ?,
		search_hints = ?
	WHERE id = ?
	;`

	return execute(ctx, r.pool, query, []any{
		string(group.Name),
		group.Attendees,
		group.MaxAttendees,
		group.HasResponded,
		group.SearchHints,
		group.ID.String(),
	})
}

func (r *groupRepository) Get(ctx context.Context, id ced.ID) (ced.Group, error) {
	query := `SELECT * FROM groups WHERE id = ?;`

	group, err := selectOne(ctx, r.pool, query,
		[]any{id.String()},
		func(stmt *sqlite.Stmt) (ced.Group, error) {
			return mapper.Group(stmt)
		},
	)
	return mustFindResult(group, err, id.String())
}

func (r *groupRepository) Delete(ctx context.Context, id ced.ID) error {
	query := `DELETE FROM groups WHERE id = ?;`

	return execute(ctx, r.pool, query, []any{id.String()})
}

func (r *groupRepository) GetAll(ctx context.Context) ([]ced.Group, error) {
	query := `SELECT * FROM groups;`

	groups, err := selectList(ctx, r.pool, query,
		[]any{},
		func(stmt *sqlite.Stmt) (ced.Group, error) {
			return mapper.Group(stmt)
		},
	)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *groupRepository) SearchByName(ctx context.Context, name string) ([]ced.Group, error) {
	query := `SELECT * FROM groups;`

	sanitize := func(n string) string {
		return strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, strings.ToLower(n))
	}

	// get all groups from database
	groups, err := selectList(ctx, r.pool, query,
		[]any{},
		func(stmt *sqlite.Stmt) (ced.Group, error) {
			return mapper.Group(stmt)
		},
	)
	if err != nil {
		return nil, err
	}

	// perform search
	searchText := sanitize(name)
	ranks := []struct {
		distance int
		group    ced.Group
	}{}
	for _, group := range groups {
		names := strings.Split(string(group.Name)+","+group.SearchHints, ",")
		for _, n := range names {
			n = sanitize(n)
			rank := fuzzy.LevenshteinDistance(n, searchText)
			if n != "" && rank <= 3 {
				ranks = append(ranks, struct {
					distance int
					group    ced.Group
				}{rank, group})
			}
		}
	}

	// sort results
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i].distance < ranks[j].distance
	})

	// prepare result
	res := []ced.Group{}
	includedGroups := []ced.ID{}

	hasExact := false
	for _, rank := range ranks {
		if rank.distance == 0 {
			hasExact = true
		}
		// If we have exact match return nothing but exact matches.
		if hasExact && rank.distance != 0 {
			break
		}
		// If group is already included then skip.
		if slices.Contains(includedGroups, rank.group.ID) {
			continue
		}

		includedGroups = append(includedGroups, rank.group.ID)
		res = append(res, rank.group)
	}

	return res, nil
}
