package slices

import "github.com/qizhaoTan/go-ext/constraints"

// Find returns the first element for which pred reports true, and whether such
// an element exists. When it does not, the zero value of T is returned.
func Find[T any](s []T, pred func(T) bool) (T, bool) {
	for _, v := range s {
		if pred(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex returns the index of the first element for which pred reports true,
// or -1 if there is none.
func FindIndex[T any](s []T, pred func(T) bool) int {
	for i, v := range s {
		if pred(v) {
			return i
		}
	}
	return -1
}

// FindLast returns the last element for which pred reports true, and whether
// such an element exists.
func FindLast[T any](s []T, pred func(T) bool) (T, bool) {
	for i := len(s) - 1; i >= 0; i-- {
		if pred(s[i]) {
			return s[i], true
		}
	}
	var zero T
	return zero, false
}

// Some reports whether pred is true for at least one element of s.
func Some[T any](s []T, pred func(T) bool) bool {
	return FindIndex(s, pred) >= 0
}

// Every reports whether pred is true for every element of s. It is true for an
// empty slice.
func Every[T any](s []T, pred func(T) bool) bool {
	for _, v := range s {
		if !pred(v) {
			return false
		}
	}
	return true
}

// None reports whether pred is false for every element of s.
func None[T any](s []T, pred func(T) bool) bool {
	return !Some(s, pred)
}

// Count returns the number of elements of s for which pred reports true.
func Count[T any](s []T, pred func(T) bool) int {
	n := 0
	for _, v := range s {
		if pred(v) {
			n++
		}
	}
	return n
}

// GroupBy partitions s into a map keyed by key, preserving the input order
// within each group.
func GroupBy[T any, K comparable](s []T, key func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range s {
		k := key(v)
		result[k] = append(result[k], v)
	}
	return result
}

// Partition splits s into the elements for which pred is true and those for
// which it is false, preserving order in both.
func Partition[T any](s []T, pred func(T) bool) (matched, unmatched []T) {
	matched = make([]T, 0, len(s))
	unmatched = make([]T, 0, len(s))
	for _, v := range s {
		if pred(v) {
			matched = append(matched, v)
		} else {
			unmatched = append(unmatched, v)
		}
	}
	return matched, unmatched
}

// KeyBy indexes s by key. When two elements share a key, the later one wins.
func KeyBy[T any, K comparable](s []T, key func(T) K) map[K]T {
	result := make(map[K]T, len(s))
	for _, v := range s {
		result[key(v)] = v
	}
	return result
}

// Sum returns the sum of the elements of s, or zero for an empty slice.
func Sum[T constraints.Number](s []T) T {
	var total T
	for _, v := range s {
		total += v
	}
	return total
}

// SumBy returns the sum of fn applied to each element of s.
func SumBy[T any, R constraints.Number](s []T, fn func(T) R) R {
	var total R
	for _, v := range s {
		total += fn(v)
	}
	return total
}

// MaxBy returns the element of s with the greatest key, and whether s was
// non-empty. Ties go to the first such element.
func MaxBy[T any, K constraints.Number | ~string](s []T, key func(T) K) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	best, bestKey := s[0], key(s[0])
	for _, v := range s[1:] {
		if k := key(v); k > bestKey {
			best, bestKey = v, k
		}
	}
	return best, true
}

// MinBy returns the element of s with the smallest key, and whether s was
// non-empty. Ties go to the first such element.
func MinBy[T any, K constraints.Number | ~string](s []T, key func(T) K) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	best, bestKey := s[0], key(s[0])
	for _, v := range s[1:] {
		if k := key(v); k < bestKey {
			best, bestKey = v, k
		}
	}
	return best, true
}
