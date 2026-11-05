package flag

import (
	"context"
	"errors"
	"fmt"
)

// Runner is an application run func called using verb list or direct
type Runner func(context.Context) error

var (
	ErrNothingToRun    = errors.New("no verb or default to run")
	ErrNoDefaultRunner = errors.New("no default runner given")
	ErrInvalidVerbList = errors.New("verblist map[string]any: any must be Runner or map[string]any")
	ErrUnknownVerb     = errors.New("unknown verb")
)

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
//   - args[0] is a verb and not a key: ErrUnknownVerb, wrapped with the verb.
//
// The error of a Runner is returned unchanged. A typed nil Runner or nil
// func(context.Context) error entry panics when it runs; an untyped nil entry
// gives ErrInvalidVerbList.
func RunVerbs(ctx context.Context, verbList map[string]any, args []string) error {
	if len(args) == 0 {
		return runDefault(ctx, verbList, ErrNothingToRun)
	}
	if !isVerb(args[0]) {
		return runDefault(ctx, verbList, ErrNoDefaultRunner)
	}
	v, ok := verbList[args[0]]
	if !ok {
		return fmt.Errorf("%w: %q", ErrUnknownVerb, args[0])
	}
	if sub, ok := v.(map[string]any); ok {
		return RunVerbs(ctx, sub, args[1:])
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
