package sqlite

import (
	"context"

	"github.com/bradenrayhorn/ced/ced"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitemigration"
	"zombiezen.com/go/sqlite/sqlitex"
)

type Pool struct {
	pool *sqlitex.Pool
}

func CreatePool(ctx context.Context, config ced.Config) (*Pool, error) {
	pool, err := sqlitex.Open(config.DbPath, 0, 10)
	if err != nil {
		return nil, err
	}

	schema := sqlitemigration.Schema{
		Migrations: migrations,
	}

	conn := pool.Get(ctx)
	defer pool.Put(conn)

	err = sqlitemigration.Migrate(ctx, conn, schema)
	if err != nil {
		return nil, err
	}
	return &Pool{pool}, nil
}

func (p *Pool) Close(ctx context.Context) error {
	return p.pool.Close()
}

func (p *Pool) Conn(ctx context.Context) (*sqlite.Conn, func()) {
	conn := p.pool.Get(ctx)

	return conn, func() {
		p.pool.Put(conn)
	}
}
