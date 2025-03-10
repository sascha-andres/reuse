package validations

import "strconv"

// CreditCard validates a given credit card number using the Luhn algorithm and returns true if the card number is valid.
func CreditCard(number string) bool {
	if len(number) != 16 {
		return false
	}
	sum := 0
	for i := range number {
		result := 0
		n, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}
		if i%2 == 0 {
			n *= 2
		}
		result = n
		if n > 9 {
			result = 1 + (n - 10)
		}
		sum += result
	}
	return sum%10 == 0
}
