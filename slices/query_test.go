package slices_test

import (
	"reflect"
	"testing"

	xslices "github.com/qizhaoTan/go-ext/slices"
)

type person struct {
	Name string
	Age  int
}

var people = []person{
	{"Alice", 30},
	{"Bob", 25},
	{"Carol", 35},
	{"Dave", 25},
}

func TestFind(t *testing.T) {
	got, ok := xslices.Find([]int{1, 2, 3}, func(n int) bool { return n > 1 })
	if !ok || got != 2 {
		t.Errorf("Find() = (%v, %v), want (2, true)", got, ok)
	}

	got, ok = xslices.Find([]int{1, 2, 3}, func(n int) bool { return n > 10 })
	if ok || got != 0 {
		t.Errorf("Find() miss = (%v, %v), want (0, false)", got, ok)
	}
}

func TestFindIndexAndFindLast(t *testing.T) {
	in := []int{1, 2, 3, 2}
	if got := xslices.FindIndex(in, func(n int) bool { return n == 2 }); got != 1 {
		t.Errorf("FindIndex() = %d, want 1", got)
	}
	if got := xslices.FindIndex(in, func(n int) bool { return n == 9 }); got != -1 {
		t.Errorf("FindIndex() miss = %d, want -1", got)
	}

	got, ok := xslices.FindLast([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
	if !ok || got != 4 {
		t.Errorf("FindLast() = (%v, %v), want (4, true)", got, ok)
	}
}

func TestSomeEveryNone(t *testing.T) {
	in := []int{2, 4, 6}
	isEven := func(n int) bool { return n%2 == 0 }
	isOdd := func(n int) bool { return n%2 != 0 }

	if !xslices.Some(in, isEven) {
		t.Error("Some() = false, want true")
	}
	if !xslices.Every(in, isEven) {
		t.Error("Every() = false, want true")
	}
	if !xslices.None(in, isOdd) {
		t.Error("None() = false, want true")
	}
	// Every is vacuously true on an empty slice.
	if !xslices.Every(nil, isOdd) {
		t.Error("Every(nil) = false, want true")
	}
}

func TestCount(t *testing.T) {
	got := xslices.Count([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
	if got != 2 {
		t.Errorf("Count() = %d, want 2", got)
	}
}

func TestGroupBy(t *testing.T) {
	got := xslices.GroupBy(people, func(p person) int { return p.Age })
	want := map[int][]person{
		30: {{"Alice", 30}},
		25: {{"Bob", 25}, {"Dave", 25}},
		35: {{"Carol", 35}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GroupBy() = %v, want %v", got, want)
	}
}

func TestPartition(t *testing.T) {
	matched, unmatched := xslices.Partition([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
	if want := []int{2, 4}; !reflect.DeepEqual(matched, want) {
		t.Errorf("Partition() matched = %v, want %v", matched, want)
	}
	if want := []int{1, 3}; !reflect.DeepEqual(unmatched, want) {
		t.Errorf("Partition() unmatched = %v, want %v", unmatched, want)
	}
}

func TestKeyBy(t *testing.T) {
	got := xslices.KeyBy(people, func(p person) string { return p.Name })
	if got["Alice"].Age != 30 {
		t.Errorf("KeyBy()[Alice].Age = %d, want 30", got["Alice"].Age)
	}
	if len(got) != 4 {
		t.Errorf("KeyBy() len = %d, want 4", len(got))
	}
}

func TestKeyByLastWins(t *testing.T) {
	got := xslices.KeyBy(people, func(p person) int { return p.Age })
	if got[25].Name != "Dave" {
		t.Errorf("KeyBy()[25].Name = %q, want %q", got[25].Name, "Dave")
	}
}

func TestSumAndSumBy(t *testing.T) {
	if got := xslices.Sum([]int{1, 2, 3}); got != 6 {
		t.Errorf("Sum() = %d, want 6", got)
	}
	if got := xslices.Sum([]float64{}); got != 0 {
		t.Errorf("Sum(empty) = %v, want 0", got)
	}
	if got := xslices.SumBy(people, func(p person) int { return p.Age }); got != 115 {
		t.Errorf("SumBy() = %d, want 115", got)
	}
}

func TestMaxByAndMinBy(t *testing.T) {
	oldest, ok := xslices.MaxBy(people, func(p person) int { return p.Age })
	if !ok || oldest.Name != "Carol" {
		t.Errorf("MaxBy() = (%v, %v), want Carol", oldest, ok)
	}

	youngest, ok := xslices.MinBy(people, func(p person) int { return p.Age })
	if !ok || youngest.Name != "Bob" {
		t.Errorf("MinBy() = (%v, %v), want Bob (first of the tie)", youngest, ok)
	}

	if _, ok := xslices.MaxBy(nil, func(p person) int { return p.Age }); ok {
		t.Error("MaxBy(nil) ok = true, want false")
	}
	if _, ok := xslices.MinBy(nil, func(p person) int { return p.Age }); ok {
		t.Error("MinBy(nil) ok = true, want false")
	}
}
