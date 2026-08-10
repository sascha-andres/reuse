# 0004 — flag: add *VarP constructors with real nil-unless-provided semantics

## Scope

`flag/flag.go` and its test only. `struct_flag.go` is untouched — the
design below is additive and doesn't change anything `struct_flag.go`
already depends on.

## Problem / motivation

`0002`/`0003` added no-default constructors (`StringP`, `Int64P`, `BoolP`,
`Float64P`, `Uint64P`, `DurationP`) that return a freshly allocated,
always-non-nil `*T` holding the zero value until parsed. That's the right
shape for `struct_flag.go` (which does its own separate nil/non-nil
decision via reflection after `Parse()`), but it's not directly useful to
a caller who wants an ordinary `*string`/`*int`/etc. variable that ends up
`nil` unless the flag was actually provided — there's no way to know
"provided vs. defaulted" until after `Parse()` has run, and the existing
`*P` functions return their pointer before that point.

This also picks up the abandoned `examples/flag/main.go` edit
(`flag.BoolVarP(boolFlagP, "another boolean flag")`, reverted last
session) — a caller trying to get exactly this "give me nil unless it was
set" behavior with no built-in support for it.

## Goal

Six new functions — `StringVarP`, `Int64VarP`, `BoolVarP`, `Float64VarP`,
`Uint64VarP`, `DurationVarP` — plus a new `ResolveP()`. Each `*VarP`
function takes the *address of a pointer variable* (`**T`) instead of a
`*T`; after the package-level `Parse()` and a call to `ResolveP()`, that
pointer variable is `nil` if the flag was never explicitly provided (CLI
or environment variable) and non-nil, pointing at the resolved value,
otherwise:

```go
var timeout *int
flag.IntVarP(&timeout, "timeout", "timeout in seconds")
flag.Parse()
flag.ResolveP()
// timeout == nil unless -timeout or TIMEOUT was set
```

## Design

**Additive, not breaking.** `StringP`/`Int64P`/`BoolP`/`Float64P`/
`Uint64P`/`DurationP` keep their existing signature and behavior exactly
as shipped in `0002`/`0003` — `struct_flag.go`'s `buildNodes` keeps
calling them unmodified. Each `*VarP` function is built *on top of* its
already-existing `*P` counterpart:

```go
type pResolution struct {
    name    string
    resolve func(provided bool)
}

var pendingResolutions []pResolution

func StringVarP(p **string, name, usage string) {
    inner := StringP(name, usage)
    pendingResolutions = append(pendingResolutions, pResolution{
        name: name,
        resolve: func(provided bool) {
            if !provided {
                *p = nil
                return
            }
            v := *inner
            *p = &v
        },
    })
}
```

`Int64VarP`, `BoolVarP`, `Float64VarP`, `Uint64VarP`, `DurationVarP`
follow the identical pattern, each wrapping its own `*P` counterpart.

`ResolveP()` walks `pendingResolutions` once, determining "provided" per
entry the same way `struct_flag.go`'s `populateNode` already does for
pointer leaves (`Visit` for explicit CLI flags, `os.LookupEnv` on
`envNameForFlagName(name)` for the environment fallback) — duplicated
here rather than shared with `struct_flag.go`, to keep this plan fully
additive and avoid touching already-shipped code for an internal DRY
cleanup that isn't required for correctness:

```go
func ResolveP() {
    set := make(map[string]bool)
    Visit(func(fl *f.Flag) { set[fl.Name] = true })
    for _, r := range pendingResolutions {
        provided := set[r.name]
        if !provided {
            _, provided = os.LookupEnv(envNameForFlagName(r.name))
        }
        r.resolve(provided)
    }
    pendingResolutions = nil
}
```

Must be called after the package-level `Parse()`, same ordering
requirement `struct_flag.Container.Parse()` already documents. Clearing
`pendingResolutions` after resolving makes it safe to call once per
`Parse()` cycle without manual cleanup between test runs, as long as each
test that registers a `*VarP` also calls `ResolveP()` — which is how the
new tests are written, so no changes are needed to the existing
`resetForStructFlagTest` test helper.

## Files touched

- `flag/flag.go` — `pResolution`, `pendingResolutions`, the six `*VarP`
  functions, `ResolveP`.
- `flag/flag_test.go` — all six absent-stays-nil, all six present-via-flag,
  all six present-via-env.
- `flag/README.md` — new section documenting `*VarP` + `ResolveP`.

## Non-goals (explicitly out of scope)

- `IntVarP`/`UintVarP` (plain, non-64-bit) — same reasoning as `0004`'s
  original draft: no existing `*P` precedent to build on (`Int`/`Uint`
  never got one).
- Changing `StringP`/`Int64P`/`BoolP`/`Float64P`/`Uint64P`/`DurationP`'s
  signature or behavior.
- Sharing the presence-check logic with `struct_flag.go` via a common
  helper — duplicated intentionally to keep this plan additive.
- Reinstating or fixing the abandoned `examples/flag/main.go` edit.

## Decisions

- [x] **D1 — signature and mechanism.** `func XVarP(p **T, name, usage
  string)`, built as a wrapper around the existing, unchanged `XP`
  constructor, resolved by a new package-level `ResolveP()` called after
  `Parse()`. `XP`'s own signature/behavior is untouched, so
  `struct_flag.go` needs no changes.
  **Settled: as described above** (the additive/wrapper option, not the
  alternative of changing `XP`'s return type).

All decisions settled — plan approved for implementation.
