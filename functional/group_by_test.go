package functional

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"slices"

	"golang.org/x/exp/maps"
)

func TestGroupBy(t *testing.T) {
	var testCases = []struct {
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

	t.Parallel()

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
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

func TestGroupByFunc(t *testing.T) {
	var testCases = []struct {
		name     string
		in       []time.Time
		out      map[string][]time.Time
		hasError bool
		gf       KeyFunc[time.Time, string]
	}{
		{
			name: "simple",
			in:   []time.Time{time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
			out: map[string][]time.Time{
				"2021": {time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
				"2022": {time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				"2023": {time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			hasError: false,
			gf: func(t time.Time) (string, error) {
				return fmt.Sprintf("%d", t.Year()), nil
			},
		},
		{
			name: "error",
			in:   []time.Time{time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
			out: map[string][]time.Time{
				"2021": {time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
				"2022": {time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				"2023": {time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			hasError: true,
			gf: func(t time.Time) (string, error) {
				if t.Year() == 2023 {
					return "", fmt.Errorf("test")
				}
				return fmt.Sprintf("%d", t.Year()), nil
			},
		},
	}

	t.Parallel()

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			result, err := GroupByFunc(testCase.in, testCase.gf)
			if err != nil {
				if !testCase.hasError {
					t.Fatalf("unexpected error: %s", err.Error())
				}
				if result != nil {
					t.Fatalf("expected nil result but got %v", result)
				}
				return
			}
			if len(result) != len(testCase.out) {
				t.Fatalf("expected %d groups but got %d", len(testCase.out), len(result))
			}
			for k := range result {
				if _, ok := testCase.out[k]; !ok {
					t.Fatalf("expected %q to be a group key [result]", k)
				}
				if len(result[k]) != len(testCase.out[k]) {
					t.Fatalf("expected %d items in group %q but got %d [result]", len(testCase.out[k]), k, len(result[k]))
				}
				if !cmp.Equal(result[k], testCase.out[k]) {
					t.Fatalf("expected %v but got %v [result]", testCase.out[k], result[k])
				}
			}
			for k := range testCase.out {
				if _, ok := result[k]; !ok {
					t.Fatalf("expected %q to be a group key [expectation]", k)
				}
				if len(result[k]) != len(testCase.out[k]) {
					t.Fatalf("expected %d items in group %q but got %d [expectation]", len(result[k]), k, len(testCase.out[k]))
				}
				if !cmp.Equal(testCase.out[k], result[k]) {
					t.Fatalf("expected %v but got %v [expectation]", result[k], testCase.out[k])
				}
			}
		})
	}
}
