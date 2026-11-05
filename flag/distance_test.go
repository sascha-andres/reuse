package flag

import "testing"

func TestOSADistance(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"", "", 0},
		{"", "run", 3},
		{"run", "", 3},
		{"run", "run", 0},
		{"rnu", "run", 1},
		{"lst", "list", 1},
		{"lst", "last", 1},
		{"list", "last", 1},
		{"help", "run", 4},
		{"ca", "abc", 3}, // OSA does not edit a swapped pair again
		{"Run", "run", 1},
	}
	for _, tt := range tests {
		t.Run(tt.a+"/"+tt.b, func(t *testing.T) {
			if got := osaDistance(tt.a, tt.b); got != tt.want {
				t.Fatalf("osaDistance(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
