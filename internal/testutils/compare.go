package testutils

import (
	"sort"

	"github.com/bradenrayhorn/ced/ced"
)

func SortSlice[T any](s []T, comp func(a, b T) bool) []T {
	sort.Slice(s, func(i, j int) bool {
		return comp(s[i], s[j])
	})

	return s
}

func CompareIndividuals(a ced.Individual, b ced.Individual) bool {
	return a.ID.String() < b.ID.String()
}
