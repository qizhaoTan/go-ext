// Package sets provides a generic hash set.
//
// Set is a thin, allocation-conscious wrapper over map[T]struct{}. It is not
// safe for concurrent use; guard it with a mutex, or give each goroutine its
// own set and merge with Union at the end.
//
// A zero Set is not usable — its backing map is nil, so Add would panic.
// Construct one with New, WithCapacity or Of.
package sets

import (
	"fmt"
	"strings"
)

// Set is an unordered collection of distinct comparable values.
type Set[T comparable] struct {
	m map[T]struct{}
}

// New returns an empty set.
func New[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]struct{})}
}

// WithCapacity returns an empty set sized to hold about n elements without
// rehashing.
func WithCapacity[T comparable](n int) *Set[T] {
	return &Set[T]{m: make(map[T]struct{}, n)}
}

// Of returns a set holding the given values.
func Of[T comparable](values ...T) *Set[T] {
	s := WithCapacity[T](len(values))
	for _, v := range values {
		s.m[v] = struct{}{}
	}
	return s
}

// FromSlice returns a set holding the elements of the given slice.
func FromSlice[T comparable](values []T) *Set[T] {
	return Of(values...)
}

// Add inserts the given values and reports how many were not already present.
func (s *Set[T]) Add(values ...T) int {
	added := 0
	for _, v := range values {
		if _, ok := s.m[v]; !ok {
			s.m[v] = struct{}{}
			added++
		}
	}
	return added
}

// Remove deletes the given values and reports how many were present.
func (s *Set[T]) Remove(values ...T) int {
	removed := 0
	for _, v := range values {
		if _, ok := s.m[v]; ok {
			delete(s.m, v)
			removed++
		}
	}
	return removed
}

// Contains reports whether v is in the set.
func (s *Set[T]) Contains(v T) bool {
	_, ok := s.m[v]
	return ok
}

// ContainsAll reports whether every given value is in the set. It is true when
// no values are given.
func (s *Set[T]) ContainsAll(values ...T) bool {
	for _, v := range values {
		if !s.Contains(v) {
			return false
		}
	}
	return true
}

// ContainsAny reports whether at least one given value is in the set.
func (s *Set[T]) ContainsAny(values ...T) bool {
	for _, v := range values {
		if s.Contains(v) {
			return true
		}
	}
	return false
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int { return len(s.m) }

// IsEmpty reports whether the set has no elements.
func (s *Set[T]) IsEmpty() bool { return len(s.m) == 0 }

// Clear removes every element, keeping the allocated capacity.
func (s *Set[T]) Clear() { clear(s.m) }

// Clone returns a copy of the set. The copy shares no state with the original.
func (s *Set[T]) Clone() *Set[T] {
	out := WithCapacity[T](s.Len())
	for v := range s.m {
		out.m[v] = struct{}{}
	}
	return out
}

// ToSlice returns the elements in an unspecified order.
func (s *Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s.m))
	for v := range s.m {
		result = append(result, v)
	}
	return result
}

// All returns an iterator over the elements, in an unspecified order. It is
// intended for use with a range statement:
//
//	for v := range s.All() {
//		// ...
//	}
func (s *Set[T]) All() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		for v := range s.m {
			if !yield(v) {
				return
			}
		}
	}
}

// Union returns a new set holding the elements of s and of every given set.
func (s *Set[T]) Union(others ...*Set[T]) *Set[T] {
	out := s.Clone()
	for _, o := range others {
		for v := range o.m {
			out.m[v] = struct{}{}
		}
	}
	return out
}

// Intersection returns a new set holding the elements present in s and in every
// given set.
func (s *Set[T]) Intersection(others ...*Set[T]) *Set[T] {
	out := New[T]()
	// Iterating the smallest input keeps the work proportional to it rather
	// than to the receiver, which matters when sizes are lopsided.
	smallest := s
	for _, o := range others {
		if o.Len() < smallest.Len() {
			smallest = o
		}
	}
	for v := range smallest.m {
		if s.Contains(v) && allContain(others, v) {
			out.m[v] = struct{}{}
		}
	}
	return out
}

// Difference returns a new set holding the elements of s that are in none of
// the given sets.
func (s *Set[T]) Difference(others ...*Set[T]) *Set[T] {
	out := New[T]()
	for v := range s.m {
		if !anyContains(others, v) {
			out.m[v] = struct{}{}
		}
	}
	return out
}

// SymmetricDifference returns a new set holding the elements in exactly one of
// s and other.
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
	out := New[T]()
	for v := range s.m {
		if !other.Contains(v) {
			out.m[v] = struct{}{}
		}
	}
	for v := range other.m {
		if !s.Contains(v) {
			out.m[v] = struct{}{}
		}
	}
	return out
}

// IsSubsetOf reports whether every element of s is in other.
func (s *Set[T]) IsSubsetOf(other *Set[T]) bool {
	if s.Len() > other.Len() {
		return false
	}
	for v := range s.m {
		if !other.Contains(v) {
			return false
		}
	}
	return true
}

// IsSupersetOf reports whether every element of other is in s.
func (s *Set[T]) IsSupersetOf(other *Set[T]) bool {
	return other.IsSubsetOf(s)
}

// Equal reports whether s and other hold the same elements.
func (s *Set[T]) Equal(other *Set[T]) bool {
	return s.Len() == other.Len() && s.IsSubsetOf(other)
}

// String implements fmt.Stringer. The element order is unspecified.
func (s *Set[T]) String() string {
	var b strings.Builder
	b.WriteString("Set{")
	first := true
	for v := range s.m {
		if !first {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%v", v)
		first = false
	}
	b.WriteByte('}')
	return b.String()
}

func allContain[T comparable](sets []*Set[T], v T) bool {
	for _, s := range sets {
		if !s.Contains(v) {
			return false
		}
	}
	return true
}

func anyContains[T comparable](sets []*Set[T], v T) bool {
	for _, s := range sets {
		if s.Contains(v) {
			return true
		}
	}
	return false
}
