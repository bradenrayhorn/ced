package testutils

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/ced/server/sqlite"
)

var InMemoryPoolPath = "file:testdb?mode=memory&cache=shared"

func StartPool(tb testing.TB) (*sqlite.Pool, func()) {
	pool, err := sqlite.CreatePool(
		context.Background(),
		InMemoryPoolPath,
	)
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
