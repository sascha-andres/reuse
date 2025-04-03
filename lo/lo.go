package lo

import (
	"log/slog"
)

// MustGet executes the provided function, logs an error if it occurs, and returns the value. Logging uses the given logger.
func MustGet[T any](logger *slog.Logger, msg string, f func() (T, error)) T {
	r, err := f()
	if err != nil {
		logger.Error(msg, "err", err)
	}
	return r
}

// MustRun executes the provided function f and logs an error with the given message if f returns an error.
func MustRun(logger *slog.Logger, msg string, f func() error) {
	err := f()
	if err != nil {
		logger.Error(msg, "err", err)
	}
}
