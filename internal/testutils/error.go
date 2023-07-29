package testutils

import (
	"errors"
	"testing"

	"github.com/bradenrayhorn/ced/ced"
	"github.com/matryer/is"
)

func IsCodeAndError(t testing.TB, err error, code string, msg string) {
	IsCode(t, err, code)
	IsError(t, err, msg)
}

func IsCode(t testing.TB, err error, code string) {
	is := is.New(t)
	is.True(err != nil)

	var cedError ced.Error
	is.True(errors.As(err, &cedError))
	foundCode, _ := cedError.CedError()
	is.Equal(foundCode, code)
}

func IsError(t testing.TB, err error, msg string) {
	is := is.New(t)
	is.True(err != nil)

	var cedError ced.Error
	is.True(errors.As(err, &cedError))
	_, foundMsg := cedError.CedError()
	is.Equal(foundMsg, msg)
}
