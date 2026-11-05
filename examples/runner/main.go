package main

import (
	"context"
	"fmt"
	"os"
	"slices"

	"go.livingit.de/reuse/flag"
)

// printVerbs prints the keys of a level, sorted, one per line. A map child is
// printed with a trailing " ...".
func printVerbs(level map[string]any) {
	keys := make([]string, 0, len(level))
	for key := range level {
		if key != "" {
			keys = append(keys, key)
		}
	}
	slices.Sort(keys)
	for _, key := range keys {
		if _, ok := level[key].(map[string]any); ok {
			fmt.Printf("%s ...\n", key)
			continue
		}
		fmt.Println(key)
	}
}

func run(ctx context.Context) error {
	var name string
	var loud bool

	flag.SetEnvPrefix("RUNNER")
	flag.StringVar(&name, "name", "world", "name to greet")
	flag.BoolVar(&loud, "loud", false, "greet loudly")

	flag.Parse()

	greet := map[string]any{
		"hello": func(_ context.Context) error {
			if loud {
				fmt.Printf("HELLO, %s!\n", name)
				return nil
			}
			fmt.Printf("hello, %s\n", name)
			return nil
		},
		"bye": func(_ context.Context) error {
			fmt.Printf("bye, %s\n", name)
			return nil
		},
	}
	greet[""] = func(_ context.Context) error {
		printVerbs(greet)
		return nil
	}

	tree := map[string]any{
		"greet": greet,
		"version": func(_ context.Context) error {
			fmt.Println("runner 1.0.0")
			return nil
		},
	}
	tree[""] = func(_ context.Context) error {
		printVerbs(tree)
		return nil
	}

	return flag.RunVerbs(ctx, tree, flag.GetVerbs())
}

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
