package ced

import (
	"errors"
	"fmt"
	"strings"
)

type Validatable interface {
	Validate() error
}

func ValidateFields(objects ...validatableField) error {
	messages := []string{}
	for _, o := range objects {
		if err := o.Validate(); err != nil {
			messages = append(messages, err.Error()+".")
		}
	}

	if len(messages) == 0 {
		return nil
	}

	return NewError(EINVALID, strings.Join(messages, " "))
}

func validate(objects ...Validatable) error {
	messages := []string{}
	for _, o := range objects {
		if err := o.Validate(); err != nil {
			messages = append(messages, err.Error())
		}
	}

	if len(messages) == 0 {
		return nil
	}

	return errors.New(strings.Join(messages, ", "))
}

// validatable field

type validatableField struct {
	name    string
	objects []Validatable
}

func (f validatableField) Validate() error {
	err := validate(f.objects...)
	if err != nil {
		return errors.New(strings.ReplaceAll(err.Error(), ":field", f.name))
	}
	return nil
}

func Field(name string, objects ...Validatable) validatableField {
	return validatableField{name, objects}
}

// required validator

type Emptiable interface {
	Empty() bool
}

type validatableEmptiable struct {
	Emptiable
}

func (e validatableEmptiable) Validate() error {
	if e.Empty() {
		return errors.New(":field is required")
	}

	return nil
}

func Required(e Emptiable) Validatable {
	return validatableEmptiable{e}
}

// max length validator

type Countable interface {
	Length() int
}

type valitableCountable struct {
	Countable
	max    int
	object string
}

func (c valitableCountable) Validate() error {
	if c.Length() > c.max {
		return fmt.Errorf(":field must be at most %d %s", c.max, c.object)
	}
	return nil
}

func Max(c Countable, max int, object string) valitableCountable {
	return valitableCountable{c, max, object}
}
