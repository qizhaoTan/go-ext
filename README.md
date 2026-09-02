# go-ext

[English](README.md) · [简体中文](README.zh-CN.md)

[![CI](https://github.com/qizhaoTan/go-ext/actions/workflows/ci.yml/badge.svg)](https://github.com/qizhaoTan/go-ext/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/qizhaoTan/go-ext.svg)](https://pkg.go.dev/github.com/qizhaoTan/go-ext)
[![Go Report Card](https://goreportcard.com/badge/github.com/qizhaoTan/go-ext)](https://goreportcard.com/report/github.com/qizhaoTan/go-ext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Generic data structures and functional helpers for Go — the parts the standard
library leaves out.

> **Status: early development.** The repository is scaffolded and the design is
> settled, but the packages below are not implemented yet. Nothing here is
> importable, and the API may still change. Watch the repo if you want to know
> when the first release lands.

## Planned packages

| Package | Role |
| --- | --- |
| `slices` | Eager functional helpers over slices — `Map`, `Filter`, `Reduce`, `GroupBy`, `Chunk`, … |
| `iters` | **Lazy** pipelines over `iter.Seq` — composes without allocating, works on infinite sequences |
| `sets` | A generic hash set with the usual algebra |
| `maps` | `Keys`, `Values`, `Filter`, `Invert`, `Merge`, … |
| `ptr` | Value ↔ pointer conversion for optional fields |
| `constraints` | Numeric type constraints |

## Design goals

These are the constraints the implementation is being written against.

**Zero dependencies.** Nothing outside the standard library, ever. A utility
package that drags in a dependency tree is worse than the code it replaces.

**Purity.** Nothing mutates its arguments. Every returned collection is freshly
allocated and does not alias the input, so a result can be handed to a caller
without surprises.

**`nil` is empty.** A `nil` slice or map ranges zero times, so it flows through
every function as an empty one — no guards, no special cases in the API.

**No standard-library duplication.** Sorting, searching and equality already
live in the standard `slices` and `maps` packages. This library only fills the
gaps, which is also why `constraints` will not define `Ordered` — `cmp.Ordered`
exists.

**Free functions, not method chains.** Go cannot introduce new type parameters
on methods, so a fluent `.Map().Filter()` API would have to fix the element type
at construction and could never transform it. The standard library made the
same choice for the same reason.

**Panic on caller bugs.** A non-positive chunk size has no correct answer.
Returning something arbitrary hides a bug at the call site; panicking does not.

## Eager or lazy

The two collection packages will offer the same operations under different
evaluation strategies, and choosing between them is the main decision this
library asks of its users.

`slices` is **eager**: each step allocates a new slice and runs to completion.
That is the right default — simple, fast on small inputs, and the intermediate
values are ordinary slices you can inspect in a debugger.

`iters` is **lazy**: nothing runs until the sequence is consumed, no
intermediate slice is ever built, and only the elements actually demanded are
touched. That is what lets a pipeline over an infinite sequence terminate.

The rule of thumb: reach for `iters` when the source is large, expensive, or
unbounded, and for `slices` the rest of the time.

## Requirements

Go 1.24 or later. The `iters` package is built on `iter.Seq`, which is why the
floor is not lower.

## Development

```bash
git clone https://github.com/qizhaoTan/go-ext.git
cd go-ext
make check   # fmt + vet + test
```

| Target | What it does |
| --- | --- |
| `make test` | `go test -race ./...` |
| `make cover` | Coverage report at `coverage.html` |
| `make lint` | `golangci-lint run` |
| `make fmt` | `gofmt -w .` |
| `make check` | fmt + vet + test — run before every commit |

CI runs the tests with the race detector on Linux, macOS and Windows, against
both Go 1.24 and the current stable release.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). While the library is still being built
out, opening an issue before a large pull request will save you time.

## License

[MIT](LICENSE) © qizhaoTan
