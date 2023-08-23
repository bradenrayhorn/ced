package sqlite

import (
	"context"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func execute(ctx context.Context, pool *Pool, query string, args []any) error {
	conn, done, err := pool.Conn(ctx)
	if err != nil {
		return err
	}
	defer done()

	return sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Args: args,
	})
}

func selectOne[T comparable](ctx context.Context, pool *Pool, query string, args []any, mapper func(stmt *sqlite.Stmt) (T, error)) (T, error) {
	var result T
	conn, done, err := pool.Conn(ctx)
	if err != nil {
		return result, err
	}
	defer done()

	err = sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Args: args,
		ResultFunc: func(stmt *sqlite.Stmt) error {
			mapped, err := mapper(stmt)
			if err != nil {
				return err
			}

			result = mapped
			return nil
		},
	})

	return result, err
}

func selectList[T any](ctx context.Context, pool *Pool, query string, args []any, mapper func(stmt *sqlite.Stmt) (T, error)) ([]T, error) {
	result := []T{}
	conn, done, err := pool.Conn(ctx)
	if err != nil {
		return result, err
	}
	defer done()

	return result, sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Args: args,
		ResultFunc: func(stmt *sqlite.Stmt) error {
			mapped, err := mapper(stmt)
			if err != nil {
				return err
			}

			result = append(result, mapped)
			return nil
		},
	})
}
