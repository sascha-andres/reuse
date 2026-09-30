# 0005 — rename module namespace to go.livingit.de/reuse

## Problem

The module is declared as `github.com/sascha-andres/reuse`. It should move to
`go.livingit.de/reuse`.

## Scope

- `go.mod`: change `module` directive.
- All internal imports referencing `github.com/sascha-andres/reuse/...`
  (`go.mod`, `examples/time/main.go`, `examples/structflag/main.go`,
  `examples/flag/main.go`, `j/marshal.go`, `functional/filter.go`) rewritten
  to `go.livingit.de/reuse/...`.
- `README.md`: state the module path explicitly, e.g. as a `go get` line,
  so the namespace is named rather than only implied by import paths. Also
  add an explicit note that the module was renamed from
  `github.com/sascha-andres/reuse`, so a reader who has the old path
  cached knows to update it.
- No behavior change. `go vet ./...`, `go build ./...`, `go test ./...` must
  stay green.

## Decisions

None open — mechanical rename, no design choice involved.

## Verification

- `grep -rn "sascha-andres/reuse"` returns nothing.
- `go build ./...` and `go test ./...` pass.
