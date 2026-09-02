# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `slices` — functional helpers over slices: `Map`, `MapIndex`, `Filter`,
  `Reject`, `FilterMap`, `FlatMap`, `Flatten`, `Reduce`, `ReduceRight`,
  `Reverse`, `Uniq`, `UniqBy`, `Chunk`, `Compact`, `Find`, `FindIndex`,
  `FindLast`, `Some`, `Every`, `None`, `Count`, `GroupBy`, `Partition`,
  `KeyBy`, `Sum`, `SumBy`, `MaxBy`, `MinBy`.
- `maps` — map helpers: `Keys`, `Values`, `Entries`, `FromEntries`, `Filter`,
  `MapValues`, `MapKeys`, `Invert`, `Merge`, `GetOr`.
- `sets` — a generic hash set with the usual algebra: `Union`, `Intersection`,
  `Difference`, `SymmetricDifference`, subset and superset tests.
- `iters` — lazy pipelines over `iter.Seq`: `Map`, `Filter`, `FilterMap`,
  `Take`, `Drop`, `TakeWhile`, `DropWhile`, `Concat`, `Collect`, `Reduce`,
  `Count`, `Find`, `Some`, `Every`, `Of`, `Range`, `Repeat`, `Iterate`.
- `ptr` — value/pointer conversion: `To`, `From`, `FromOr`, `ToNilable`,
  `IsNil`.
- `constraints` — numeric type constraints: `Signed`, `Unsigned`, `Integer`,
  `Float`, `Complex`, `Number`.

[Unreleased]: https://github.com/qizhaoTan/go-ext/compare/main...HEAD
