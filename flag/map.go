package flag

import (
	"os"
	"strings"
)

// GetMapForEnvironmentVariable returns a map of strings from env variables
// that start with the given name.
func GetMapForEnvironmentVariable(name string) map[string]string {
	environmentVariables := os.Environ()
	result := make(map[string]string)
	for _, environmentVariable := range environmentVariables {
		if strings.HasPrefix(environmentVariable, name) {
			parts := strings.SplitN(environmentVariable, "=", 2)
			if len(parts) == 2 {
				parts[0] = strings.TrimPrefix(parts[0], name+"_")
				parts[0] = strings.ToLower(parts[0])
				result[parts[0]] = parts[1]
			}
		}
	}
	return result
}
