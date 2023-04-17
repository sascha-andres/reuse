package flag

import (
	"os"
	"testing"
)

func TestMap(t *testing.T) {
	_ = os.Setenv("ENV_TEST_VAL1", "1")
	_ = os.Setenv("ENV_TEST_VAL2", "2")
	defer func() { _ = os.Unsetenv("ENV_TEST_VAL1") }()
	defer func() { _ = os.Unsetenv("ENV_TES2_VAL1") }()

	data := GetMapForEnvironmentVariable("ENV_TEST")
	if data["val1"] != "1" {
		t.Fatalf("expected val1 to be 1, got %s", data["val1"])
	}
	if data["val2"] != "2" {
		t.Fatalf("expected val2 to be 2, got %s", data["val2"])
	}
}
