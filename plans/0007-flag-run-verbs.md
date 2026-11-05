# 0007 — flag: implement RunVerbs

## Scope

`flag/runner.go`, a new `flag/runner_test.go`, and `flag/README.md`.

## Problem / motivation

`flag/runner.go` has a skeleton of `RunVerbs` and `isVerb`. The missing
parts are described only in comments. Because of this, the `flag`
package does not build today (`go build ./...`):

- `runner.go:20`, `:27`, `:33`: `declared and not used: v`
- `runner.go:30`: `undefined: ErrInvalidVerblist` (the declared name is
  `ErrInvalidVerbList`)

Nothing calls `RunVerbs` or `isVerb` yet.

## Intended behaviour (from the comments)

`RunVerbs(ctx, verbList, args)` dispatches on `args`, which is the verb
list from `GetVerbs()`. `verbList` maps a verb to a `Runner` or to a
nested `map[string]any`. The key `""` is the default of that level.

Rules as amended by D1 to D4 ("Runner" means a `Runner` or a plain
`func(context.Context) error`, D2):

1. `args` empty:
   - `""` is a Runner: run it.
   - `""` exists, not a Runner: `ErrInvalidVerbList` (D4).
   - no `""`: `ErrNothingToRun`.
2. `args[0]` is a verb:
   - key in `verbList`, Runner: run it.
   - key in `verbList`, `map[string]any`: `RunVerbs(ctx, sub, args[1:])`.
   - key in `verbList`, other type: `ErrInvalidVerbList`.
   - not a key: `ErrUnknownVerb`, wrapped with the verb
     (`fmt.Errorf("%w: %q", ErrUnknownVerb, args[0])`) (D1).
3. `args[0]` is not a verb: same as rule 1 for the `""` entry, but a
   missing `""` gives `ErrNoDefaultRunner`.

`isVerb` (D3): true when the value is not empty, the first character is
ASCII `[A-Za-z0-9]`, and all other characters are ASCII `[A-Za-z0-9_-]`.

Original skeleton rules, superseded by D1 and D4:
~~1. `""` exists, not a Runner: `ErrNoDefaultRunner`.~~
~~3. Otherwise (not a verb, or unknown verb): run the `""` Runner of this
level. If there is none: `ErrNoDefaultRunner`.~~
~~`isVerb` returns false for any value that is not only alphanumeric.~~

## Fix

- Fill in the three commented branches as described above, with a type
  switch for the Runner/map check (see D2).
- Add the exported sentinel `ErrUnknownVerb` (D1).
- Use `ErrInvalidVerbList` (fix the misspelled reference).
- Implement the alphanumeric check in `isVerb`.
- Extend the `RunVerbs` doc comment: `args` is normally `GetVerbs()`,
  and the rules 1 to 3 above.

## Files touched

- `flag/runner.go` — the fix.
- `flag/runner_test.go` — new table test: each rule 1 to 3, nested map,
  ~~unknown verb falls to default~~ unknown verb gives `ErrUnknownVerb`
  (D1), invalid entry type, runner error is
  returned unchanged, `isVerb` cases.
- `flag/README.md` — short `RunVerbs` section with an example.

## Non-goals

- Changes to `Parse()` or `GetVerbs()`.
- Passing remaining args to a Runner (the `Runner` signature has no
  args; changing it is a separate API decision).
- Renaming exported errors.

## Decisions

- [x] **D1 — unknown verb.** ~~Rule 3 as the skeleton has it: an unknown
  verb runs the default Runner of that level (at a nested level, the
  default of that level).~~
  Alternative: return a new `ErrUnknownVerb` when `args[0]` is a verb
  but not a key. Tradeoff: the skeleton behaviour lets the default
  treat the verb as data; but a typo (`app bulid`) runs the default
  silently.
  ~~**Proposed: keep the skeleton behaviour.**~~
  **Settled: `ErrUnknownVerb` (user).**
- [x] **D2 — accepted Runner types.** A plain func literal stored in
  `map[string]any` has the type `func(context.Context) error`, not
  `Runner`. A type assertion to `Runner` does not match it, so
  `map[string]any{"x": func(ctx context.Context) error {...}}` would
  give `ErrInvalidVerbList`.
  Options: a) accept both `Runner` and `func(context.Context) error`;
  b) accept only `Runner` and document that callers must convert.
  **Proposed: a.** **Settled: a, accept both (user).**
- [x] **D3 — "alphanumeric".**
  a) ASCII only, `[A-Za-z0-9]`;
  b) Unicode `unicode.IsLetter || unicode.IsDigit`;
  c) a or b, plus `-` and `_` after the first character (allows
  `dry-run`, `list_all`).
  Tradeoff: a and b follow the comment exactly but make `dry-run` not
  a verb, so it silently goes to the default (D1). c is more useful for
  CLIs but goes past the comment.
  ~~**Proposed: a.**~~
  **Settled: c with ASCII base: first character `[A-Za-z0-9]`, others
  `[A-Za-z0-9_-]` (user).**
- [x] **D4 — default entry that is a map.** ~~`""` mapped to a
  `map[string]any` returns `ErrNoDefaultRunner` (skeleton behaviour)
  rather than `ErrInvalidVerbList`.~~
  ~~**Proposed: keep the skeleton behaviour.**~~
  **Settled: alternative (user). A `""` entry that is not a Runner
  returns `ErrInvalidVerbList`. `ErrNoDefaultRunner` only means `""` is
  missing.**
- [x] **D5 — nil Runner (found during implementation).** A nil
  `Runner` or nil `func(context.Context) error` entry would panic when
  called. The implementation treats it as "not a Runner" and returns
  `ErrInvalidVerbList`. The Opus 5 verifier judged this acceptable.
  Alternative: let it panic (caller bug).
  ~~**Proposed: `ErrInvalidVerbList` (as implemented).**~~
  **Settled: panic (user). A typed nil `Runner` or nil
  `func(context.Context) error` is called and panics. An untyped `nil`
  value matches neither type and still gives `ErrInvalidVerbList`.**

~~All decisions settled — plan approved for implementation.~~
~~D1 to D4 settled and implemented. D5 is open.~~
~~All decisions settled. D5 needs a change: remove the nil guards in
`asRunner`, nil-Runner tests expect a panic.~~
All decisions settled and implemented. Verified by Opus 5.
