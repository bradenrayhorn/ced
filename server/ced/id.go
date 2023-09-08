package ced

import (
	"encoding/json"
	"fmt"

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
		return ID(ksuid.Nil), fmt.Errorf("invalid id: %s", string)
	}

	return ID(id), nil
}

func (id *ID) UnmarshalJSON(b []byte) error {
	var idString string
	if err := json.Unmarshal(b, &idString); err != nil {
		return err
	}

	parsedID, err := IDFromString(idString)
	if err != nil {
		return err
	}
	*id = parsedID
	return nil
}

func (id ID) MarshalJSON() ([]byte, error) {
	if id.Empty() {
		return json.Marshal(nil)
	}
	return json.Marshal(id.String())
}
