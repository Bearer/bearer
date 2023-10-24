package slices

// Except returns a copy of the slice with the specified value removed
func Except[T comparable](slice []T, value T) []T {
	var result []T

	for _, candidate := range slice {
		if candidate != value {
			result = append(result, candidate)
		}
	}

	return result
}
