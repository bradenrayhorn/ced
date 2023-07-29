package sqlite

import (
	"fmt"

	"github.com/bradenrayhorn/ced/ced"
)

func mustFindResult[T comparable](result T, err error) (T, error) {
	if err != nil {
		return result, err
	}

	if result == *new(T) {
		return result, ced.NewError(ced.ENOTFOUND, fmt.Sprintf("%T not found.", result))
	}

	return result, nil
}
