package sqlite

import (
	"context"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitemigration"
	"zombiezen.com/go/sqlite/sqlitex"
)

type Pool struct {
	pool *sqlitemigration.Pool
}

func CreatePool(ctx context.Context, uri string) (*Pool, error) {
	schema := sqlitemigration.Schema{
		Migrations: migrations,
	}

	var poolError error
	ready := make(chan int)
	pool := sqlitemigration.NewPool(uri, schema, sqlitemigration.Options{
		PoolSize: 20,
		PrepareConn: func(conn *sqlite.Conn) error {
			return sqlitex.ExecuteTransient(conn, "PRAGMA foreign_keys = ON;", nil)
		},
		OnReady: func() {
			ready <- 1
		},
		OnError: func(err error) {
			poolError = err
			ready <- 1
		},
	})

	<-ready
	if poolError != nil {
		return nil, poolError
	}

	return &Pool{pool}, nil
}

func (p *Pool) Close(ctx context.Context) error {
	return p.pool.Close()
}

func (p *Pool) Conn(ctx context.Context) (*sqlite.Conn, func(), error) {
	conn, err := p.pool.Get(ctx)

	return conn, func() {
		if err == nil && conn != nil {
			p.pool.Put(conn)
		}
	}, err
}
