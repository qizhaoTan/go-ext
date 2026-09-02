package slices_test

import (
	"fmt"
	"sort"
	"strings"

	xslices "github.com/qizhaoTan/go-ext/slices"
)

func ExampleMap() {
	names := []string{"alice", "bob"}
	fmt.Println(xslices.Map(names, strings.ToUpper))
	// Output: [ALICE BOB]
}

func ExampleFilter() {
	fmt.Println(xslices.Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 }))
	// Output: [2 4]
}

func ExampleGroupBy() {
	words := []string{"apple", "avocado", "banana"}
	groups := xslices.GroupBy(words, func(s string) string { return s[:1] })

	// Map iteration order is unspecified, so sort the keys for a stable output.
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %v\n", k, groups[k])
	}
	// Output:
	// a: [apple avocado]
	// b: [banana]
}

func ExampleReduce() {
	total := xslices.Reduce([]int{1, 2, 3, 4}, 0, func(acc, n int) int { return acc + n })
	fmt.Println(total)
	// Output: 10
}

func ExampleChunk() {
	fmt.Println(xslices.Chunk([]int{1, 2, 3, 4, 5}, 2))
	// Output: [[1 2] [3 4] [5]]
}
