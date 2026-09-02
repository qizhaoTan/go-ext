// Package maps provides generic helpers for Go maps.
//
// The standard library's [maps] package covers iteration and cloning. This
// package adds the transformations built on top of that: filtering, inverting
// and merging.
//
// Go maps have no defined iteration order, so no function here promises one.
// Callers that need deterministic output should sort the result.
package maps

// Keys returns the keys of m in an unspecified order.
func Keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// Values returns the values of m in an unspecified order.
func Values[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

// Entries returns the key/value pairs of m in an unspecified order.
func Entries[K comparable, V any](m map[K]V) []Entry[K, V] {
	result := make([]Entry[K, V], 0, len(m))
	for k, v := range m {
		result = append(result, Entry[K, V]{Key: k, Value: v})
	}
	return result
}

// Entry is a single key/value pair, used to move map contents through
// slice-shaped code.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// FromEntries builds a map from the given pairs. When two entries share a key,
// the later one wins.
func FromEntries[K comparable, V any](entries []Entry[K, V]) map[K]V {
	result := make(map[K]V, len(entries))
	for _, e := range entries {
		result[e.Key] = e.Value
	}
	return result
}

// Filter returns the entries of m for which pred reports true.
func Filter[K comparable, V any](m map[K]V, pred func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if pred(k, v) {
			result[k] = v
		}
	}
	return result
}

// MapValues returns a map with the same keys as m and each value replaced by
// fn applied to it.
func MapValues[K comparable, V, R any](m map[K]V, fn func(K, V) R) map[K]R {
	result := make(map[K]R, len(m))
	for k, v := range m {
		result[k] = fn(k, v)
	}
	return result
}

// MapKeys returns a map with each key replaced by fn applied to it. When two
// keys collide, which entry survives is unspecified, because map iteration
// order is unspecified.
func MapKeys[K comparable, V any, R comparable](m map[K]V, fn func(K, V) R) map[R]V {
	result := make(map[R]V, len(m))
	for k, v := range m {
		result[fn(k, v)] = v
	}
	return result
}

// Invert swaps keys and values. When two keys share a value, which one survives
// is unspecified.
func Invert[K, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K, len(m))
	for k, v := range m {
		result[v] = k
	}
	return result
}

// Merge returns a new map holding the entries of every given map. Later maps
// take precedence over earlier ones.
func Merge[K comparable, V any](ms ...map[K]V) map[K]V {
	size := 0
	for _, m := range ms {
		size += len(m)
	}
	result := make(map[K]V, size)
	for _, m := range ms {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// GetOr returns m[key] if present, and fallback otherwise.
func GetOr[K comparable, V any](m map[K]V, key K, fallback V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return fallback
}
