package maps

import (
	"cmp"
	"maps"
	"slices"
)

func SortedKeys[M ~map[K]V, K cmp.Ordered, V any](m M) []K {
	keys := slices.Collect(maps.Keys(m))
	slices.Sort(keys)
	return keys
}
