---
project: reuse
task: flag-verb-flag-interleaving
topic: flag-parsing
---

`flag.Parse()` silently dropped any flag that appeared after the first
verb, and everything after it (further flags, further verbs). Reported
by the user as: `app -flag 1 verb` works, `app verb -flag 1` doesn't, and
"also flags between verbs" (`app verb1 -flag 1 verb2`) needs to work too.

Root cause: in the loop over `os.Args[1:]`, the branch for a `-`-prefixed
token started with `if len(verbs) > 0 { break }`. `break` exits the
whole `for` loop, not just that branch — the intent look liked "stop
collecting new verbs once flags resume", but the effect was "stop parsing
anything at all the moment a flag follows a verb". The token that
triggered it, and every token after, never reached
`f.CommandLine.Parse(arguments)`, so a flag like `-flag` silently kept
its zero/default value with no error.

The `previousIsFlag` state machine in the rest of the loop already
generalizes correctly to flags interleaved with verbs — it was only ever
gated off by this one `break`. Fix was deletion of the guard, no other
structural change (`flag/flag.go:479` before the fix, in `Parse()`).

Test isolation gap found along the way: `resetForStructFlagTest` (in
`struct_flag_test.go`) saved/restored `booleanFlags`, `envPrefix`,
`overriddenEnvPrefixes`, `os.Args`, `f.CommandLine`, but not the
package-level `verbs` slice. Any test that calls `Parse()` and checks
`GetVerbs()` would have accumulated verbs left over from earlier tests
in the same run. Added `verbs` to the same save/reset/restore. This
means `TestGetVerbs` (already `t.Skip`ped, unrelated issue: it mutates
`os.Args` directly rather than going through the reset helper) is still
the odd one out — new verb/flag-interleaving tests use
`resetForStructFlagTest` instead.

Added coverage: flag after a verb, flag between two verbs, boolean flag
between two verbs (must not consume the second verb as its value), and a
string flag whose value is a single argument containing an embedded
space (`"test test"`) — the last one to confirm parsing works on the
*token* level and doesn't do any accidental splitting that verb/flag
interleaving might tempt.

Opus 5 verification of the above (implemented by Sonnet 5) found the
same silent-drop failure through three more syntaxes in the same branch,
one of which (`-test=1 verb -other 2`) reproduces the originally
reported symptom exactly, just via `=` instead of a space:

- `-flag=value` and `-bool=false` — the inline `=` value wasn't
  recognized, so `previousIsFlag` got set from the bare flag-name check
  and stole the next token (a verb, or another flag) into `arguments`.
- `--flag` (double dash) — `strings.TrimPrefix(v, "-")` only stripped one
  `-`, so double-dash booleans never matched `booleanFlags`.

Fixed by switching to `strings.TrimLeft` (strips all leading `-`) and
splitting on the first `=` to detect an inline value, short-circuiting
`previousIsFlag` when one is present regardless of whether the flag is
boolean.

Explicitly left unfixed (follow-up, not part of this plan): a flag value
that itself looks like a flag, e.g. `-num -5` between verbs — the
`previousIsFlag` state machine mis-classifies `-5` as a new flag rather
than `-num`'s value, and fixing it needs per-flag arity/type tracking
this package doesn't currently keep outside stdlib's own `FlagSet`. Also
left unfixed: `Parse()` never resets the package-level `verbs` slice
across repeated calls in the same process (verbs accumulate) — the
test-helper reset added alongside the original fix only covers test
isolation, not this.

Implemented by Claude (Sonnet 5); verified by a separate Claude (Opus 5)
agent reviewing the diff, per this repo's implement/verify split. The
`=`/`--` amendment was applied by the same Sonnet 5 implementer after
the verifier's report, with the user choosing the amendment's scope
(fix `=`/`--`/`-bool=false`, defer the negative-number-value case).

See `plans/0006-flag-parse-flags-after-verbs.md`.
