package ced

import (
	"github.com/segmentio/ksuid"
)

type ID ksuid.KSUID

func (id ID) String() string {
	return ksuid.KSUID(id).String()
}

func (id ID) Empty() bool {
	return ksuid.KSUID(id).IsNil()
}

func NewID() ID {
	return ID(ksuid.New())
}

func IDFromString(string string) (ID, error) {
	id, err := ksuid.Parse(string)
	if err != nil {
		return ID(ksuid.Nil), err
	}

	return ID(id), nil
}
