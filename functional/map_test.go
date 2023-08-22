package functional

import (
	"fmt"
	"testing"

	"slices"
)

var mapTestCases = []struct {
	name string
	in   []any
	out  []any
	err  error
	f    func(any) (any, error)
}{
	{
		name: "simple",
		in:   []any{1},
		out:  []any{"str 1"},
		err:  nil,
		f: func(val any) (any, error) {
			return fmt.Sprintf("str %v", val), nil
		},
	},
}

func TestMap(t *testing.T) {
	for _, test := range mapTestCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			res, err := Map(test.in, test.f)
			if err != nil {
				if test.err == nil {
					t.Fatalf("expected no error, received %q", err)
				} else {
					if test.err != err {
						t.Fatalf("expected %q, received %q", test.err, err)
					}
				}
			}
			for _, re := range res {
				if !slices.Contains(test.out, re) {
					t.Fatalf("expected %q to be in result", re)
				}
			}
			for _, a := range test.out {
				if !slices.Contains(res, a) {
					t.Fatalf("missing %q from results", a)
				}
			}
		})
	}
}
