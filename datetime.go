package reuse

import (
	"time"
)

// Date returns a string representation of the date in format YYYY-MM-DD
// It is RFC 3339, ISO-8601 and W3C(HTML) compliant
func Date(t time.Time) string {
	return t.Format("2006-01-02")
}

// DateTime returns a string representation of the date and time in format YYYY-MM-DDTHH:MM:SSZ
// It is RFC 3339, ISO-8601 and W3C(HTML) compliant
func DateTime(t time.Time) string {
	return t.Format("2006-01-02T15:04:05Z07:00")
}

// DateTimeWithTimeZone returns a string representation of the date and time in format YYYY-MM-DDTHH:MM:SSZ07:00
// It is RFC 3339, ISO-8601 and W3C(HTML) compliant
func DateTimeWithTimeZone(t time.Time) string {
	return t.Format(time.RFC3339)
}

// CurrentDate returns a string representation of the current date in format YYYY-MM-DD
// It is RFC 3339, ISO-8601 and W3C(HTML) compliant
func CurrentDate() string {
	return Date(time.Now())
}

// CurrentDateTime returns a string representation of the current date and time in format YYYY-MM-DDTHH:MM:SSZ
// It is RFC 3339, ISO-8601 and W3C(HTML) compliant
func CurrentDateTime() string {
	return DateTime(time.Now())
}

// CurrentDateTimeWithTimeZone returns a string representation of the current date and time in format YYYY-MM-DDTHH:MM:SSZ07:00
// It is RFC 3339, ISO-8601 and W3C(HTML) compliant
func CurrentDateTimeWithTimeZone() string {
	return DateTimeWithTimeZone(time.Now())
}

// UTCCurrentDateTimeWithTimeZone returns a string representation of the current date and time (UTC) in format YYYY-MM-DDTHH:MM:SSZ
// It is RFC 3339, ISO-8601 and W3C(HTML) compliant
func UTCCurrentDateTimeWithTimeZone() string {
	return DateTimeWithTimeZone(time.Now().UTC())
}

// UTCCurrentDate returns a string representation of the current date (UTC) in format YYYY-MM-DD
// It is RFC 3339, ISO-8601 and W3C(HTML) compliant
func UTCCurrentDate() string {
	return Date(time.Now().UTC())
}

// UTCCurrentDateTime returns a string representation of the current date and time (UTC) in format YYYY-MM-DDTHH:MM:SSZ
// It is RFC 3339, ISO-8601 and W3C(HTML) compliant
func UTCCurrentDateTime() string {
	return DateTime(time.Now().UTC())
}
