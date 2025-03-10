package validations

import "testing"

func TestCreditCard(t *testing.T) {
	tests := []struct {
		number   string
		expected bool
	}{
		{"4111111111111111", true},
		{"4111111111111110", false},
		{"4137894711755904", true},
	}
	t.Parallel()
	for _, test := range tests {
		t.Run(test.number, func(t *testing.T) {
			if result := CreditCard(test.number); result != test.expected {
				t.Fatalf("expected %t, got %t", test.expected, result)
			}
		})
	}
}
