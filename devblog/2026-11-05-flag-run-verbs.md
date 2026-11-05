---
project: reuse
task: flag-run-verbs
topic: verb-dispatch
---

Plan: `plans/0007-flag-run-verbs.md`. Implemented by Sonnet 5, verified
by Opus 5.

`flag/runner.go` was a skeleton: the branches of `RunVerbs` and the
alphanumeric check of `isVerb` existed only as comments. Because of
this, the whole `flag` package did not build (`declared and not used:
v` three times, and a reference to `ErrInvalidVerblist` while the
declared name is `ErrInvalidVerbList`). Nothing called `RunVerbs` yet,
so no caller depended on the skeleton behaviour.

`RunVerbs(ctx, verbList, args)` dispatches on the output of
`GetVerbs()`. A value in `verbList` is a Runner or a nested
`map[string]any`; the key `""` is the default of that level.

The user changed three skeleton behaviours (D1, D3, D4 in the plan):

- An unknown verb returns `ErrUnknownVerb`, wrapped with the verb. The
  skeleton ran the default instead, so a typo (`app bulid`) ran the
  default with no error. The default still runs when `args[0]` is not
  a verb at all.
- A verb is ASCII `[A-Za-z0-9]` followed by `[A-Za-z0-9_-]`, so
  `dry-run` and `list_all` are verbs. With the strict alphanumeric
  check of the comment, `dry-run` would not be a verb and would go to
  the default.
- A `""` entry that is not a Runner returns `ErrInvalidVerbList`.
  `ErrNoDefaultRunner` now only means that `""` is missing.

D2: a func literal stored in `map[string]any` has the type
`func(context.Context) error`, not `Runner`. A type assertion to
`Runner` does not match it. `asRunner` accepts both types.

Nested levels do not fall back to the parent's default. When args run
out at a nested level, rule 1 applies at that level (`ErrNothingToRun`
if it has no `""`).

Found during implementation, D5: the first version returned
`ErrInvalidVerbList` for a nil Runner or nil func entry. Opus 5 judged
that acceptable, but the user chose a panic: a nil Runner is a bug in
the calling code and must not look like a configuration error. A
typed nil is now called and panics. An untyped `nil` matches neither
Runner type and still gives `ErrInvalidVerbList`. The panic test
checks for a nil-pointer `runtime.Error`, so an unrelated panic
before the call does not make it pass.

Verification round 1 (Opus 5): approved on condition. The tests for
"runner error returned unchanged" used `errors.Is`, which also passes
when the implementation wraps the error; they now check identity.
Nil cases were added, so removing a `r != nil` guard in `asRunner`
fails a test. Round 2: approved.
