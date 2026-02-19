package set

type Set[T comparable] map[T]struct{}

func New[T comparable]() Set[T] {
	return make(Set[T])
}

func (set Set[T]) Add(item T) bool {
	if set.Has(item) {
		return false
	}

	set[item] = struct{}{}
	return true
}

func (set Set[T]) AddAll(items []T) {
	for _, item := range items {
		set[item] = struct{}{}
	}
}

func (set Set[T]) Has(item T) bool {
	_, exists := set[item]
	return exists
}

func (set Set[T]) Items() []T {
	result := make([]T, 0, len(set))
	for item := range set {
		result = append(result, item)
	}
	return result
}
