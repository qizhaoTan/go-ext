package slices_test

import (
	"reflect"
	"strconv"
	"testing"

	xslices "github.com/qizhaoTan/go-ext/slices"
)

func TestMap(t *testing.T) {
	got := xslices.Map([]int{1, 2, 3}, strconv.Itoa)
	want := []string{"1", "2", "3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Map() = %v, want %v", got, want)
	}
	if got := xslices.Map(nil, strconv.Itoa); len(got) != 0 {
		t.Errorf("Map(nil) = %v, want empty", got)
	}
}

func TestMapIndex(t *testing.T) {
	got := xslices.MapIndex([]string{"a", "b"}, func(i int, s string) string {
		return strconv.Itoa(i) + s
	})
	want := []string{"0a", "1b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MapIndex() = %v, want %v", got, want)
	}
}

func TestFilterAndReject(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	in := []int{1, 2, 3, 4, 5, 6}

	if got, want := xslices.Filter(in, isEven), []int{2, 4, 6}; !reflect.DeepEqual(got, want) {
		t.Errorf("Filter() = %v, want %v", got, want)
	}
	if got, want := xslices.Reject(in, isEven), []int{1, 3, 5}; !reflect.DeepEqual(got, want) {
		t.Errorf("Reject() = %v, want %v", got, want)
	}
}

func TestFilterDoesNotMutateInput(t *testing.T) {
	in := []int{1, 2, 3, 4}
	out := xslices.Filter(in, func(n int) bool { return n > 2 })
	out[0] = 99
	if want := []int{1, 2, 3, 4}; !reflect.DeepEqual(in, want) {
		t.Errorf("input mutated: got %v, want %v", in, want)
	}
}

func TestFilterMap(t *testing.T) {
	got := xslices.FilterMap([]string{"1", "x", "3"}, func(s string) (int, bool) {
		n, err := strconv.Atoi(s)
		return n, err == nil
	})
	if want := []int{1, 3}; !reflect.DeepEqual(got, want) {
		t.Errorf("FilterMap() = %v, want %v", got, want)
	}
}

func TestFlatMapAndFlatten(t *testing.T) {
	got := xslices.FlatMap([]int{1, 2}, func(n int) []int { return []int{n, n * 10} })
	if want := []int{1, 10, 2, 20}; !reflect.DeepEqual(got, want) {
		t.Errorf("FlatMap() = %v, want %v", got, want)
	}

	flat := xslices.Flatten([][]int{{1, 2}, nil, {3}})
	if want := []int{1, 2, 3}; !reflect.DeepEqual(flat, want) {
		t.Errorf("Flatten() = %v, want %v", flat, want)
	}
}

func TestReduceAndReduceRight(t *testing.T) {
	concat := func(acc string, s string) string { return acc + s }
	in := []string{"a", "b", "c"}

	if got := xslices.Reduce(in, "", concat); got != "abc" {
		t.Errorf("Reduce() = %q, want %q", got, "abc")
	}
	if got := xslices.ReduceRight(in, "", concat); got != "cba" {
		t.Errorf("ReduceRight() = %q, want %q", got, "cba")
	}
}

func TestReverse(t *testing.T) {
	in := []int{1, 2, 3}
	got := xslices.Reverse(in)
	if want := []int{3, 2, 1}; !reflect.DeepEqual(got, want) {
		t.Errorf("Reverse() = %v, want %v", got, want)
	}
	if want := []int{1, 2, 3}; !reflect.DeepEqual(in, want) {
		t.Errorf("Reverse mutated input: got %v, want %v", in, want)
	}
}

func TestUniqAndUniqBy(t *testing.T) {
	if got, want := xslices.Uniq([]int{1, 2, 1, 3, 2}), []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
		t.Errorf("Uniq() = %v, want %v", got, want)
	}

	got := xslices.UniqBy([]string{"apple", "avocado", "banana"}, func(s string) byte { return s[0] })
	if want := []string{"apple", "banana"}; !reflect.DeepEqual(got, want) {
		t.Errorf("UniqBy() = %v, want %v", got, want)
	}
}

func TestChunk(t *testing.T) {
	got := xslices.Chunk([]int{1, 2, 3, 4, 5}, 2)
	want := [][]int{{1, 2}, {3, 4}, {5}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Chunk() = %v, want %v", got, want)
	}
}

func TestChunkDoesNotAliasInput(t *testing.T) {
	in := []int{1, 2, 3, 4}
	chunks := xslices.Chunk(in, 2)
	chunks[0][0] = 99
	if want := []int{1, 2, 3, 4}; !reflect.DeepEqual(in, want) {
		t.Errorf("Chunk aliased input: got %v, want %v", in, want)
	}
}

func TestChunkPanicsOnNonPositiveSize(t *testing.T) {
	for _, size := range []int{0, -1} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf("Chunk(_, %d) did not panic", size)
				}
			}()
			xslices.Chunk([]int{1, 2}, size)
		})
	}
}

func TestCompact(t *testing.T) {
	if got, want := xslices.Compact([]int{0, 1, 0, 2}), []int{1, 2}; !reflect.DeepEqual(got, want) {
		t.Errorf("Compact() = %v, want %v", got, want)
	}
	if got, want := xslices.Compact([]string{"", "a"}), []string{"a"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Compact() = %v, want %v", got, want)
	}
}
