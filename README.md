# some re-usable snippets

Think of this as my collection of snippets in a go gettable way

```
go get go.livingit.de/reuse
```

> **Module renamed:** this module was previously `github.com/sascha-andres/reuse`.
> It is now `go.livingit.de/reuse`. Update any existing imports and your
> `go.mod` requirement accordingly.

## Releases

- [v0.17.1](https://github.com/sascha-andres/reuse/releases/tag/v0.17.1) — flag: stop dropping flags after, between, and inline with verbs
- [v0.17.0](https://github.com/sascha-andres/reuse/releases/tag/v0.17.0) — namespace change
- [v0.16.0](https://github.com/sascha-andres/reuse/releases/tag/v0.16.0) — Struct config rework (#9)
- [v0.15.1](https://github.com/sascha-andres/reuse/releases/tag/v0.15.1) — feat: test: do not re-add first
- [v0.15.0](https://github.com/sascha-andres/reuse/releases/tag/v0.15.0) — add stringslice (#7)
- [v0.14.2](https://github.com/sascha-andres/reuse/releases/tag/v0.14.2) — fix: re-add first argument
- [v0.14.1](https://github.com/sascha-andres/reuse/releases/tag/v0.14.1) — feat: add Args()
- [v0.14.0](https://github.com/sascha-andres/reuse/releases/tag/v0.14.0) — feat: add package j
- [v0.13.0](https://github.com/sascha-andres/reuse/releases/tag/v0.13.0) — feat: add ErrNotImplemented
- [v0.12.0](https://github.com/sascha-andres/reuse/releases/tag/v0.12.0) — Merge pull request #6 from sascha-andres/async-await
- [v0.11.1](https://github.com/sascha-andres/reuse/releases/tag/v0.11.1) — feat: break loop if flag found after verb
- [v0.11.0](https://github.com/sascha-andres/reuse/releases/tag/v0.11.0) — feat: add test with struct
- [v0.10.0](https://github.com/sascha-andres/reuse/releases/tag/v0.10.0) — add backoff (#5)
- [v0.9.1](https://github.com/sascha-andres/reuse/releases/tag/v0.9.1) — fix(flag): use correct env prefix for reading env var
- [v0.9.0](https://github.com/sascha-andres/reuse/releases/tag/v0.9.0) — feat(flag): allow override prefix for selected flags
- [v0.8.1](https://github.com/sascha-andres/reuse/releases/tag/v0.8.1) — feat(lo): add lo package
- [v0.8.0](https://github.com/sascha-andres/reuse/releases/tag/v0.8.0) — Merge pull request #4 from sascha-andres/add-testing-package
- [v0.7.0](https://github.com/sascha-andres/reuse/releases/tag/v0.7.0) — feat(circuitbreaker): initial version
- [v0.6.2](https://github.com/sascha-andres/reuse/releases/tag/v0.6.2) — feat: add exists helper for dir and file
- [v0.6.1](https://github.com/sascha-andres/reuse/releases/tag/v0.6.1) — feat: add test for filter
- [v0.6.0](https://github.com/sascha-andres/reuse/releases/tag/v0.6.0) — feat: add filter func
- [v0.5.2](https://github.com/sascha-andres/reuse/releases/tag/v0.5.2) — fix: correct function call
- [v0.5.1](https://github.com/sascha-andres/reuse/releases/tag/v0.5.1) — fix: pass values to UintVarWithoutEnv
- [v0.5.0](https://github.com/sascha-andres/reuse/releases/tag/v0.5.0) — feat: add struct flag config
- [v0.4.0](https://github.com/sascha-andres/reuse/releases/tag/v0.4.0) — feat: add date time time functions
- [v0.3.0](https://github.com/sascha-andres/reuse/releases/tag/v0.3.0) — fix: use correct method in example
- [v0.2.0](https://github.com/sascha-andres/reuse/releases/tag/v0.2.0) — feat: add group by func
- [v0.1.0](https://github.com/sascha-andres/reuse/releases/tag/v0.1.0) — feat: add separate feature
- [v0.0.6](https://github.com/sascha-andres/reuse/releases/tag/v0.0.6) — feat: prevent first error acurring twice
- [v0.0.5](https://github.com/sascha-andres/reuse/releases/tag/v0.0.5) — feat: fix verb detection
- [v0.0.4](https://github.com/sascha-andres/reuse/releases/tag/v0.0.4) — feat: add flag package
- [v0.0.3](https://github.com/sascha-andres/reuse/releases/tag/v0.0.3) — feat: update dependencies
- [v0.0.2](https://github.com/sascha-andres/reuse/releases/tag/v0.0.2) — feat: update dependencies
- [v0.0.1](https://github.com/sascha-andres/reuse/releases/tag/v0.0.1) — feat: add memoize function

## Packages

- **[reuse](.)** (root) — grab-bag of standalone helpers: `Memoize`,
  retry/`Do`/`DoCtx` with configurable backoff delays, date/time
  formatting helpers, `FileExists`/`DirectoryExists`, `DiscardError` for
  ignoring a `Close()`-style error deliberately, `MultiError`.
- **[async](async/)** — `Future[T]`/`Async` for running a function in a
  goroutine and awaiting its result or error. See [async/README.md](async/README.md).
- **[circuitbreaker](circuitbreaker/)** — `CircuitBreaker` with
  closed/open/half-open states around a `WorkFunc`.
- **[flag](flag/)** — drop-in replacement for stdlib `flag` adding verb
  (positional-argument) collection interleaved with flags, environment
  variable fallback/prefixing, `AddFlagsForStruct` for binding a struct's
  fields to flags via `flag:"..."` tags, and `*VarP` constructors that
  resolve to `nil` unless explicitly provided. See [flag/README.md](flag/README.md).
- **[functional](functional/)** — generic slice/map helpers: `Filter`,
  `Map`, `GroupBy`/`GroupByFunc`, `Head`/`Tail`, plus `ValueCell`/
  `FormulaCell` for a small reactive-cell primitive.
- **[j](j/)** — `UnmarshalFile` reads and JSON-decodes a file into `T`.
- **[lo](lo/)** — `slog`-logging variants of the "must" pattern:
  `Must`, `MustGet`, `MustRun`.
- **[te](te/)** — `testing`-failing variants of the same pattern:
  `MustGet`, `MustRun`, `ShouldGet`, `ShouldRun`.
- **[validations](validations/)** — `CreditCard` (Luhn check).

## Plans

- [0001 — struct_flag: support nested structs](plans/0001-struct-flag-nested-support.md)
- [0002 — struct_flag: support pointer-to-scalar leaf fields](plans/0002-struct-flag-pointer-scalar-leaves.md)
- [0003 — struct_flag: support time.Duration](plans/0003-struct-flag-duration-support.md)
- [0004 — flag: add *VarP constructors](plans/0004-flag-varp-constructors.md)
- [0005 — rename module namespace to go.livingit.de/reuse](plans/0005-module-namespace-rename.md)
- [0006 — flag: stop dropping flags that appear after or between verbs](plans/0006-flag-parse-flags-after-verbs.md)