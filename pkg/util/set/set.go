package set

import "golang.org/x/exp/maps"

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
		set.Add(item)
	}
}

func (set Set[T]) Has(item T) bool {
	_, exists := set[item]
	return exists
}

func (set Set[T]) Items() []T {
	return maps.Keys(set)
}
