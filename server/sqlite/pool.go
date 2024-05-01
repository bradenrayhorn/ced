package sqlite

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	pool := sqlitemigration.NewPool(uri, schema, sqlitemigration.Options{
		PoolSize: 20,
		PrepareConn: func(conn *sqlite.Conn) error {
			conn.SetBusyTimeout(time.Second * 10)

			return errors.Join(
				sqlitex.ExecuteTransient(conn, "PRAGMA foreign_keys = ON;", nil),
				sqlitex.ExecuteTransient(conn, "PRAGMA synchronous = NORMAL;", nil),
			)
		},
		OnError: func(err error) {
			poolError = err
		},
	})

	// wait for pool to be ready
	conn, err := pool.Get(ctx)
	if err != nil {
		// if there is an error, it might have been due to the pool failing to connect for some reason
		code := sqlite.ErrCode(err)
		return nil, errors.Join(err, poolError, fmt.Errorf("sqlite code: %s, msg: %s", code.String(), code.Message()))
	}

	pool.Put(conn)

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
