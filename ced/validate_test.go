package ced

import (
	"errors"
	"testing"

	"github.com/matryer/is"
)

// validate field

type StructPass struct{}

func (s StructPass) Validate() error {
	return nil
}

type StructError1 struct{}

func (s StructError1) Validate() error {
	return errors.New(":field error 1")
}

type StructError2 struct{}

func (s StructError2) Validate() error {
	return errors.New(":field error 2")
}

func TestValidateField(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("one passing", func(t *testing.T) {
		result := ValidateFields(Field("Field 1", StructPass{}))
		is.NoErr(result)
	})

	t.Run("one failing", func(t *testing.T) {
		result := ValidateFields(Field("Field 1", StructError1{}))
		_, msg := result.(Error).CedError()
		is.Equal("Field 1 error 1.", msg)
	})

	t.Run("one failing and one passing", func(t *testing.T) {
		result := ValidateFields(Field("Field 1", StructError1{}, StructPass{}))
		_, msg := result.(Error).CedError()
		is.Equal("Field 1 error 1.", msg)
	})

	t.Run("multiple failing", func(t *testing.T) {
		result := ValidateFields(Field("Field 1", StructError1{}, StructPass{}, StructError2{}))
		_, msg := result.(Error).CedError()
		is.Equal("Field 1 error 1, Field 1 error 2.", msg)
	})

	t.Run("multiple fields", func(t *testing.T) {
		result := ValidateFields(
			Field("Field 1", StructError1{}, StructError2{}),
			Field("Field 2", StructError2{}),
		)
		_, msg := result.(Error).CedError()
		is.Equal("Field 1 error 1, Field 1 error 2. Field 2 error 2.", msg)
	})
}

// is required

type StructEmpty struct{ empty bool }

func (s StructEmpty) Empty() bool {
	return s.empty
}

func TestRequired(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("fails if empty", func(t *testing.T) {
		err := Required(StructEmpty{empty: true}).Validate()
		is.Equal(":field is required", err.Error())
	})

	t.Run("succeeds if not empty", func(t *testing.T) {
		err := Required(StructEmpty{empty: false}).Validate()
		is.NoErr(err)
	})
}

// max

type StructMax struct{ length int }

func (s StructMax) Length() int {
	return s.length
}

func TestMaxLength(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("fails if over max", func(t *testing.T) {
		err := Max(StructMax{length: 2}, 1, "characters").Validate()
		is.Equal(":field must be at most 1 characters", err.Error())
	})

	t.Run("succeeds if at max", func(t *testing.T) {
		err := Max(StructMax{length: 2}, 2, "characters").Validate()
		is.NoErr(err)
	})

	t.Run("succeeds if below max", func(t *testing.T) {
		err := Max(StructMax{length: 0}, 1, "characters").Validate()
		is.NoErr(err)
	})
}
