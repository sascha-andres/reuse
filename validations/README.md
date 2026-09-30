# validations

One helper: `CreditCard` validates a 16-digit credit card number string
using the Luhn algorithm.

```go
validations.CreditCard("4111111111111111") // true
validations.CreditCard("1234567890123456") // false
```

Requires exactly 16 digits; any other length or a non-digit character
returns `false`.
