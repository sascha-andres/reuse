# 0008 — examples: RunVerbs example

## Scope

New `examples/runner/main.go`. No change to the `flag` package.
Amendment (D4, D5): one pointer line in `flag/README.md`.

## Motivation

`RunVerbs` (plan 0007) has no runnable example. The other `examples/`
directories show one feature each with a single `main.go`.

## Design

A small verb tree:

```
""            prints the valid verbs of the root level
greet
  ""          prints the valid verbs of the greet level
  hello       leaf, flags -name (string), -loud (bool)
  bye         leaf, flag -name (string)
version       leaf, no flags
```

- **Flags of the leaves.** The `flag` package has one global flag set
  and one global `Parse()`. A leaf cannot own a flag set. So `main`
  defines all flags with `flag.StringVar`/`flag.BoolVar` before
  `flag.Parse()`, and the leaf runner reads the variables. Flags can be
  before, between or after verbs (plan 0006), for example
  `runner greet -name Ada hello -loud`. `flag.SetEnvPrefix("RUNNER")`
  so that a `NAME` environment variable does not set `-name`.
- **Default runners.** Each level with child nodes gets a `""` entry
  that prints the keys of that level, sorted, one per line, without
  `""`. A child that is a map is printed with a trailing ` ...`
  (for example `greet ...`). The `""` entry is a closure over its
  map, so the map is built first and the `""` entry is added after.
- **Errors.** `main` prints the error of `RunVerbs` to stderr and
  exits with code 1. For `ErrUnknownVerb` the message names the verb.

## Verification

`go vet ./...`, `go build ./...`. Run by hand (no test, like the other
examples):

- `runner` and `runner greet`: print the verbs of that level.
- `runner greet hello -name Ada -loud`, `runner greet -name Ada bye`,
  `runner version`: run the leaf, flags applied.
- `runner greet nope`: `unknown verb: "nope"`, exit code 1.

## Non-goals

- Per-verb flag sets or per-verb usage text.
- Changes to `flag/`.

## Decisions

- [x] **D1 — verb tree.** As in Design (`greet hello|bye`, `version`).
  **Proposed: as shown.** **Settled: as shown (user).**
- [x] **D2 — flags.** Global flags defined in `main`, read by the
  leaves. This is the only option the package supports today; a
  per-verb flag set would need a change to `flag/` (non-goal).
  **Proposed: global flags.** **Settled: global flags (user).**
- [x] **D3 — default output.** Sorted keys, one per line, map children
  with a trailing ` ...`.
  **Proposed: as shown.** **Settled: as shown (user).**

~~All decisions settled — plan approved for implementation.~~

Amendment after the first implementation (user):

- [x] **D4 — README pointer.** Add one line to the "Run verbs" section
  of `flag/README.md` that points to `examples/runner`.
  **Settled: add it (user).**
- [x] **D5 — `main` → `run`.** `main` only creates the context, calls
  `run(ctx)`, prints the error to stderr and exits with code 1. `run`
  sets the env prefix, builds the verb list, defines the flags, calls
  `flag.Parse()` and `RunVerbs`. D2 stays: flags are global and can be
  before, between or after verbs.
  Rejected alternative: each leaf defines its flags and calls
  `flag.Parse()` itself; flags would only work after the leaf verb.
  **Settled: `run` function (user).**

All decisions settled — amendment approved for implementation.
