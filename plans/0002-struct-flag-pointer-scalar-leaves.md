# 0002 — struct_flag: support pointer-to-scalar leaf fields

## Scope

`flag` submodule only (`flag/flag.go`, `flag/struct_flag.go`, their tests,
`flag/README.md`). No other module in this repo is touched.

## Problem

Plan `0001` explicitly left pointer-to-scalar leaf fields (e.g. `*string`,
`*int`, `*uint`) out of scope — `buildNodes` in `struct_flag.go` hits
`isPtr && ft.Kind() != reflect.Struct` and silently skips the field, so an
optional scalar such as

```go
type Config struct {
    Timeout *int `flag:"timeout"`
}
```

gets no flag at all today; `Timeout` is always `nil` regardless of what's
passed on the command line or in the environment.

## Goal

Support pointer-to-scalar leaf fields for the same five kind buckets
`struct_flag.go` already handles for non-pointer leaves: string, signed
integers, bool, floats, unsigned integers. Semantics mirror the
pointer-to-struct behavior already settled in `0001`'s D1: the field is
`nil` after `Parse()` unless a flag or environment variable for it was
explicitly provided (presence, not value — `-timeout=0` or `TIMEOUT=0`
still counts as provided); otherwise it's a pointer to the resolved value.

This applies at any nesting depth, inside both value-typed and
pointer-typed struct branches — a `*uint` leaf's nil-ness is independent of
whether its containing branch ends up nil.

## Design

**`flag.go`**: add five new exported constructors, one per kind bucket
`struct_flag.go` uses, each a thin wrapper around the existing constructor
with a fixed zero-value default — a custom default is meaningless for a
value whose whole point is "nil when absent":

```go
func StringP(name, usage string) *string   { return String(name, "", usage) }
func Int64P(name, usage string) *int64     { return Int64(name, 0, usage) }
func BoolP(name, usage string) *bool       { return Bool(name, false, usage) }
func Float64P(name, usage string) *float64 { return Float64(name, 0.0, usage) }
func Uint64P(name, usage string) *uint64   { return Uint64(name, 0, usage) }
```

No `*WithoutEnv` or `*Var` counterparts are added — nothing in this plan
needs them, and `struct_flag.go` is the only caller today.

**`struct_flag.go`**:

- `structNode` already has `isPtr` and `elemType` fields (used today only
  for struct/pointer-to-struct branches). Reuse both for pointer-scalar
  leaves: `isPtr = true`, `elemType` = the field's pointee type (e.g.
  `int32` for a `*int32` field), `kind` = that pointee type's `Kind()`
  (exactly as already computed for non-pointer leaves).
- In `buildNodes`, the branch that currently does
  `if isPtr { continue }` for non-struct pointer fields becomes a second
  leaf-registration switch, structurally identical to the existing one but
  calling `StringP`/`Int64P`/`BoolP`/`Float64P`/`Uint64P` instead of
  `String`/`Int64`/`Bool`/`Float64`/`Uint64`, and setting `isPtr`/`elemType`
  on the resulting node. The underlying value maps (`stringVariables`,
  `intVariables`, ...) are unchanged — a pointer-scalar leaf's registered
  flag is still stored as the widened type (`*int64`, `*uint64`,
  `*float64`) exactly like a non-pointer leaf; only the eventual struct
  field differs in whether it's a direct scalar or an allocated pointer.
- The scalar-writing switch in `populateNode` (the five `case`s keyed off
  `n.kind`) is factored into a small helper that writes into an arbitrary
  `reflect.Value`, so it can target either the real struct field
  (non-pointer leaf, as today) or a scratch value that then gets addressed
  and assigned to the field (pointer leaf):

  ```go
  if n.isPtr {
      if !provided {
          return false // field stays nil (its zero value)
      }
      target := reflect.New(n.elemType).Elem()
      writeLeafScalar(target, n, c)
      parent.Field(n.index).Set(target.Addr())
      return true
  }
  writeLeafScalar(parent.Field(n.index), n, c)
  return provided
  ```

  `provided` is computed exactly as it is today (`setFlags[n.fullName]` or
  `os.LookupEnv(envNameForFlagName(n.fullName))`) — no change to presence
  detection, just a new consumer of it at the leaf level in addition to the
  existing branch level.

## Files touched

- `flag/flag.go` — five new constructors (`StringP`, `Int64P`, `BoolP`,
  `Float64P`, `Uint64P`).
- `flag/struct_flag.go` — pointer-scalar leaf registration and population,
  as described above.
- `flag/flag_test.go` or a new small test — coverage for the five new
  constructors in isolation (presence via `flag.Visit`/env, default is
  always zero).
- `flag/struct_flag_test.go` — pointer-scalar leaf present (via flag),
  present (via env), absent (stays nil), and a pointer-scalar leaf nested
  inside both a value-typed and a pointer-typed branch.
- `flag/README.md` — document pointer-scalar leaf support alongside the
  existing pointer-struct documentation.

## Non-goals (explicitly out of scope)

- `*WithoutEnv` or `*Var` variants of the new `*P` constructors.
- Pointer-scalar support for kinds `struct_flag.go` doesn't otherwise
  support (`time.Duration`, complex, etc.) — unchanged from `0001`.
- Slices, arrays, or maps of scalars or structs — unchanged from `0001`.

## Decisions

- [x] **D1 — new public API in `flag.go` vs. purely internal reuse.**
  Two options: (a) add `StringP`/`Int64P`/`BoolP`/`Float64P`/`Uint64P` to
  `flag.go` as new public, no-default constructors that `struct_flag.go`
  calls; or (b) keep everything inside `struct_flag.go`, reusing the
  existing `String`/`Int64`/`Bool`/`Float64`/`Uint64` with a zero-value
  default passed in directly, relying on the same presence check to decide
  whether the default is ever observed.
  **Settled: option (a)**, with exactly the five names and signatures
  shown above.

All decisions settled — plan approved for implementation.
