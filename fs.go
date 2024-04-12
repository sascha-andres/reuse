package reuse

import "os"

// DirectoryExists checks whether a directory exists at the given directory path.
// It returns true if the directory exists, otherwise false.
func DirectoryExists(directory string) bool {
	if i, err := os.Stat(directory); err == nil && i.IsDir() {
		return true
	}
	return false
}

// FileExists checks whether a file or directory exists at the given file path.
// It returns true if the file or directory exists, otherwise false.
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}
