package functional

import (
	"reflect"
	"testing"

	"golang.org/x/exp/maps"
	"slices"
)

var groupByTestCases = []struct {
	name string
	in   []map[string]string
	out  map[string][]map[string]string
	keys []string
	err  reflect.Type
}{
	{
		name: "simple",
		keys: []string{"name"},
		in: []map[string]string{
			{
				"name":   "alex",
				"friend": "berta",
			},
			{
				"name":   "berta",
				"friend": "alex",
			},
			{
				"name":   "berta",
				"friend": "cesar",
			},
		},
		out: map[string][]map[string]string{
			"berta": {
				{
					"name":   "berta",
					"friend": "alex",
				},
				{
					"name":   "berta",
					"friend": "cesar",
				},
			},
			"alex": {
				{
					"name":   "alex",
					"friend": "berta",
				},
			},
		},
		err: nil,
	},
	{
		name: "group_missing",
		keys: []string{"name"},
		in: []map[string]string{
			{
				"name":   "alex",
				"friend": "berta",
			},
			{
				"name":   "berta",
				"friend": "alex",
			},
			{
				"name":   "berta",
				"friend": "cesar",
			},
		},
		out: map[string][]map[string]string{
			"berta": {
				{
					"name":   "berta",
					"friend": "alex",
				},
				{
					"name":   "berta",
					"friend": "cesar",
				},
			},
			"alex": {
				{
					"name":   "alex",
					"friend": "berta",
				},
			},
		},
		err: nil,
	},
}

func TestGroupBy(t *testing.T) {
	for _, testCase := range groupByTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			testCase := testCase
			grouped, err := GroupBy(testCase.in, testCase.keys...)
			// TODO err case
			_ = err
			for key := range grouped {
				if !slices.Contains(maps.Keys(testCase.out), key) {
					t.Fatalf("expected %q not to be a group key", key)
				}
			}
			keys := maps.Keys(grouped)
			for _, k := range maps.Keys(testCase.out) {
				if !slices.Contains(keys, k) {
					t.Fatalf("expected %q to be a group key", k)
				}
			}
		})
	}
}
