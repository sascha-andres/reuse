# flag

is an opinionated replacement for the flag package. It aims to be a drop in replacement for 80% of the use cases
while providing an easy way for falling back to environment variables if present

## Exclude flags from env

This adds a bunch of methods *WithoutEnv which does not query environment variables.   

## Verbs

This flag package provides a list of verbs. That is something passed to the command line without a previous flag. Use `GetVerbs()` to retrieve.

    cmd verb -bool verb2 -comment text

`GetVerbs()` will return `[]string{"verb", "verb2"}`

The boolFlag will be set to true and the commentFlag will be set to "text".

Flags are recognized in any of the usual forms — `-flag value`,
`-flag=value`, `--flag` — anywhere in the argument list, before, between,
or after verbs.

## Run verbs

`RunVerbs` dispatches on the verbs. It takes a nested `map[string]any` and the
result of `GetVerbs()`. A value is a `Runner`, a plain
`func(context.Context) error` or a nested `map[string]any`. The key `""` is
the default of a level.

```go
flag.Parse()
err := flag.RunVerbs(ctx, map[string]any{
    "":      runServer,
    "serve": runServer,
    "db": map[string]any{
        "up":   dbUp,
        "down": dbDown,
    },
}, flag.GetVerbs())
```

- No verb (or the first argument is not a verb): the `""` entry runs. If it
  is missing, the result is `ErrNothingToRun` (no arguments) or
  `ErrNoDefaultRunner`.
- A verb that is a key: a `Runner` runs, a nested map is used for the next
  verb. At the nested level, no verb runs the `""` entry of that level.
- A verb that is not a key: an `*UnknownVerbError` that wraps
  `ErrUnknownVerb`. Check with `errors.Is` or `errors.As`. `Path` holds the
  verbs of the levels above, `Suggestions` the verb keys of the level within
  the optimal string alignment distance (Levenshtein plus swap of two
  adjacent characters), nearest first.
- An entry of another type (also a `""` entry that is not a `Runner`):
  `ErrInvalidVerbList`.
- The error of a `Runner` is returned unchanged.
- A typed nil `Runner` or nil `func(context.Context) error` entry panics when
  it runs. An untyped nil entry gives `ErrInvalidVerbList`.

Options change the suggestions:

- `WithMaxDistance(n)`: max distance, default 1. A negative value disables
  suggestions.
- `WithIgnoreCase()`: compare without case. The lookup stays case-sensitive.

```go
var uv *flag.UnknownVerbError
if errors.As(err, &uv) && len(uv.Suggestions) > 0 {
    fmt.Printf("did you mean: %s\n", strings.Join(uv.Suggestions, ", "))
}
```

A runnable example is in `examples/runner`.

A verb starts with `[A-Za-z0-9]` and continues with `[A-Za-z0-9_-]`, for
example `build`, `dry-run`, `list_all`.

## Separated

If you want to pass arguments for something like a sub command you can use the separate feature. Activate it using `flag.SetSeparated()`. Everything after `--` will be treated as separate from the command line and not parsed as verbs or flags.

    cmd verb -bool verb2 -comment text -- separated from command line

`GetVerbs()` will return `[]string{"verb", "verb2"}`

`GetSeparated()` will return `[]string{"separated", "from", "command", "line"}`

The boolFlag will be set to true and the commentFlag will be set to "text".

## Struct flags

`AddFlagsForStruct` derives flags (and, via the usual env var fallback,
environment variables) from a struct's `flag:"..."` tags. Struct-typed and
pointer-to-struct-typed fields are supported and recursed into, to
arbitrary depth — their own tag becomes the prefix for the fields nested
inside them:

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

or with a defined struct type, reused by value or by pointer:

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

Both produce `-name -http-ip -http-port -grpc-ip -grpc-port` (prefixed by
whatever prefix is passed to `AddFlagsForStruct`). Environment variable
names are derived the same way as for any other flag: the flag name
upper-cased with `-` replaced by `_`, optionally prefixed via
`SetEnvPrefix`.

A pointer-to-struct field (`*Binding` above) is left `nil` by `Parse()`
unless at least one flag or environment variable inside it was actually
supplied; a value-typed field (`Binding`) is always populated. `Parse()` on
the container must be called after the package-level `flag.Parse()`, since
presence is determined in part by which flags were explicitly set on the
command line.

```go
a := &Config{}
c, err := flag.AddFlagsForStruct("app", a)
if err != nil {
    panic(err)
}
flag.Parse()
cfg := c.Parse()
```

Pointer-to-scalar leaf fields (`*string`, `*int`/`*int32`/..., `*bool`,
`*float64`/`*float32`, `*uint`/`*uint32`/...) are supported the same way:
`nil` unless the flag or its environment variable was explicitly provided,
regardless of nesting depth or whether the containing branch is itself a
value or a pointer:

```go
type Config struct {
    Timeout *int `flag:"timeout"`
}
```

`cfg.Timeout` stays `nil` unless `-timeout` (or the corresponding env var)
was set — even `-timeout=0` makes it a non-nil pointer to `0`.

`time.Duration` leaf fields are parsed as duration strings (`"5s"`,
`"1h30m"`, ...) rather than raw nanosecond counts, both as a value
(`time.Duration`, always populated like any other value leaf) and as a
pointer (`*time.Duration`, nil-unless-provided like any other pointer
leaf):

```go
type Config struct {
    Timeout  time.Duration  `flag:"timeout"`
    Deadline *time.Duration `flag:"deadline"`
}
```

`-timeout 5s` (or `TIMEOUT=5s`) sets `Timeout`; `Deadline` stays `nil`
unless `-deadline` (or `DEADLINE`) was explicitly given.

## Nilable flags

Outside of struct flags, `StringVarP`/`Int64VarP`/`BoolVarP`/
`Float64VarP`/`Uint64VarP`/`DurationVarP` give an ordinary variable the
same nil-unless-provided behavior. Unlike `StringVar`/`Int64Var`/etc.,
which write into an already-allocated `*T`, these take the address of a
pointer variable (`**T`) and leave it `nil` unless the flag was explicitly
set via CLI or environment variable:

```go
var timeout *int64
flag.Int64VarP(&timeout, "timeout", "timeout in seconds")
flag.Parse()
flag.ResolveP()
// timeout == nil unless -timeout or TIMEOUT was set
```

`ResolveP` must be called once, after `Parse`, and resolves every pointer
registered via a `*VarP` call up to that point.