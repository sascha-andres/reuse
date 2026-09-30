# j

One helper: read and JSON-decode a file into a Go type.

```go
type Config struct {
    Name string `json:"name"`
}

cfg, err := j.UnmarshalFile[Config]("config.json")
```

Returns `os.ErrNotExist` (via `reuse.FileExists`) if the file doesn't
exist, rather than whatever `os.ReadFile` would produce.
