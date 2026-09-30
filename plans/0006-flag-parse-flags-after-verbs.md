# 0006 — flag: stop dropping flags that appear after or between verbs

## Scope

`flag/flag.go`'s `Parse()` and its test only.

## Problem / motivation

`Parse()` splits `os.Args` into verbs (bare tokens) and flag arguments
(handed to `f.CommandLine.Parse`). Once at least one verb had been seen,
any subsequent token starting with `-` hit:

```go
if strings.HasPrefix(v, "-") {
    if len(verbs) > 0 {
        break
    }
    ...
```

`break` exits the whole loop, not just the flag/verb branch — so every
argument after the first flag-following-a-verb is silently discarded,
including that flag itself and anything after it (more flags, more
verbs). Concretely:

- `app -flag 1 verb` worked (no verb seen yet when `-flag` is hit).
- `app verb -flag 1` dropped `-flag 1` entirely: `-flag` never reaches
  `f.CommandLine.Parse`, so it silently keeps its zero/default value.
- `app verb1 -flag 1 verb2` dropped `-flag 1 verb2` outright, losing a
  second verb too.

Boolean flags aren't special-cased in the bug — the `break` fires before
the existing `booleanFlags` check ever runs — but they're worth calling
out because they're the case most likely to be used as a valueless
sub-flag between verbs (e.g. `app verb1 -v verb2`), which must not
consume `verb2` as the flag's value.

## Root cause

The `len(verbs) > 0` guard was presumably meant to stop verb-collection
once flags start, on the assumption verbs are a contiguous prefix. That
assumption isn't what the rest of the function implements (the
`previousIsFlag` state machine already tolerates flags anywhere in the
stream) and isn't what callers need: verbs and flags are interleaved on
a real CLI (`app verb1 -flag 1 verb2`).

## Fix

Delete the guard. The surrounding state machine already does the right
thing without it:

```go
if strings.HasPrefix(v, "-") {
    if !slices.Contains(booleanFlags, strings.TrimPrefix(v, "-")) {
        previousIsFlag = true
    }
    arguments = append(arguments, v)
} else {
    if !previousIsFlag {
        verbs = append(verbs, v)
    } else {
        arguments = append(arguments, v)
    }
    previousIsFlag = false
}
```

- A `-`-prefixed token is always appended to `arguments`; `previousIsFlag`
  is set unless the flag is boolean, exactly like the existing
  before-first-verb case already does.
- A bare token becomes a verb unless it's the value of the immediately
  preceding non-boolean flag.

This makes flags after a verb, and further verbs after that, parse
normally: `app verb1 -flag 1 verb2` yields `verbs=["verb1","verb2"]`,
`arguments=["-flag","1"]`.

Also update `Parse`'s doc comment, which described the old (and buggy)
"verbs before the first flag" restriction, and drop a stale commented-out
`re-add the first argument` block that referenced removed code.

## Files touched

- `flag/flag.go` — remove the `break`, update doc comment, drop dead
  comment block.
- `flag/flag_test.go` — new case(s) covering a flag after a verb and a
  flag between two verbs (including a boolean flag between verbs, to
  confirm it doesn't consume the next verb as a value). The existing
  `TestGetVerbs` table is `t.Skip`ped (documented as broken because it
  mutates the shared `os.Args`/`CommandLine` globals) — new coverage goes
  in a `Parse`-driven test using the same `resetForStructFlagTest`
  pattern as the `*VarP` tests, not that table.

## Amendment: `-flag=value` / `--flag` / `-flag=false` also silently drop

Opus 5 verification of the D1 fix (below) found the same silent-drop
failure through three more syntaxes, all in the same branch:

- `strings.TrimPrefix(v, "-")` only strips one leading `-`, so `--on`
  (double-dash boolean) never matches `booleanFlags`; `previousIsFlag`
  gets set and the next verb is swallowed.
- The inline-value form `-flag=value` (and `-bool=false`) isn't
  recognized as already having its value attached, so `previousIsFlag`
  gets set from the flag-name-check alone and the next token — a verb or
  another flag — is stolen into `arguments` and everything after it can
  cascade into the same drop.

`cmd -test=1 verb -other 2` reproduced the exact originally-reported
symptom (`-other`'s value silently lost, no error) through the `=` form
instead of the space form.

### Fix

```go
if strings.HasPrefix(v, "-") {
    name := strings.TrimLeft(v, "-")
    hasInlineValue := false
    if idx := strings.Index(name, "="); idx >= 0 {
        name, hasInlineValue = name[:idx], true
    }
    if !hasInlineValue && !slices.Contains(booleanFlags, name) {
        previousIsFlag = true
    }
    arguments = append(arguments, v)
}
```

`TrimLeft` (not `TrimPrefix`) strips all leading `-`, covering `--flag`.
Splitting on the first `=` and short-circuiting `previousIsFlag` when one
is present covers `-flag=value` and `-bool=false` — the token already
carries its own value, so the next token must not be treated as that
flag's value regardless of whether the flag is boolean.

Not fixed here (explicit non-goal, see below): a flag value that itself
looks like a flag, e.g. `-num -5` between verbs — stdlib's own parser
handles this token opaquely once it reaches `f.CommandLine.Parse`, and
this package's own `previousIsFlag` state machine currently
mis-classifies `-5` as introducing a new flag rather than being the
value of `-num`, which loses verbs/flags after it. Fixing this requires
knowing each flag's arity/type at this point in `Parse`, which the
package doesn't track today (`booleanFlags` is the only per-flag
metadata kept outside stdlib's own `FlagSet`). Left as a follow-up.

## Non-goals

- Fixing or re-enabling `TestGetVerbs` — pre-existing, unrelated to this
  bug.
- Changing `GetSeparated`/`--` handling.
- Fixing `-flagvalue -5`-shaped ambiguity (flag value that looks like a
  new flag) — see amendment above.
- Fixing `Parse()` never resetting the package-level `verbs` slice across
  repeated calls (found during verification, `flag/flag.go:463`) — a
  real but separate pre-existing bug, not part of the reported symptom.

## Decisions

- [x] **D1 — fix shape.** Delete the `len(verbs) > 0` guard rather than
  reworking the state machine; the existing `previousIsFlag` logic
  already generalizes correctly once the early `break` is gone.
  **Settled: as described above.**
- [x] **D2 — amendment scope.** Fix `-flag=value`, `--flag`, and
  `-flag=false` (all found by verification, all same root cause/branch)
  in this same plan; leave `-flag -5`-shaped ambiguity and the
  `verbs`-not-reset-across-`Parse()`-calls bug for a follow-up.
  **Settled: user chose "fix #1-#3" when asked.**

All decisions settled — plan approved for implementation.
