# 0001 — struct_flag: support nested structs

## Scope

`flag` submodule only (`flag/struct_flag.go` and its test/README). No other
module in this repo is touched.

## Problem

`flag/struct_flag.go` builds CLI flags (and, via the existing env-var
machinery in `flag/flag.go`, environment variable fallbacks) from a struct's
`flag:"..."` tags. Two things are missing/broken today:

1. **Nested structs are not implemented.** `AddFlagsForStruct` hits
   `field.Type.Kind() == reflect.Struct` and just does
   `fmt.Println("struct implementation")` — no flags are created, the field
   is silently dropped, `Parse()` never touches it.
2. **Latent key collision.** The five value maps on `Container`
   (`stringVariables`, `intVariables`, ...) are keyed by the *leaf* tag
   (e.g. `"ip"`), not the full flag name. Two sibling structs that both have
   a field tagged `ip` (exactly the `BindHttp`/`BindGrpc` example in the
   request) already collide today for any field added inside a struct-typed
   branch, and will collide immediately once nesting is implemented naively.

## Goal

Support both of the shapes in the request:

```go
type Config struct {
    Name     string `flag:"name"`
    BindHttp struct {
        IP   string `flag:"ip"`
        Port uint   `flag:"port"`
    } `flag:"http"`
    BindGrpc struct {
        IP   string `flag:"ip"`
        Port uint   `flag:"port"`
    } `flag:"grpc"`
}
```

and

```go
type Binding struct {
    IP   string `flag:"ip"`
    Port uint   `flag:"port"`
}

type Config struct {
    Name     string   `flag:"name"`
    BindHttp Binding  `flag:"http"`
    BindGrpc *Binding `flag:"grpc"`
}
```

Both produce flags `-name -http-ip -http-port -grpc-ip -grpc-port` (given an
empty top-level prefix, matching the existing `fmt.Sprintf("%v-%v", prefix,
tag)` naming already used for flat fields — unchanged).

Nesting depth is arbitrary (struct containing a struct containing a struct),
not just one level.

Env var names need no new logic: `flag/flag.go`'s `String`/`Int64`/`Bool`/
`Float64`/`Uint64` already resolve environment variables through
`envNameForFlagName` (env prefix + `-`→`_` + uppercase) whenever they're
called with a name. Since `struct_flag.go` is in the same package, calling
those constructors with the fully-qualified nested name (`"http-ip"`, etc.)
gets env support for free — no change needed to `flag.go`.

`*Binding` (pointer to a named or anonymous struct) is `nil` after `Parse()`
unless at least one field inside it was actually supplied, via CLI flag or
env var. Value-typed nested structs (`Binding`, not `*Binding`) are always
populated (fields default to their normal zero-value-or-env-or-flag
resolution, same as today for flat fields).

## Design

Replace the flat `map[tag]*T` value maps (keyed by leaf tag — the source of
the collision bug) with a small node tree built once in `AddFlagsForStruct`
and walked again in `Parse()`:

```go
type node struct {
    field    reflect.StructField
    index    int          // field index within its immediate parent struct
    isPtr    bool         // field type is a pointer to a struct
    children []*node      // non-nil => this node is a struct/*struct branch
    kind     reflect.Kind // meaningful only for leaves (children == nil)
    fullName string       // e.g. "http-ip" (leaf) or "http" (branch)
}
```

`Container[T]` keeps the five typed value maps as before, but keyed by
`fullName` (globally unique per registered flag) instead of the leaf tag —
that alone fixes the collision. `Container` additionally holds the
top-level `[]*node` produced by registration.

**Registration** (`AddFlagsForStruct`): recursively walk struct fields. A
field with kind `Struct`, or kind `Pointer` to `Struct`, becomes a branch
node: recurse into its (dereferenced) type with the accumulated name as the
new prefix, collecting child nodes. Anything else becomes a leaf: call the
existing `String`/`Int64`/`Bool`/`Float64`/`Uint64` constructor with the
full name and store the returned pointer in the appropriate map under
`fullName`, exactly as today but keyed correctly.

**Parse()**: walk the node tree in lockstep with `reflect.Value`. For a
leaf, set the field from the map entry (as today, but looked up by
`fullName`) and report whether it was "provided" (see below). For a branch:
- non-pointer: recurse directly into `parent.Field(index)`, populate its
  leaves in place, ignore the "provided" result (value structs are always
  written).
- pointer: recurse into a scratch `reflect.Value` of the element type;
  if any descendant leaf was "provided", allocate (`reflect.New`), copy the
  scratch value in, and set the field; otherwise leave the field at its
  zero value (`nil`). Propagate "any descendant provided" up to the parent
  branch so pointer-to-struct-containing-pointer-to-struct nests correctly.

"Provided" for a leaf = `true` if the flag was explicitly set on the command
line (`Visit`, which — per `flag/flag.go` — only calls back for flags set
during `Parse()`) **or** its resolved env var name is present in the
environment (`os.LookupEnv` on the same name `envNameForFlagName` would
compute — reusable directly since it's an unexported function in the same
package).

## Files touched

- `flag/struct_flag.go` — rework described above.
- `flag/struct_flag_test.go` — new; there is no existing test file for this
  code. Covers: flat fields (regression), anonymous nested struct, named
  nested struct (value), named nested struct (pointer, present), named
  nested struct (pointer, absent → nil), multi-level nesting, env-var-only
  provision of a pointer branch, and the sibling-collision case
  (`BindHttp.IP` / `BindGrpc.IP`) that is broken today.
- `flag/README.md` — document nested-struct support with the two examples
  from the request.

`examples/flag/main.go` is not touched — it demonstrates the flat
`BoolVar`/`StringVar` API, unrelated to `struct_flag.go`.

## Non-goals (explicitly out of scope)

- Slices, arrays, or maps of structs.
- Pointers to scalar leaf fields (e.g. `*string` leaf) — only pointer-to-
  struct is handled, per the request.
- `time.Duration` fields, at any nesting level — not supported by
  `struct_flag.go` today either; unrelated to this rework.
- Anything about unexported fields beyond "don't panic" (see Decision 4).

## Decisions

- [x] **D1 — "provided" semantics for nilable pointer branches.** A leaf
  counts as provided if explicitly set via CLI (`flag.Visit`) or if its env
  var is present (`os.LookupEnv`), regardless of value — e.g.
  `-grpc-active=false` or `GRPC_ACTIVE=false` both count as "provided" (the
  pointer gets allocated), matching "nil if no flag (or env variable) are
  present" read literally as presence, not truthiness. Confirm this is the
  intended reading.
  **Settled: presence-based, as proposed.** Explicitly configuring a leaf
  (CLI flag present per `flag.Visit`, or its env var present per
  `os.LookupEnv`) makes the branch non-nil even if the resulting value
  equals that field's zero value (e.g. `-grpc-port=0` or `GRPC_PORT=0`
  alone makes `BindGrpc` non-nil with `Port=0`). Rejected the alternative
  (value-based: non-nil only if some leaf's value differs from zero) as
  unable to distinguish "explicitly set to zero" from "unset."
- [x] **D2 — breaking change to `Container`'s private layout.** The five
  value maps move from keyed-by-leaf-tag to keyed-by-full-name, and
  `Container` gains a node tree. `AddFlagsForStruct`'s and `Parse()`'s
  public signatures are unchanged. Grepped the repo: nothing outside
  `flag/struct_flag.go` and its (nonexistent) tests references `Container`'s
  private fields, so this is safe. Confirm proceeding without a
  deprecation/compat shim (project convention is to just change the code,
  no back-compat hacks).
  **Settled: proceed as proposed**, keyed-by-full-flag-name plus a node
  tree, no compat shim.
- [x] **D3 — unexported struct fields.** Today, an unexported field with a
  `flag` tag would panic on `reflect.Value.SetString`/etc. (can't `Set` an
  unexported field). Proposed: skip unexported fields silently during
  registration (`field.PkgPath != ""` / `!field.IsExported()`), same
  treatment as an untagged field, rather than returning an error. Confirm
  silent-skip is preferred over `AddFlagsForStruct` returning an error.
  **Settled: silent skip, as proposed.**

All decisions settled — plan approved for implementation.
