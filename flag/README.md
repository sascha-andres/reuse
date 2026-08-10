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

## Separated

If you want to pass arguments for something like a sub command you can use the separate feature. Activate it using `flag.SetSeparate()`. Everything after `--` will be treated as a separate from command line and not parsed as verbs or flags. 

    cmd verb -bool verb2 -comment text -- separated from command line

`GetVerbs()` will return `[]string{"verb", "verb2"}`

`GetSeparate()` will return `[]string{"separated", "from", "command", "line"}`

`GetBool("bool")` will return `true`

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