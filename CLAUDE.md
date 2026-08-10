# reuse — CLAUDE.md

## Process rules

- **No implementation without a plan.** Every unit of work must have a
  corresponding Markdown file in `plans/` before any code is written. If a
  plan doesn't exist for the work you're about to do, write one first. Plan
  files are named `NNNN-<slug>.md`, where `NNNN` is a zero-padded four-digit
  sequence number, incremented from the highest number already present in
  `plans/` (first plan is `0001`).
- **A plan is approved, not merely reviewed, before implementation starts.**
  Review is not the gate; approval is. Until the user approves it, a plan is a
  proposal.
- **Each decision in a plan is separately approved.** A plan that still has an
  open decision is not approved, and no code is written against the parts that
  depend on it. Record decisions in a `## Decisions` section, strike each one
  as it is settled, and keep the superseded text visible rather than deleting
  it. `plans/0012-config-binding-migration.md` is the worked example: four
  decisions, each settled individually before the plan became implementable.
- **Read existing plans lazily, not exhaustively.** `plans/` is thousands of
  lines; reading all of it before writing a plan is not affordable and is not
  asked for. Instead: grep `plans/` for keywords relevant to the work, read the
  hits, then ask the user whether to read more. Do not read the directory
  speculatively, and do not skip the grep - a plan that duplicates or
  contradicts an existing one is the failure this rule exists to prevent.
- **Every plan is linked from `README.md`.** A plan a fresh agent cannot find
  is a plan that gets rewritten from scratch.
- **Fine-grained commits.** One commit per file/class touched or per
  logical unit of change. Never squash a body of work into one commit. This
  rule outranks the TODO-list grouping below: a group of TODO items maps to a
  logical unit and may well be several commits. Grouping is for planning, not
  permission to squash.
- **A TODO list per plan.** On approval, a plan gets an accompanying TODO list,
  never a single global list shared across plans. Name it after the plan and
  put it alongside: `plans/0009-foo.md` gets `plans/0009-foo.todo.md`.
  - Items are short slugs referencing the plan with a line number, e.g.
    `- [ ] trim-three-read-sites (0009:114)`. No long text; the plan holds the
    detail.
  - Group related items. A group is one logical unit of change, subject to the
    commit rule above.
  - Lifecycle: `[ ]` open, `[/]` in progress with brief context appended,
    `[x]` done with that context removed.
  - Keep `README.md` current when a TODO list changes.
- **Where each kind of record lives.** Four files, one job each, and a fact
  belongs in exactly one of them. Two copies drift, and then neither can be
  trusted.
  - `plans/NNNN-*.md` - the plan and its decisions.
  - `plans/NNNN-*.todo.md` - item state for that plan.
  - `plans/0001-progress.md` - phase state and sign-offs. State only.
  - `plans/memory.md` - durable findings and lessons that outlive the plan
    that produced them. Undated and distilled: what was measured, what it
    established or disproved, and the `file:line` or command that shows it.
    Read this instead of the whole `devblog/` directory. The devblog stays the
    dated, per-step narrative; `memory.md` is the still-true residue. An entry
    must be falsifiable, and a superseded entry is corrected in place with the
    correction visible.
- **Lint gate.** Work is not done until `go vet ./...` and `go build ./...`
  pass and `go test ./...` is green for the module you touched. Pre-existing
  findings in code you didn't touch are baseline, not yours to fix as a
  side effect.
- **Devblog.** Every relevant step (a component finished, a blocker
  resolved, a non-obvious decision made) gets a dated entry in
  `devblog/YYYY-MM-DD-<slug>.md`. Front matter tags use three namespaces:
  `project:<name>`, `task:<name>` (business task/domain), and
  `topic:<name>` (cross-cutting concern). Write these as memory for a
  future agent with no other context — include the reasoning behind a
  decision, not just what changed. This is not a changelog. Record which
  model implemented and which verified (see the next rule).
- **Sonnet 5 implements, Opus 5 verifies.** Implementation agents run on
  Sonnet 5; the agent that reviews a diff and signs a phase off runs on
  Opus 5. No agent verifies its own output, and the verifier reports rather
  than fixes — the implementer applies the fix and it goes round again.
  Escalating implementation to Opus 5 is an exception that needs a reason
  recorded in the devblog. Each devblog entry names the model that produced
  the code and the model that reviewed it.
- **README per phase.** Every phase of work leaves a `README.md` that a new
  developer *and* a fresh agent can onboard from: what this component is, how
  to build/run/test it locally, the prerequisites, the gotchas, and which of
  the AI tooling below applies to it. Update the nearest existing README
  rather than adding a parallel one. A phase is not done without it.

## Go style

- Keep code as simple as possible. Do not introduce an abstraction
  (interface, factory, base class, DI seam) unless there's a concrete
  second use case or a genuine unit-testing need for it. When unsure
  whether something needs an abstraction, stop and ask rather than
  guessing.
