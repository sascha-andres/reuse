package j

import (
	"encoding/json"
	"os"

	"github.com/sascha-andres/reuse"
)

// UnmarshalFile reads a JSON file, parses its content into a Go type T, and returns a pointer to the unmarshalled object.
func UnmarshalFile[T any](filename string) (*T, error) {
	if !reuse.FileExists(filename) {
		return nil, os.ErrNotExist
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var a T
	err = json.Unmarshal(data, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
