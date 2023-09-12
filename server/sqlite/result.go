package sqlite

import (
	"fmt"

	"github.com/bradenrayhorn/ced/server/ced"
)

func mustFindResult[T comparable](result T, err error, identifier string) (T, error) {
	if err != nil {
		return result, err
	}

	if result == *new(T) {
		return result, ced.NewError(ced.ENOTFOUND, fmt.Sprintf("%T [%s] not found", result, identifier))
	}

	return result, nil
}
