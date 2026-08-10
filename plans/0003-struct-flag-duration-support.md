# 0003 — struct_flag: support time.Duration

## Scope

`flag` submodule only (`flag/flag.go`, `flag/struct_flag.go`, their tests,
`flag/README.md`). No other module in this repo is touched.

## Problem

`time.Duration` is defined as `type Duration int64`, so `reflect.Kind()` of
a `time.Duration` field is `reflect.Int64` — indistinguishable, by Kind
alone, from a plain `int64` field. `struct_flag.go`'s registration switch
keys off `Kind()`, so a field like

```go
type Config struct {
    Timeout time.Duration `flag:"timeout"`
}
```

is already silently caught by `isIntKind` today and registered via
`Int64(name, 0, usage)`. It "works" in the sense that it compiles and
builds a flag, but that flag expects a raw nanosecond count
(`-timeout 5000000000`), not a duration string (`-timeout 5s`) — the
opposite of what `time.Duration` fields are for, and inconsistent with
`flag.go`'s own `Duration`/`DurationVar` constructors, which already do
proper `time.ParseDuration`-based parsing for flat (non-struct) use and
have done since before this rework. Both prior plans (`0001`, `0002`)
carried this gap forward without calling it out.

## Goal

A `time.Duration` leaf field (value or pointer) gets duration-string
parsing (`"5s"`, `"1h30m"`, ...) via `flag.go`'s existing `Duration`
machinery, at any nesting depth, instead of being treated as a plain
`int64`:

```go
type Config struct {
    Timeout  time.Duration  `flag:"timeout"`
    Deadline *time.Duration `flag:"deadline"`
}
```

`Timeout` is always populated (parsed from `-timeout`/`TIMEOUT`, defaulting
to `0` like any other value-typed leaf). `Deadline` follows the same
nil-unless-provided rule as every other pointer-scalar leaf from `0002`:
`nil` unless `-deadline` or `DEADLINE` was explicitly given, regardless of
nesting depth or whether the containing branch is itself a value or a
pointer.

## Design

**`flag.go`**: add `DurationP(name, usage string) *time.Duration`,
mirroring `StringP`/`Int64P`/`BoolP`/`Float64P`/`Uint64P` from `0002` — a
thin wrapper (`Duration(name, 0, usage)`) with no default parameter, for
the same reason: a nilable pointer leaf has no meaningful custom default.

**`struct_flag.go`**:

- `Container[T]` gains a sixth map, `durationVariables map[string]*time.Duration`.
  It can't reuse `intVariables` (`map[string]*int64`) — `*time.Duration`
  and `*int64` are distinct pointer types even though `time.Duration`'s
  underlying representation is `int64`.
- `structNode` gains `isDuration bool`.
- In `buildNodes`, **before** the existing `isIntKind`/`String`/`Bool`/
  `isFloatKind`/`isUintKind` switch (in both the pointer and non-pointer
  leaf branches), check `ft == reflect.TypeFor[time.Duration]()` first and
  register via `Duration`/`DurationP` into `durationVariables`, setting
  `isDuration: true` on the node. This ordering matters: `time.Duration`
  would otherwise be caught by `isIntKind` first, since `Kind()` can't
  tell them apart — only an exact type comparison can. A field of some
  other `int64`-based named type (not literally `time.Duration`) still
  falls through to the plain integer path, unchanged.
- `writeLeafScalar` gains an `n.isDuration` case, writing
  `v.SetInt(int64(*c.durationVariables[n.fullName]))` — `SetInt` only
  requires `v.Kind() == reflect.Int64`, which holds whether `v` is a plain
  `int64` field or a `time.Duration` field, so this works unchanged for
  both the value-leaf and pointer-leaf (scratch-value) cases already
  wired up by `populateNode`.
- Presence detection (`setFlags` / `envNameForFlagName` lookup) is
  unchanged — a duration leaf is just another leaf kind as far as
  `populateNode`'s pointer/value branching is concerned.

## Files touched

- `flag/flag.go` — `DurationP`.
- `flag/struct_flag.go` — duration detection/registration/writing as above.
- `flag/flag_test.go` — coverage for `DurationP` (zero default, explicit
  value parses a duration string, `Visit` reports it set).
- `flag/struct_flag_test.go` — value-typed duration leaf (regression that
  it now parses `"5s"` rather than requiring a raw nanosecond count),
  pointer-typed duration leaf present via flag, present via env, absent
  (nil), and a duration leaf nested inside both a value and a pointer
  branch (mirroring the `0002` pointer-scalar nesting tests).
- `flag/README.md` — document duration leaf support.

## Non-goals (explicitly out of scope)

- Named types derived from `time.Duration` (e.g. `type Backoff
  time.Duration`) — only the exact `time.Duration` type is special-cased;
  a derived type falls through to the plain-integer path, same as today.
- Any other `time` types (`time.Time`, etc.).

## Decisions

- [x] **D1 — duration-string parsing vs. leaving the current (accidental)
  raw-nanosecond-int64 behavior alone.** Confirms the entire point of this
  plan: route `time.Duration` fields through `flag.go`'s existing
  `time.ParseDuration`-based `Duration` constructors instead of the plain
  `Int64` path they fall into today. Rejected alternative: do nothing,
  since `time.Duration` fields already "work" today as raw-nanosecond
  int64 flags — rejected because that silently contradicts the type's own
  purpose and this package's existing flat-flag `Duration` support.
  **Settled: as proposed.**
- [x] **D2 — scope includes both value and pointer forms.** `time.Duration`
  (always populated, like any other value leaf) and `*time.Duration`
  (nil-unless-provided, like every other pointer-scalar leaf from `0002`)
  are both in scope, for consistency with every other supported kind.
  Confirm, or say if you want value-only for this increment.
  **Settled: both, as proposed.**
- [x] **D3 — `Container` gains a dedicated `durationVariables` map.**
  Mechanical consequence of D1 (a `*time.Duration` can't live in the
  existing `intVariables` map). Confirm no objection to the extra map
  field.
  **Settled: as proposed.**

All decisions settled — plan approved for implementation.
