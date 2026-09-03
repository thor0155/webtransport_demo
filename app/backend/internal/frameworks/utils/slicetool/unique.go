package slicetool

func Unique[T comparable](input []T) []T {
	return UniqueBy(input, func(v T) T { return v })
}

func UniqueBy[T any, K comparable](input []T, key func(T) K) []T {
	seen := make(map[K]struct{}, len(input))
	result := make([]T, 0, len(input))

	for _, v := range input {
		k := key(v)

		if _, exists := seen[k]; exists {
			continue
		}

		seen[k] = struct{}{}
		result = append(result, v)
	}

	return result
}

func MapUnique[T any, R comparable](input []T, fn func(T) R) []R {
	seen := make(map[R]struct{}, len(input))
	result := make([]R, 0, len(input))

	for _, v := range input {
		value := fn(v)

		if _, exists := seen[value]; exists {
			continue
		}

		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result
}
