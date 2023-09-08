package ced

import "strings"

type ValidatableString string

func (s ValidatableString) Empty() bool {
	return strings.TrimSpace(string(s)) == ""
}

func (s ValidatableString) Length() int {
	return len(string(s))
}
