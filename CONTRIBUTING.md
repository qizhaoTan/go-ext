# Contributing

Thanks for taking the time to contribute.

## Getting started

```bash
git clone https://github.com/qizhaoTan/go-ext.git
cd go-ext
make check
```

The library has no dependencies outside the standard library, and that is
deliberate — please do not add one.

## What belongs here

A function earns its place when it is:

- **general** — useful across projects, not specific to one domain;
- **pure** — it does not mutate its arguments, and returns fresh values;
- **predictable** — it behaves sensibly on `nil` and empty inputs;
- **not already in the standard library** — if `slices` or `maps` covers it,
  use those instead.

## Pull requests

1. Add tests. New code should cover the empty and `nil` cases as well as the
   happy path.
2. Add a runnable `Example` for anything whose usage is not obvious — these show
   up on pkg.go.dev.
3. Document every exported identifier, starting the comment with its name.
4. Run `make check` before pushing.
5. Keep the commit history readable; squash fixup commits.

## Compatibility

Anything exported is covered by [semantic versioning](https://semver.org/).
Breaking changes wait for a major release, so new behaviour usually means a new
function rather than a changed signature.
