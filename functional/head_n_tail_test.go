package functional

import "testing"

// TestIntArray verifies the behavior of Head and Tail functions with an integer slice as input.
func TestIntArray(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}
	out := Head(in)
	if out != 1 {
		t.Fatalf("expected 1, got %d", out)
	}
	tail := Tail(in)
	if len(tail) != 4 {
		t.Fatalf("expected 4, got %d", len(tail))
	}
	if tail[0] != 2 {
		t.Fatalf("expected 2, got %d", tail[0])
	}
	if tail[1] != 3 {
		t.Fatalf("expected 3, got %d", tail[1])
	}
	if tail[2] != 4 {
		t.Fatalf("expected 4, got %d", tail[2])
	}
	if tail[3] != 5 {
		t.Fatalf("expected 5, got %d", tail[3])
	}
}
