// Package slices provides generic, functional-style helpers for Go slices.
//
// The functions here complement the standard library's [slices] package rather
// than replacing it: sorting, searching and equality already live there, so this
// package focuses on the transformations Go's standard library leaves out.
//
// All functions are pure. Nothing mutates its input, and every returned slice is
// freshly allocated, so a result can be handed to a caller without aliasing
// surprises. A nil input is always treated as an empty slice.
package slices

// Map returns a new slice holding fn applied to each element of s.
func Map[T, R any](s []T, fn func(T) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// MapIndex is Map with the element index passed to fn.
func MapIndex[T, R any](s []T, fn func(int, T) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = fn(i, v)
	}
	return result
}

// Filter returns the elements of s for which pred reports true, in order.
func Filter[T any](s []T, pred func(T) bool) []T {
	result := make([]T, 0, len(s))
	for _, v := range s {
		if pred(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reject is the complement of Filter: it keeps the elements for which pred
// reports false.
func Reject[T any](s []T, pred func(T) bool) []T {
	return Filter(s, func(v T) bool { return !pred(v) })
}

// FilterMap applies fn to each element and keeps the results for which fn
// reports true. It fuses a Filter and a Map into a single pass, which matters
// when the mapping is the expensive part of the predicate.
func FilterMap[T, R any](s []T, fn func(T) (R, bool)) []R {
	result := make([]R, 0, len(s))
	for _, v := range s {
		if r, ok := fn(v); ok {
			result = append(result, r)
		}
	}
	return result
}

// FlatMap applies fn to each element and concatenates the resulting slices.
func FlatMap[T, R any](s []T, fn func(T) []R) []R {
	result := make([]R, 0, len(s))
	for _, v := range s {
		result = append(result, fn(v)...)
	}
	return result
}

// Flatten concatenates the slices of s into a single slice.
func Flatten[T any](s [][]T) []T {
	size := 0
	for _, inner := range s {
		size += len(inner)
	}
	result := make([]T, 0, size)
	for _, inner := range s {
		result = append(result, inner...)
	}
	return result
}

// Reduce folds s into a single value, starting from initial and combining
// left to right.
func Reduce[T, R any](s []T, initial R, fn func(R, T) R) R {
	acc := initial
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

// ReduceRight folds s right to left.
func ReduceRight[T, R any](s []T, initial R, fn func(R, T) R) R {
	acc := initial
	for i := len(s) - 1; i >= 0; i-- {
		acc = fn(acc, s[i])
	}
	return acc
}

// Reverse returns a new slice with the elements of s in reverse order.
//
// The standard library's slices.Reverse reverses in place; this is the
// non-mutating counterpart.
func Reverse[T any](s []T) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}

// Uniq returns the elements of s with duplicates removed, keeping the first
// occurrence of each value.
func Uniq[T comparable](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	result := make([]T, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

// UniqBy returns the elements of s whose key, as computed by key, is seen for
// the first time.
func UniqBy[T any, K comparable](s []T, key func(T) K) []T {
	seen := make(map[K]struct{}, len(s))
	result := make([]T, 0, len(s))
	for _, v := range s {
		k := key(v)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		result = append(result, v)
	}
	return result
}

// Chunk splits s into consecutive chunks of at most size elements. The final
// chunk is short when len(s) is not a multiple of size.
//
// Chunk panics if size is not positive: a non-positive chunk size has no
// meaningful result, and returning an empty slice would hide the caller's bug.
func Chunk[T any](s []T, size int) [][]T {
	if size <= 0 {
		panic("go-ext/slices: Chunk called with non-positive size")
	}
	result := make([][]T, 0, (len(s)+size-1)/size)
	for i := 0; i < len(s); i += size {
		end := min(i+size, len(s))
		// Copy rather than reslice so the chunks do not alias s.
		chunk := make([]T, end-i)
		copy(chunk, s[i:end])
		result = append(result, chunk)
	}
	return result
}

// Compact returns the elements of s that are not the zero value of T.
func Compact[T comparable](s []T) []T {
	var zero T
	return Filter(s, func(v T) bool { return v != zero })
}
