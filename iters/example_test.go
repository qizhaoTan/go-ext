package iters_test

import (
	"fmt"

	"github.com/qizhaoTan/go-ext/iters"
)

// The pipeline is lazy, so it does only the work the Take actually demands,
// even though the underlying sequence is infinite.
func ExampleTake() {
	naturals := iters.Iterate(1, func(n int) int { return n + 1 })
	squaresOfMultiplesOfThree := iters.Map(
		iters.Filter(naturals, func(n int) bool { return n%3 == 0 }),
		func(n int) int { return n * n },
	)
	fmt.Println(iters.Collect(iters.Take(squaresOfMultiplesOfThree, 4)))
	// Output: [9 36 81 144]
}

func ExampleFilter() {
	evens := iters.Filter(iters.Range(0, 10), func(n int) bool { return n%2 == 0 })
	for v := range evens {
		fmt.Print(v, " ")
	}
	fmt.Println()
	// Output: 0 2 4 6 8
}

func ExampleReduce() {
	sum := iters.Reduce(iters.Range(1, 101), 0, func(acc, n int) int { return acc + n })
	fmt.Println(sum)
	// Output: 5050
}
