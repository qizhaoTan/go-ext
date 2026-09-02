package sets_test

import (
	"fmt"
	"sort"

	"github.com/qizhaoTan/go-ext/sets"
)

func ExampleSet() {
	s := sets.Of(1, 2, 3)
	s.Add(3, 4) // 3 is already present

	fmt.Println(s.Len(), s.Contains(4))
	// Output: 4 true
}

func ExampleSet_Intersection() {
	a := sets.Of("go", "rust", "zig")
	b := sets.Of("go", "python", "zig")

	common := a.Intersection(b).ToSlice()
	sort.Strings(common) // iteration order is unspecified
	fmt.Println(common)
	// Output: [go zig]
}
