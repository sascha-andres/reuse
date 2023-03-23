package reuse

// DiscardError will just discard the error returned
//
// best used for defer calls:
// defer file.Close() -> an error unhandled message will be displayed by some dev envs
// defer reuse.DiscardError(file.Close) -> all ok
func DiscardError(closer func() error) { _ = closer() }
