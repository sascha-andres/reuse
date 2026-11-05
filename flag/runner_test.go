package flag

import (
	"context"
	"errors"
	"runtime"
	"slices"
	"strings"
	"testing"
)

var errBoom = errors.New("boom")

func TestRunVerbs(t *testing.T) {
	var got string
	mk := func(name string) Runner {
		return func(context.Context) error { got = name; return nil }
	}
	plain := func(name string) func(context.Context) error {
		return func(context.Context) error { got = name; return nil }
	}

	tests := []struct {
		name    string
		list    map[string]any
		args    []string
		want    string
		wantErr error
		same    bool // error must be returned unchanged (identity)
	}{
		{name: "empty args, default Runner", list: map[string]any{"": mk("def")}, want: "def"},
		{name: "empty args, default plain func", list: map[string]any{"": plain("def")}, want: "def"},
		{name: "empty args, default not a Runner", list: map[string]any{"": map[string]any{}}, wantErr: ErrInvalidVerbList},
		{name: "empty args, default wrong type", list: map[string]any{"": 1}, wantErr: ErrInvalidVerbList},
		{name: "empty args, no default", list: map[string]any{"a": mk("a")}, wantErr: ErrNothingToRun},
		{name: "nil args, no default", list: map[string]any{}, args: nil, wantErr: ErrNothingToRun},
		{name: "verb, Runner", list: map[string]any{"a": mk("a"), "": mk("def")}, args: []string{"a"}, want: "a"},
		{name: "verb, plain func", list: map[string]any{"a": plain("a")}, args: []string{"a"}, want: "a"},
		{name: "verb, untyped nil", list: map[string]any{"a": nil}, args: []string{"a"}, wantErr: ErrInvalidVerbList},
		{name: "verb, wrong type", list: map[string]any{"a": "x"}, args: []string{"a"}, wantErr: ErrInvalidVerbList},
		{name: "unknown verb", list: map[string]any{"a": mk("a"), "": mk("def")}, args: []string{"b"}, wantErr: ErrUnknownVerb},
		{name: "unknown verb, no default", list: map[string]any{}, args: []string{"b"}, wantErr: ErrUnknownVerb},
		{name: "not a verb, default Runner", list: map[string]any{"": mk("def")}, args: []string{"-x"}, want: "def"},
		{name: "not a verb, default plain func", list: map[string]any{"": plain("def")}, args: []string{"_x"}, want: "def"},
		{name: "not a verb, default not a Runner", list: map[string]any{"": map[string]any{}}, args: []string{"-x"}, wantErr: ErrInvalidVerbList},
		{name: "not a verb, no default", list: map[string]any{"a": mk("a")}, args: []string{"-x"}, wantErr: ErrNoDefaultRunner},
		{name: "empty string arg, no default", list: map[string]any{}, args: []string{""}, wantErr: ErrNoDefaultRunner},
		{
			name: "nested verb",
			list: map[string]any{"db": map[string]any{"up": mk("db up"), "": mk("db def")}},
			args: []string{"db", "up"},
			want: "db up",
		},
		{
			name: "nested, args run out, default of nested level",
			list: map[string]any{"db": map[string]any{"up": mk("db up"), "": mk("db def")}, "": mk("def")},
			args: []string{"db"},
			want: "db def",
		},
		{
			name:    "nested, args run out, no default at nested level",
			list:    map[string]any{"db": map[string]any{"up": mk("db up")}, "": mk("def")},
			args:    []string{"db"},
			wantErr: ErrNothingToRun,
		},
		{
			name:    "nested, unknown verb",
			list:    map[string]any{"db": map[string]any{"up": mk("db up")}},
			args:    []string{"db", "down"},
			wantErr: ErrUnknownVerb,
		},
		{
			name:    "nested, not a verb, no default at nested level",
			list:    map[string]any{"db": map[string]any{"up": mk("db up")}},
			args:    []string{"db", "-x"},
			wantErr: ErrNoDefaultRunner,
		},
		{
			name:    "nested, invalid entry",
			list:    map[string]any{"db": map[string]any{"up": 5}},
			args:    []string{"db", "up"},
			wantErr: ErrInvalidVerbList,
		},
		{
			name: "extra args after Runner are ignored",
			list: map[string]any{"a": mk("a")},
			args: []string{"a", "b"},
			want: "a",
		},
		{
			name:    "Runner error returned unchanged",
			list:    map[string]any{"a": Runner(func(context.Context) error { return errBoom })},
			args:    []string{"a"},
			wantErr: errBoom,
			same:    true,
		},
		{
			name:    "plain func error returned unchanged",
			list:    map[string]any{"": func(context.Context) error { return errBoom }},
			wantErr: errBoom,
			same:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got = ""
			err := RunVerbs(context.Background(), tt.list, tt.args)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if tt.same && err != tt.wantErr {
				t.Fatalf("error = %v, want identical %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("ran %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRunVerbs_TypedNilRunnerPanics(t *testing.T) {
	tests := []struct {
		name string
		list map[string]any
		args []string
	}{
		{"default nil Runner", map[string]any{"": Runner(nil)}, nil},
		{"verb nil Runner", map[string]any{"a": Runner(nil)}, []string{"a"}},
		{"verb nil plain func", map[string]any{"a": (func(context.Context) error)(nil)}, []string{"a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				err, ok := r.(runtime.Error)
				if !ok || !strings.Contains(err.Error(), "nil pointer dereference") {
					t.Fatalf("recovered %v, want runtime nil pointer dereference", r)
				}
			}()
			_ = RunVerbs(context.Background(), tt.list, tt.args)
		})
	}
}

func TestRunVerbs_UnknownVerbNamesVerb(t *testing.T) {
	err := RunVerbs(context.Background(), map[string]any{}, []string{"bulid"})
	if err == nil || err.Error() != `unknown verb: "bulid"` {
		t.Fatalf("got %v", err)
	}
}

func TestIsVerb(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{"", false},
		{"-x", false},
		{"--flag", false},
		{"_x", false},
		{"build", true},
		{"dry-run", true},
		{"list_all", true},
		{"abc123", true},
		{"123", true},
		{"A", true},
		{"größe", false},
		{"über", false},
		{"a b", false},
		{" a", false},
		{"a ", false},
		{"a=b", false},
		{"a.b", false},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := isVerb(tt.value); got != tt.want {
				t.Fatalf("isVerb(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestRunVerbs_Suggestions(t *testing.T) {
	noop := Runner(func(context.Context) error { return nil })
	top := map[string]any{"run": noop, "help": noop, "list": noop, "last": noop, "": noop, "-x": noop}
	nested := map[string]any{"remote": map[string]any{"add": noop, "rm": noop}}

	tests := []struct {
		name     string
		list     map[string]any
		args     []string
		opts     []VerbOption
		wantPath []string
		want     []string
	}{
		{name: "transposition", list: top, args: []string{"rnu"}, want: []string{"run"}},
		{name: "two matches, sorted", list: top, args: []string{"lst"}, want: []string{"last", "list"}},
		{name: "no match", list: top, args: []string{"xyz"}},
		{name: "empty and non-verb keys excluded", list: top, args: []string{"x"}},
		{name: "case counts as an edit", list: top, args: []string{"RUN"}},
		{name: "ignore case", list: top, args: []string{"RUN"}, opts: []VerbOption{WithIgnoreCase()}, want: []string{"run"}},
		{name: "ignore case plus edit", list: top, args: []string{"RNU"}, opts: []VerbOption{WithIgnoreCase()}, want: []string{"run"}},
		{
			name: "max distance 2, nearest first",
			list: top, args: []string{"lis"}, opts: []VerbOption{WithMaxDistance(2)},
			want: []string{"list", "last"},
		},
		{name: "negative max distance disables", list: top, args: []string{"rnu"}, opts: []VerbOption{WithMaxDistance(-1)}},
		{name: "nested", list: nested, args: []string{"remote", "ad"}, wantPath: []string{"remote"}, want: []string{"add"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunVerbs(context.Background(), tt.list, tt.args, tt.opts...)
			var uv *UnknownVerbError
			if !errors.As(err, &uv) || !errors.Is(err, ErrUnknownVerb) {
				t.Fatalf("error = %v, want *UnknownVerbError", err)
			}
			if uv.Verb != tt.args[len(tt.args)-1] {
				t.Fatalf("Verb = %q, want %q", uv.Verb, tt.args[len(tt.args)-1])
			}
			if !slices.Equal(uv.Path, tt.wantPath) {
				t.Fatalf("Path = %q, want %q", uv.Path, tt.wantPath)
			}
			if !slices.Equal(uv.Suggestions, tt.want) {
				t.Fatalf("Suggestions = %q, want %q", uv.Suggestions, tt.want)
			}
		})
	}
}
