package maps_test

import (
	"reflect"
	"sort"
	"strconv"
	"testing"

	xmaps "github.com/qizhaoTan/go-ext/maps"
)

func sorted[T ~int | ~string](s []T) []T {
	out := append([]T(nil), s...)
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

var base = map[string]int{"a": 1, "b": 2, "c": 3}

func TestKeysAndValues(t *testing.T) {
	if got, want := sorted(xmaps.Keys(base)), []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Keys() = %v, want %v", got, want)
	}
	if got, want := sorted(xmaps.Values(base)), []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
		t.Errorf("Values() = %v, want %v", got, want)
	}
	if got := xmaps.Keys(map[string]int{}); len(got) != 0 {
		t.Errorf("Keys(empty) = %v, want empty", got)
	}
}

func TestEntriesRoundTrip(t *testing.T) {
	got := xmaps.FromEntries(xmaps.Entries(base))
	if !reflect.DeepEqual(got, base) {
		t.Errorf("round trip = %v, want %v", got, base)
	}
}

func TestFilter(t *testing.T) {
	got := xmaps.Filter(base, func(_ string, v int) bool { return v > 1 })
	want := map[string]int{"b": 2, "c": 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filter() = %v, want %v", got, want)
	}
}

func TestMapValuesAndMapKeys(t *testing.T) {
	values := xmaps.MapValues(base, func(_ string, v int) string { return strconv.Itoa(v) })
	want := map[string]string{"a": "1", "b": "2", "c": "3"}
	if !reflect.DeepEqual(values, want) {
		t.Errorf("MapValues() = %v, want %v", values, want)
	}

	keys := xmaps.MapKeys(base, func(k string, _ int) string { return k + "!" })
	if len(keys) != 3 || keys["a!"] != 1 {
		t.Errorf("MapKeys() = %v, want keys suffixed with !", keys)
	}
}

func TestInvert(t *testing.T) {
	got := xmaps.Invert(base)
	want := map[int]string{1: "a", 2: "b", 3: "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Invert() = %v, want %v", got, want)
	}
}

func TestMergeLaterWins(t *testing.T) {
	got := xmaps.Merge(
		map[string]int{"a": 1, "b": 2},
		map[string]int{"b": 20, "c": 3},
	)
	want := map[string]int{"a": 1, "b": 20, "c": 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Merge() = %v, want %v", got, want)
	}
}

func TestMergeDoesNotMutateInputs(t *testing.T) {
	first := map[string]int{"a": 1}
	xmaps.Merge(first, map[string]int{"a": 2})
	if first["a"] != 1 {
		t.Errorf("Merge mutated its first input: %v", first)
	}
}

func TestGetOr(t *testing.T) {
	if got := xmaps.GetOr(base, "a", -1); got != 1 {
		t.Errorf("GetOr(hit) = %d, want 1", got)
	}
	if got := xmaps.GetOr(base, "zzz", -1); got != -1 {
		t.Errorf("GetOr(miss) = %d, want -1", got)
	}
}
