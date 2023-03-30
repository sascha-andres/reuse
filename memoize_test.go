package reuse

import (
	"errors"
	"testing"
)

func TestMemoization(t *testing.T) {
	g := Memoize(func(t int) (int, error) {
		if t%2 == 0 {
			return 1, nil
		}
		if t == 3 {
			return 0, errors.New("just to test error case")
		}
		return 0, nil
	})

	v, cache, err := g(1)
	if err != nil {
		t.Errorf("[1] expected no error, got %v", err)
	}
	if cache {
		t.Error("[1] on first query cache should be empty")
	}
	if v != 0 {
		t.Error("[1] my func is wrong")
	}
	v, cache, err = g(2)
	if err != nil {
		t.Errorf("[2] expected no error, got %v", err)
	}
	if cache {
		t.Error("[2] on first query cache should not contain result for2")
	}
	if v != 1 {
		t.Error("[2] my func is wrong")
	}
	v, cache, err = g(1)
	if err != nil {
		t.Errorf("[1:2] expected no error, got %v", err)
	}
	if !cache {
		t.Error("[1:2] on second query cache should be used")
	}
	if v != 0 {
		t.Error("[1:2] my func is wrong")
	}
	v, cache, err = g(3)
	if err == nil {
		t.Error("[3] expected error")
	}
}
