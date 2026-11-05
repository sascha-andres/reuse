package flag

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// Runner is an application run func called using verb list or direct
type Runner func(context.Context) error

var (
	ErrNothingToRun    = errors.New("no verb or default to run")
	ErrNoDefaultRunner = errors.New("no default runner given")
	ErrInvalidVerbList = errors.New("verblist map[string]any: any must be Runner or map[string]any")
	ErrUnknownVerb     = errors.New("unknown verb")
)

// UnknownVerbError is returned by RunVerbs when a verb is not a key of its
// level. errors.Is(err, ErrUnknownVerb) is true.
type UnknownVerbError struct {
	Path        []string // verbs before Verb, outermost first
	Verb        string
	Suggestions []string // verb keys of the level within the max distance, nearest first
}

func (e *UnknownVerbError) Error() string { return fmt.Sprintf("%v: %q", ErrUnknownVerb, e.Verb) }

func (e *UnknownVerbError) Unwrap() error { return ErrUnknownVerb }

// VerbOption changes how RunVerbs finds suggestions for an unknown verb.
type VerbOption func(*verbOptions)

type verbOptions struct {
	maxDistance int
	ignoreCase  bool
}

// WithMaxDistance sets the max edit distance of a suggestion. The default is
// 1. A negative value disables suggestions.
func WithMaxDistance(n int) VerbOption {
	return func(o *verbOptions) { o.maxDistance = n }
}

// WithIgnoreCase compares the verb and the keys without case to find
// suggestions. The lookup of the verb stays case-sensitive.
func WithIgnoreCase() VerbOption {
	return func(o *verbOptions) { o.ignoreCase = true }
}

// RunVerbs takes a nested dict of string:any and runs the matching Runner.
// args is normally the result of GetVerbs(). A Runner is a Runner or a plain
// func(context.Context) error. The key "" is the default of a level.
//
//   - args is empty or args[0] is not a verb: the "" entry runs. If "" is
//     missing, the result is ErrNothingToRun (args empty) or
//     ErrNoDefaultRunner (args not empty). If "" is not a Runner, the result
//     is ErrInvalidVerbList.
//   - args[0] is a verb and a key: a Runner runs, a map[string]any is used
//     with args[1:] at the next level, any other type gives ErrInvalidVerbList.
//   - args[0] is a verb and not a key: an *UnknownVerbError that wraps
//     ErrUnknownVerb. Its Suggestions are the verb keys of that level within
//     the optimal string alignment distance of opts (default 1).
//
// The error of a Runner is returned unchanged. A typed nil Runner or nil
// func(context.Context) error entry panics when it runs; an untyped nil entry
// gives ErrInvalidVerbList.
func RunVerbs(ctx context.Context, verbList map[string]any, args []string, opts ...VerbOption) error {
	o := verbOptions{maxDistance: 1}
	for _, opt := range opts {
		opt(&o)
	}
	return runVerbs(ctx, verbList, args, nil, o)
}

// runVerbs is RunVerbs for one level. path is the verbs of the levels above.
func runVerbs(ctx context.Context, verbList map[string]any, args, path []string, o verbOptions) error {
	if len(args) == 0 {
		return runDefault(ctx, verbList, ErrNothingToRun)
	}
	if !isVerb(args[0]) {
		return runDefault(ctx, verbList, ErrNoDefaultRunner)
	}
	v, ok := verbList[args[0]]
	if !ok {
		return &UnknownVerbError{Path: path, Verb: args[0], Suggestions: suggestVerbs(verbList, args[0], o)}
	}
	if sub, ok := v.(map[string]any); ok {
		return runVerbs(ctx, sub, args[1:], append(path, args[0]), o)
	}
	if r, ok := asRunner(v); ok {
		return r(ctx)
	}
	return ErrInvalidVerbList
}

// runDefault runs the "" entry of verbList. missing is returned if there is no
// "" entry.
func runDefault(ctx context.Context, verbList map[string]any, missing error) error {
	v, ok := verbList[""]
	if !ok {
		return missing
	}
	if r, ok := asRunner(v); ok {
		return r(ctx)
	}
	return ErrInvalidVerbList
}

// suggestVerbs returns the verb keys of verbList within o.maxDistance of verb,
// sorted by distance, then by name.
func suggestVerbs(verbList map[string]any, verb string, o verbOptions) []string {
	if o.ignoreCase {
		verb = strings.ToLower(verb)
	}
	dist := map[string]int{}
	var out []string
	for key := range verbList {
		if !isVerb(key) { // a key that is not a verb can never run
			continue
		}
		k := key
		if o.ignoreCase {
			k = strings.ToLower(key)
		}
		if d := osaDistance(verb, k); d <= o.maxDistance {
			dist[key] = d
			out = append(out, key)
		}
	}
	slices.SortFunc(out, func(a, b string) int {
		return cmp.Or(cmp.Compare(dist[a], dist[b]), strings.Compare(a, b))
	})
	return out
}

func asRunner(v any) (Runner, bool) {
	switch r := v.(type) {
	case Runner:
		return r, true
	case func(context.Context) error:
		return r, true
	}
	return nil, false
}

// isVerb returns true if value is a verb: not empty, the first character is
// ASCII [A-Za-z0-9] and all other characters are ASCII [A-Za-z0-9_-].
func isVerb(value string) bool {
	if value == "" {
		return false
	}
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
		case i > 0 && (c == '-' || c == '_'):
		default:
			return false
		}
	}
	return true
}
