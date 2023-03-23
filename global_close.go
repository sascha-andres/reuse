package reuse

// Close will just discard the error returned
//
// best used for defer calls:
// defer file.Close() -> an error unhandled message will be displayed by some dev envs
// defer reuse.Close(file.Close) -> all ok
func Close(closer func() error) { _ = closer() }
