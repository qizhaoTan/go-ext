# go-ext

Generic data structures and functional helpers for Go. Zero dependencies
outside the standard library — keep it that way.

## Layout

| Package | Role |
| --- | --- |
| `slices/` | Eager functional helpers. `transform.go` returns new collections; `query.go` inspects and aggregates. |
| `iters/` | Lazy pipelines over `iter.Seq`. The differentiator: composes without allocating, handles infinite sequences. |
| `sets/` | Generic hash set over `map[T]struct{}`. The only stateful type here. |
| `maps/` | Map transformations. |
| `ptr/` | Value ↔ pointer conversion. |
| `constraints/` | Numeric constraints. Deliberately no `Ordered` — `cmp.Ordered` exists. |

## Invariants

These hold across the library; breaking one is a bug, not a style choice.

- **Purity.** Nothing mutates an argument. Every returned collection is freshly
  allocated and does not alias the input (`Chunk` copies for this reason).
- **`nil` is empty.** A `nil` slice or map ranges zero times, so it flows
  through without a guard. Never add a `nil` special case.
- **No standard-library duplication.** Sorting, searching and equality belong to
  the standard `slices`/`maps`. This library only fills gaps.
- **Free functions, not methods.** Go cannot introduce type parameters on
  methods, so a chained API could not stay generic over the element type.
- **Panic on caller bugs.** A non-positive `Chunk` size has no correct answer;
  panicking beats inventing one.

## Conventions

- Every exported identifier is documented, comment starting with its name.
- Tests are external (`package foo_test`) — they exercise the public API only.
- Import the local `slices`/`maps` as `xslices`/`xmaps` in tests to avoid
  shadowing the standard library.
- Cover the empty and `nil` cases alongside the happy path.
- Add a runnable `Example` for anything non-obvious; they render on pkg.go.dev
  and are verified by `go test`.

## Commands

```bash
make check   # fmt + vet + test — run before every commit
make test    # go test -race ./...
make cover   # coverage report at coverage.html
make lint    # golangci-lint
```
