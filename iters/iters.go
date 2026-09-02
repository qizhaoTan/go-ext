// Package iters provides lazy sequence operations over Go 1.23 iterators.
//
// Where the slices package materialises a new slice at every step, the
// functions here compose without allocating: a Map over a Filter over a Take
// does one pass and touches only the elements actually consumed. That is the
// difference that matters on large or infinite sequences.
//
//	// Reads exactly three elements, however long src is.
//	first := iters.Collect(iters.Take(iters.Filter(src, isEven), 3))
//
// Every function returns an [iter.Seq] that may be ranged over. Sequences are
// consumed lazily, so nothing runs until the result is ranged over or passed to
// a terminal operation such as Collect, Reduce or Count.
package iters

import "iter"

// Map returns a sequence of fn applied to each element of seq.
func Map[T, R any](seq iter.Seq[T], fn func(T) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

// Filter returns the elements of seq for which pred reports true.
func Filter[T any](seq iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if pred(v) && !yield(v) {
				return
			}
		}
	}
}

// FilterMap applies fn to each element and yields the results for which fn
// reports true.
func FilterMap[T, R any](seq iter.Seq[T], fn func(T) (R, bool)) iter.Seq[R] {
	return func(yield func(R) bool) {
		for v := range seq {
			if r, ok := fn(v); ok && !yield(r) {
				return
			}
		}
	}
}

// Take returns at most the first n elements of seq. It stops pulling from seq
// once n elements have been yielded, which is what makes it safe over an
// infinite sequence.
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		if n <= 0 {
			return
		}
		count := 0
		for v := range seq {
			if !yield(v) {
				return
			}
			count++
			if count >= n {
				return
			}
		}
	}
}

// Drop returns seq with its first n elements skipped.
func Drop[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		skipped := 0
		for v := range seq {
			if skipped < n {
				skipped++
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// TakeWhile returns the leading elements of seq for which pred reports true,
// stopping at the first element that does not.
func TakeWhile[T any](seq iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if !pred(v) || !yield(v) {
				return
			}
		}
	}
}

// DropWhile returns seq with its leading elements for which pred reports true
// skipped.
func DropWhile[T any](seq iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		dropping := true
		for v := range seq {
			if dropping {
				if pred(v) {
					continue
				}
				dropping = false
			}
			if !yield(v) {
				return
			}
		}
	}
}

// Concat returns the given sequences one after another.
func Concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Collect materialises seq into a slice.
//
// It never returns from an infinite sequence: bound one with Take first.
func Collect[T any](seq iter.Seq[T]) []T {
	var result []T
	for v := range seq {
		result = append(result, v)
	}
	return result
}

// Reduce folds seq into a single value, starting from initial.
func Reduce[T, R any](seq iter.Seq[T], initial R, fn func(R, T) R) R {
	acc := initial
	for v := range seq {
		acc = fn(acc, v)
	}
	return acc
}

// Count returns the number of elements in seq.
func Count[T any](seq iter.Seq[T]) int {
	n := 0
	for range seq {
		n++
	}
	return n
}

// Find returns the first element for which pred reports true, and whether such
// an element exists. It stops consuming seq as soon as one is found.
func Find[T any](seq iter.Seq[T], pred func(T) bool) (T, bool) {
	for v := range seq {
		if pred(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Some reports whether pred is true for at least one element of seq.
func Some[T any](seq iter.Seq[T], pred func(T) bool) bool {
	_, found := Find(seq, pred)
	return found
}

// Every reports whether pred is true for every element of seq.
func Every[T any](seq iter.Seq[T], pred func(T) bool) bool {
	for v := range seq {
		if !pred(v) {
			return false
		}
	}
	return true
}

// Of returns a sequence over the given values.
func Of[T any](values ...T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
}

// Range returns the integers from start (inclusive) to stop (exclusive).
func Range(start, stop int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := start; i < stop; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Repeat returns an infinite sequence of v. Bound it with Take.
func Repeat[T any](v T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for yield(v) {
		}
	}
}

// Iterate returns the infinite sequence seed, fn(seed), fn(fn(seed)), and so
// on. Bound it with Take or TakeWhile.
func Iterate[T any](seed T, fn func(T) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		v := seed
		for yield(v) {
			v = fn(v)
		}
	}
}
