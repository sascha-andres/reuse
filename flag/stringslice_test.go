package flag

import (
	"flag"
	"os"
	"testing"
)

func resetCommandLine(t *testing.T) {
	t.Helper()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

func TestStringSliceWithoutEnv_Default(t *testing.T) {
	resetCommandLine(t)
	get := StringSliceWithoutEnv("tags", []string{"a", "b"}, "usage")
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		t.Fatal(err)
	}
	got := get()
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("expected [a b], got %v", got)
	}
}

func TestStringSliceWithoutEnv_FlagValue(t *testing.T) {
	resetCommandLine(t)
	get := StringSliceWithoutEnv("tags", []string{"a", "b"}, "usage")
	if err := flag.CommandLine.Parse([]string{"-tags", "x,y,z"}); err != nil {
		t.Fatal(err)
	}
	got := get()
	if len(got) != 3 || got[0] != "x" || got[1] != "y" || got[2] != "z" {
		t.Errorf("expected [x y z], got %v", got)
	}
}

func TestStringSliceWithoutEnv_IgnoresEnv(t *testing.T) {
	resetCommandLine(t)
	setEnv(t, "TAGS", "env1,env2")
	get := StringSliceWithoutEnv("tags", []string{"default"}, "usage")
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		t.Fatal(err)
	}
	got := get()
	// WithoutEnv should ignore TAGS env var and use the provided default
	if len(got) != 1 || got[0] != "default" {
		t.Errorf("expected [default], got %v", got)
	}
}

func TestStringSliceWithoutEnv_SingleValue(t *testing.T) {
	resetCommandLine(t)
	get := StringSliceWithoutEnv("items", []string{"only"}, "usage")
	if err := flag.CommandLine.Parse([]string{"-items", "single"}); err != nil {
		t.Fatal(err)
	}
	got := get()
	if len(got) != 1 || got[0] != "single" {
		t.Errorf("expected [single], got %v", got)
	}
}

func TestStringSlice_DefaultValue(t *testing.T) {
	resetCommandLine(t)
	envPrefix = ""
	get := StringSlice("fruits", []string{"apple", "banana"}, "usage")
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		t.Fatal(err)
	}
	got := get()
	if len(got) != 2 || got[0] != "apple" || got[1] != "banana" {
		t.Errorf("expected [apple banana], got %v", got)
	}
}

func TestStringSlice_EnvOverridesDefault(t *testing.T) {
	resetCommandLine(t)
	envPrefix = ""
	setEnv(t, "FRUITS", "mango,kiwi,grape")
	get := StringSlice("fruits", []string{"apple"}, "usage")
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		t.Fatal(err)
	}
	got := get()
	if len(got) != 3 || got[0] != "mango" || got[1] != "kiwi" || got[2] != "grape" {
		t.Errorf("expected [mango kiwi grape], got %v", got)
	}
}

func TestStringSlice_FlagOverridesEnv(t *testing.T) {
	resetCommandLine(t)
	envPrefix = ""
	setEnv(t, "FRUITS", "mango,kiwi")
	get := StringSlice("fruits", []string{"apple"}, "usage")
	if err := flag.CommandLine.Parse([]string{"-fruits", "pear,peach"}); err != nil {
		t.Fatal(err)
	}
	got := get()
	if len(got) != 2 || got[0] != "pear" || got[1] != "peach" {
		t.Errorf("expected [pear peach], got %v", got)
	}
}

func TestStringSlice_WithEnvPrefix(t *testing.T) {
	resetCommandLine(t)
	envPrefix = "APP"
	defer func() { envPrefix = "" }()
	setEnv(t, "APP_COLORS", "red,green,blue")
	get := StringSlice("colors", []string{"white"}, "usage")
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		t.Fatal(err)
	}
	got := get()
	if len(got) != 3 || got[0] != "red" || got[1] != "green" || got[2] != "blue" {
		t.Errorf("expected [red green blue], got %v", got)
	}
}

func TestStringSliceVarWithoutEnv_Default(t *testing.T) {
	resetCommandLine(t)
	get := StringSliceVarWithoutEnv("modes", []string{"fast", "slow"}, "usage")
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		t.Fatal(err)
	}
	got := get()
	if len(got) != 2 || got[0] != "fast" || got[1] != "slow" {
		t.Errorf("expected [fast slow], got %v", got)
	}
}

func TestStringSliceVarWithoutEnv_FlagValue(t *testing.T) {
	resetCommandLine(t)
	get := StringSliceVarWithoutEnv("modes", []string{"fast"}, "usage")
	if err := flag.CommandLine.Parse([]string{"-modes", "turbo,ultra"}); err != nil {
		t.Fatal(err)
	}
	got := get()
	if len(got) != 2 || got[0] != "turbo" || got[1] != "ultra" {
		t.Errorf("expected [turbo ultra], got %v", got)
	}
}

func TestStringSliceVarWithoutEnv_IgnoresEnv(t *testing.T) {
	resetCommandLine(t)
	setEnv(t, "MODES", "env-a,env-b")
	get := StringSliceVarWithoutEnv("modes", []string{"default-mode"}, "usage")
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		t.Fatal(err)
	}
	got := get()
	if len(got) != 1 || got[0] != "default-mode" {
		t.Errorf("expected [default-mode], got %v", got)
	}
}

func TestStringSliceVar_EnvOverridesDefault(t *testing.T) {
	resetCommandLine(t)
	envPrefix = ""
	setEnv(t, "MODES", "alpha,beta")
	get := StringSliceVar("modes", []string{"default"}, "usage")
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		t.Fatal(err)
	}
	got := get()
	if len(got) != 2 || got[0] != "alpha" || got[1] != "beta" {
		t.Errorf("expected [alpha beta], got %v", got)
	}
}

func TestStringSliceVar_FlagValue(t *testing.T) {
	resetCommandLine(t)
	envPrefix = ""
	get := StringSliceVar("modes", []string{"default"}, "usage")
	if err := flag.CommandLine.Parse([]string{"-modes", "x,y"}); err != nil {
		t.Fatal(err)
	}
	got := get()
	if len(got) != 2 || got[0] != "x" || got[1] != "y" {
		t.Errorf("expected [x y], got %v", got)
	}
}

func TestStringSlice_EmptyDefault(t *testing.T) {
	resetCommandLine(t)
	envPrefix = ""
	get := StringSliceWithoutEnv("empty", []string{}, "usage")
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		t.Fatal(err)
	}
	got := get()
	// strings.Split("", ",") returns [""] — one empty string element
	if len(got) != 1 || got[0] != "" {
		t.Errorf("expected [\"\"], got %v", got)
	}
}
