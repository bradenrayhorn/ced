package responses

func MapSlice[T any, K any](in []T, mapper func(T) K) []K {
	var out []K
	for _, m := range in {
		mapped := mapper(m)

		out = append(out, mapped)
	}

	return out
}
