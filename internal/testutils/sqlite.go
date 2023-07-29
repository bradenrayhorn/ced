package testutils

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/bradenrayhorn/ced/sqlite"
)

func StartPool(tb testing.TB) (*sqlite.Pool, func()) {
	pool, err := sqlite.CreatePool(context.Background(), ced.Config{
		DbPath: "file:testdb?mode=memory&cache=shared",
	})
	if err != nil {
		tb.Fatal(err)
	}

	return pool, func() {
		err = pool.Close(context.Background())
		if err != nil {
			tb.Fatal(err)
		}
	}
}
