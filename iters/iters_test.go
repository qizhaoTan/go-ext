package iters_test

import (
	"reflect"
	"testing"

	"github.com/qizhaoTan/go-ext/iters"
)

func TestOfAndCollect(t *testing.T) {
	got := iters.Collect(iters.Of(1, 2, 3))
	if want := []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
		t.Errorf("Collect(Of()) = %v, want %v", got, want)
	}
	if got := iters.Collect(iters.Of[int]()); got != nil {
		t.Errorf("Collect(empty) = %v, want nil", got)
	}
}

func TestMapAndFilter(t *testing.T) {
	seq := iters.Map(iters.Of(1, 2, 3), func(n int) int { return n * 2 })
	if got, want := iters.Collect(seq), []int{2, 4, 6}; !reflect.DeepEqual(got, want) {
		t.Errorf("Map() = %v, want %v", got, want)
	}

	seq = iters.Filter(iters.Of(1, 2, 3, 4), func(n int) bool { return n%2 == 0 })
	if got, want := iters.Collect(seq), []int{2, 4}; !reflect.DeepEqual(got, want) {
		t.Errorf("Filter() = %v, want %v", got, want)
	}
}

func TestFilterMap(t *testing.T) {
	seq := iters.FilterMap(iters.Of(1, 2, 3, 4), func(n int) (string, bool) {
		if n%2 != 0 {
			return "", false
		}
		return "even", true
	})
	if got, want := iters.Collect(seq), []string{"even", "even"}; !reflect.DeepEqual(got, want) {
		t.Errorf("FilterMap() = %v, want %v", got, want)
	}
}

func TestTakeAndDrop(t *testing.T) {
	if got, want := iters.Collect(iters.Take(iters.Range(0, 10), 3)), []int{0, 1, 2}; !reflect.DeepEqual(got, want) {
		t.Errorf("Take() = %v, want %v", got, want)
	}
	if got := iters.Collect(iters.Take(iters.Range(0, 10), 0)); got != nil {
		t.Errorf("Take(0) = %v, want nil", got)
	}
	if got := iters.Collect(iters.Take(iters.Range(0, 10), -1)); got != nil {
		t.Errorf("Take(-1) = %v, want nil", got)
	}
	// Asking for more than there is yields everything.
	if got, want := iters.Collect(iters.Take(iters.Of(1, 2), 5)), []int{1, 2}; !reflect.DeepEqual(got, want) {
		t.Errorf("Take(more than available) = %v, want %v", got, want)
	}

	if got, want := iters.Collect(iters.Drop(iters.Range(0, 5), 3)), []int{3, 4}; !reflect.DeepEqual(got, want) {
		t.Errorf("Drop() = %v, want %v", got, want)
	}
}

func TestTakeIsLazyOverAnInfiniteSequence(t *testing.T) {
	// The point of the package: this terminates.
	got := iters.Collect(iters.Take(iters.Repeat("x"), 3))
	if want := []string{"x", "x", "x"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Take(Repeat()) = %v, want %v", got, want)
	}
}

func TestTakeOnlyPullsWhatItNeeds(t *testing.T) {
	pulled := 0
	counting := iters.Map(iters.Range(0, 1000), func(n int) int {
		pulled++
		return n
	})
	iters.Collect(iters.Take(counting, 3))
	if pulled != 3 {
		t.Errorf("pulled %d elements, want 3 — the pipeline is not lazy", pulled)
	}
}

func TestTakeWhileAndDropWhile(t *testing.T) {
	lessThan3 := func(n int) bool { return n < 3 }

	got := iters.Collect(iters.TakeWhile(iters.Of(1, 2, 3, 1), lessThan3))
	if want := []int{1, 2}; !reflect.DeepEqual(got, want) {
		t.Errorf("TakeWhile() = %v, want %v (stops at the first failure)", got, want)
	}

	got = iters.Collect(iters.DropWhile(iters.Of(1, 2, 3, 1), lessThan3))
	if want := []int{3, 1}; !reflect.DeepEqual(got, want) {
		t.Errorf("DropWhile() = %v, want %v (resumes after the first failure)", got, want)
	}
}

func TestConcat(t *testing.T) {
	got := iters.Collect(iters.Concat(iters.Of(1, 2), iters.Of[int](), iters.Of(3)))
	if want := []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
		t.Errorf("Concat() = %v, want %v", got, want)
	}
}

func TestReduceAndCount(t *testing.T) {
	sum := iters.Reduce(iters.Range(1, 5), 0, func(acc, n int) int { return acc + n })
	if sum != 10 {
		t.Errorf("Reduce() = %d, want 10", sum)
	}
	if got := iters.Count(iters.Range(0, 7)); got != 7 {
		t.Errorf("Count() = %d, want 7", got)
	}
}

func TestFindSomeEvery(t *testing.T) {
	got, ok := iters.Find(iters.Range(0, 10), func(n int) bool { return n > 5 })
	if !ok || got != 6 {
		t.Errorf("Find() = (%d, %v), want (6, true)", got, ok)
	}
	if _, ok := iters.Find(iters.Range(0, 3), func(n int) bool { return n > 5 }); ok {
		t.Error("Find() miss ok = true, want false")
	}

	if !iters.Some(iters.Of(1, 2), func(n int) bool { return n == 2 }) {
		t.Error("Some() = false, want true")
	}
	if !iters.Every(iters.Of(2, 4), func(n int) bool { return n%2 == 0 }) {
		t.Error("Every() = false, want true")
	}
	if iters.Every(iters.Of(2, 3), func(n int) bool { return n%2 == 0 }) {
		t.Error("Every() = true, want false")
	}
}

func TestFindStopsEarly(t *testing.T) {
	pulled := 0
	counting := iters.Map(iters.Range(0, 1000), func(n int) int {
		pulled++
		return n
	})
	iters.Find(counting, func(n int) bool { return n == 2 })
	if pulled != 3 {
		t.Errorf("pulled %d elements, want 3 — Find did not stop early", pulled)
	}
}

func TestRangeAndIterate(t *testing.T) {
	if got := iters.Collect(iters.Range(3, 3)); got != nil {
		t.Errorf("Range(3,3) = %v, want nil", got)
	}

	powers := iters.Collect(iters.Take(iters.Iterate(1, func(n int) int { return n * 2 }), 5))
	if want := []int{1, 2, 4, 8, 16}; !reflect.DeepEqual(powers, want) {
		t.Errorf("Iterate() = %v, want %v", powers, want)
	}
}

func TestPipelineComposes(t *testing.T) {
	// Map over Filter over an infinite sequence, bounded by Take.
	seq := iters.Take(
		iters.Map(
			iters.Filter(iters.Iterate(1, func(n int) int { return n + 1 }), func(n int) bool { return n%3 == 0 }),
			func(n int) int { return n * n },
		),
		4,
	)
	if got, want := iters.Collect(seq), []int{9, 36, 81, 144}; !reflect.DeepEqual(got, want) {
		t.Errorf("pipeline = %v, want %v", got, want)
	}
}
