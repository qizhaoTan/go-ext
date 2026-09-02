# go-ext

[![CI](https://github.com/qizhaoTan/go-ext/actions/workflows/ci.yml/badge.svg)](https://github.com/qizhaoTan/go-ext/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/qizhaoTan/go-ext.svg)](https://pkg.go.dev/github.com/qizhaoTan/go-ext)
[![Go Report Card](https://goreportcard.com/badge/github.com/qizhaoTan/go-ext)](https://goreportcard.com/report/github.com/qizhaoTan/go-ext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Generic data structures and functional helpers for Go — the parts the standard
library leaves out.

Zero dependencies. Everything outside `sets` is a pure function: nothing mutates
its input, `nil` is always treated as empty, and every returned collection is
freshly allocated.

```bash
go get github.com/qizhaoTan/go-ext
```

Requires Go 1.24 or later.

## Packages

| Package | What it is |
| --- | --- |
| [`slices`](https://pkg.go.dev/github.com/qizhaoTan/go-ext/slices) | Eager functional helpers over slices — `Map`, `Filter`, `Reduce`, `GroupBy`, `Chunk`, … |
| [`iters`](https://pkg.go.dev/github.com/qizhaoTan/go-ext/iters) | **Lazy** pipelines over `iter.Seq` — composes without allocating, works on infinite sequences |
| [`sets`](https://pkg.go.dev/github.com/qizhaoTan/go-ext/sets) | A generic hash set with the usual algebra |
| [`maps`](https://pkg.go.dev/github.com/qizhaoTan/go-ext/maps) | `Keys`, `Values`, `Filter`, `Invert`, `Merge`, … |
| [`ptr`](https://pkg.go.dev/github.com/qizhaoTan/go-ext/ptr) | Value ↔ pointer conversion for optional fields |
| [`constraints`](https://pkg.go.dev/github.com/qizhaoTan/go-ext/constraints) | Numeric type constraints |

## Eager or lazy

The two collection packages are the same operations with different evaluation
strategies, and picking between them is the main decision this library asks you
to make.

`slices` is **eager**. Each step allocates a new slice and runs to completion.
That is the right default: it is simple, it is fast on small inputs, and the
intermediate values are ordinary slices you can print in a debugger.

```go
import xslices "github.com/qizhaoTan/go-ext/slices"

users := []User{ /* ... */ }

names := xslices.Map(
    xslices.Filter(users, func(u User) bool { return u.Active }),
    func(u User) string { return u.Name },
)
```

`iters` is **lazy**. Nothing runs until the sequence is consumed, no
intermediate slice is ever built, and only the elements actually demanded are
touched. That is what makes this terminate:

```go
import "github.com/qizhaoTan/go-ext/iters"

naturals := iters.Iterate(1, func(n int) int { return n + 1 }) // infinite

first4 := iters.Collect(iters.Take(
    iters.Map(
        iters.Filter(naturals, func(n int) bool { return n%3 == 0 }),
        func(n int) int { return n * n },
    ),
    4,
))
// [9 36 81 144] — the pipeline pulled 12 elements, not 1000
```

Reach for `iters` when the source is large, expensive, or unbounded, and for
`slices` the rest of the time.

## Examples

### Grouping and aggregating

```go
byDept := xslices.GroupBy(employees, func(e Employee) string { return e.Dept })
payroll := xslices.SumBy(employees, func(e Employee) int { return e.Salary })
senior, junior := xslices.Partition(employees, func(e Employee) bool { return e.Years > 5 })
```

### Sets

```go
import "github.com/qizhaoTan/go-ext/sets"

current := sets.Of("read", "write")
granted := sets.Of("read", "admin")

missing := granted.Difference(current) // Set{admin}
if current.IsSubsetOf(granted) { /* ... */ }

for perm := range current.All() { // iter.Seq, works with range
    fmt.Println(perm)
}
```

A `Set` is a wrapper over `map[T]struct{}` and is **not** safe for concurrent
use. Guard it with a mutex, or give each goroutine its own and merge with
`Union`.

### Optional fields

```go
import "github.com/qizhaoTan/go-ext/ptr"

req := UpdateRequest{
    Name:  ptr.To("alice"),      // &"alice" is not valid Go
    Email: ptr.ToNilable(email), // nil when email is ""
}

name := ptr.FromOr(resp.Name, "(unknown)")
```

## Design notes

**Only what the standard library omits.** Sorting, searching and equality live
in the standard `slices` and `maps` packages; this library does not duplicate
them. `constraints` deliberately omits `Ordered`, because `cmp.Ordered` exists.

**Errors over silent fallbacks.** `Chunk` panics on a non-positive size rather
than returning something arbitrary — a zero chunk size is a bug at the call
site, and hiding it costs more than it saves.

**No method chaining.** Go's type parameters cannot be introduced on methods, so
a fluent `.Map().Filter()` API would have to fix the element type at
construction and could never change it. Free functions keep the transformations
fully generic, which is why the standard library made the same choice.

**No `nil` special cases in the API.** A `nil` slice or map ranges zero times,
so it flows through every function as an empty one without a guard.

## Development

```bash
make test    # go test -race ./...
make cover   # coverage report at coverage.html
make lint    # golangci-lint
make check   # fmt + vet + test
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Bug reports and pull requests are
welcome.

## License

[MIT](LICENSE) © qizhaoTan
