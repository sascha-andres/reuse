---
project: reuse
task: runner-example
topic: verb-dispatch
---

Plan: `plans/0008-runner-example.md`. Implemented by Sonnet 5, verified
by Opus 5 (approved in one round, nits only).

`examples/runner/main.go` shows `flag.RunVerbs` (plan 0007) with the
tree `greet hello|bye` and `version`.

Leaf flags: the `flag` package has one global flag set and one global
`Parse()`, so a verb cannot own its flags. `main` defines all flags
(`-name`, `-loud`) before `flag.Parse()`, and each leaf reads the
variables it needs. Because of plan 0006, flags can be before, between
or after verbs (`runner greet -name Ada hello -loud` works). Per-verb
flag sets would need a change to `flag/` and were a non-goal.

`flag.SetEnvPrefix("RUNNER")` is set because the package reads an
environment variable for each flag; without a prefix, `NAME` would set
`-name`. The verifier confirmed that `NAME=Bob` does not leak.

Each level with child maps has a `""` entry that prints the keys of
that level (sorted, map children with ` ...`). The entry is a closure
over its own map, so the map is built first and `""` is added after.

Gotcha for manual runs in zsh: `for a in "greet hello -name Ada"; do
runner $a` passes one argument, because zsh does not word-split `$a`.
Both agents hit this; use `${=a}`.

Amendment (plan 0008 D4, D5, user): the setup moved from `main` into
`run(ctx) error`; `main` only calls `run`, prints the error and exits
with code 1. `run` sets the env prefix, defines the flags, parses,
builds the verb list and calls `RunVerbs`. Flags stay global (D2). The
user's first wording was "parse the flags inside the runner"; a
per-leaf `flag.Parse()` was rejected because flags would then only work
after the leaf verb, and the root default could not take flags. A
`flag/README.md` line points to the example. Sonnet 5 implemented,
Opus 5 verified the same 13 manual cases as before: no behaviour change.

Behaviour worth knowing: extra args after a leaf are ignored
(`runner greet hello extra` runs `hello`), as documented in plan 0007.
An undefined flag exits with code 2 from the `flag` package's own
`Parse`, before `RunVerbs` runs.
