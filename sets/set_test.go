package sets_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/qizhaoTan/go-ext/sets"
)

func sortedSlice(s *sets.Set[int]) []int {
	out := s.ToSlice()
	sort.Ints(out)
	return out
}

func TestNewAndOf(t *testing.T) {
	if s := sets.New[int](); !s.IsEmpty() || s.Len() != 0 {
		t.Errorf("New() = %v, want empty", s)
	}
	if s := sets.Of(1, 2, 2, 3); s.Len() != 3 {
		t.Errorf("Of() len = %d, want 3 (duplicates collapse)", s.Len())
	}
	if s := sets.FromSlice([]string{"a", "a"}); s.Len() != 1 {
		t.Errorf("FromSlice() len = %d, want 1", s.Len())
	}
}

func TestAddReportsNewElements(t *testing.T) {
	s := sets.New[int]()
	if got := s.Add(1, 2, 2); got != 2 {
		t.Errorf("Add() = %d, want 2", got)
	}
	if got := s.Add(2, 3); got != 1 {
		t.Errorf("Add() again = %d, want 1", got)
	}
}

func TestRemoveReportsPresentElements(t *testing.T) {
	s := sets.Of(1, 2, 3)
	if got := s.Remove(1, 99); got != 1 {
		t.Errorf("Remove() = %d, want 1", got)
	}
	if s.Contains(1) {
		t.Error("Contains(1) = true after Remove, want false")
	}
}

func TestContainsAllAndAny(t *testing.T) {
	s := sets.Of(1, 2, 3)
	if !s.ContainsAll(1, 2) {
		t.Error("ContainsAll(1,2) = false, want true")
	}
	if s.ContainsAll(1, 9) {
		t.Error("ContainsAll(1,9) = true, want false")
	}
	if !s.ContainsAny(9, 3) {
		t.Error("ContainsAny(9,3) = false, want true")
	}
	if s.ContainsAny(8, 9) {
		t.Error("ContainsAny(8,9) = true, want false")
	}
	// Vacuously true with no arguments.
	if !s.ContainsAll() {
		t.Error("ContainsAll() = false, want true")
	}
}

func TestClearAndClone(t *testing.T) {
	s := sets.Of(1, 2, 3)
	clone := s.Clone()
	s.Clear()

	if !s.IsEmpty() {
		t.Errorf("Clear() left %d elements", s.Len())
	}
	if clone.Len() != 3 {
		t.Errorf("Clone was affected by Clear: len = %d, want 3", clone.Len())
	}
}

func TestUnion(t *testing.T) {
	got := sets.Of(1, 2).Union(sets.Of(2, 3), sets.Of(4))
	if want := []int{1, 2, 3, 4}; !reflect.DeepEqual(sortedSlice(got), want) {
		t.Errorf("Union() = %v, want %v", sortedSlice(got), want)
	}
}

func TestUnionDoesNotMutateReceiver(t *testing.T) {
	s := sets.Of(1, 2)
	s.Union(sets.Of(3))
	if s.Len() != 2 {
		t.Errorf("Union mutated receiver: len = %d, want 2", s.Len())
	}
}

func TestIntersection(t *testing.T) {
	got := sets.Of(1, 2, 3, 4).Intersection(sets.Of(2, 3, 4), sets.Of(3, 4, 5))
	if want := []int{3, 4}; !reflect.DeepEqual(sortedSlice(got), want) {
		t.Errorf("Intersection() = %v, want %v", sortedSlice(got), want)
	}
	// With no arguments the receiver is unchanged.
	if got := sets.Of(1, 2).Intersection(); got.Len() != 2 {
		t.Errorf("Intersection() with no args len = %d, want 2", got.Len())
	}
}

func TestIntersectionIsCorrectWhenReceiverIsLarger(t *testing.T) {
	// Exercises the "iterate the smallest input" path.
	large := sets.Of(1, 2, 3, 4, 5, 6, 7, 8)
	got := large.Intersection(sets.Of(3, 99))
	if want := []int{3}; !reflect.DeepEqual(sortedSlice(got), want) {
		t.Errorf("Intersection() = %v, want %v", sortedSlice(got), want)
	}
}

func TestDifference(t *testing.T) {
	got := sets.Of(1, 2, 3, 4).Difference(sets.Of(2), sets.Of(4))
	if want := []int{1, 3}; !reflect.DeepEqual(sortedSlice(got), want) {
		t.Errorf("Difference() = %v, want %v", sortedSlice(got), want)
	}
}

func TestSymmetricDifference(t *testing.T) {
	got := sets.Of(1, 2, 3).SymmetricDifference(sets.Of(3, 4))
	if want := []int{1, 2, 4}; !reflect.DeepEqual(sortedSlice(got), want) {
		t.Errorf("SymmetricDifference() = %v, want %v", sortedSlice(got), want)
	}
}

func TestSubsetSupersetEqual(t *testing.T) {
	small, large := sets.Of(1, 2), sets.Of(1, 2, 3)

	if !small.IsSubsetOf(large) {
		t.Error("IsSubsetOf() = false, want true")
	}
	if large.IsSubsetOf(small) {
		t.Error("IsSubsetOf() reversed = true, want false")
	}
	if !large.IsSupersetOf(small) {
		t.Error("IsSupersetOf() = false, want true")
	}
	if small.Equal(large) {
		t.Error("Equal() = true, want false")
	}
	if !sets.Of(1, 2).Equal(sets.Of(2, 1)) {
		t.Error("Equal() on same elements = false, want true")
	}
	// The empty set is a subset of everything.
	if !sets.New[int]().IsSubsetOf(large) {
		t.Error("empty IsSubsetOf() = false, want true")
	}
}

func TestAllIterates(t *testing.T) {
	s := sets.Of(1, 2, 3)
	var got []int
	for v := range s.All() {
		got = append(got, v)
	}
	sort.Ints(got)
	if want := []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
		t.Errorf("All() yielded %v, want %v", got, want)
	}
}

func TestAllStopsOnBreak(t *testing.T) {
	s := sets.Of(1, 2, 3, 4, 5)
	count := 0
	for range s.All() {
		count++
		if count == 2 {
			break
		}
	}
	if count != 2 {
		t.Errorf("All() yielded %d elements after break, want 2", count)
	}
}

func TestString(t *testing.T) {
	if got := sets.Of(1).String(); got != "Set{1}" {
		t.Errorf("String() = %q, want %q", got, "Set{1}")
	}
	if got := sets.New[int]().String(); got != "Set{}" {
		t.Errorf("String() = %q, want %q", got, "Set{}")
	}
}
