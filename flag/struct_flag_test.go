package flag

import (
	f "flag"
	"os"
	"testing"
)

// resetForStructFlagTest gives each test a clean stdlib FlagSet and resets
// the package-level state that AddFlagsForStruct/Parse mutate, restoring
// everything afterwards so tests don't leak into each other.
func resetForStructFlagTest(t *testing.T) {
	t.Helper()
	oldCommandLine := f.CommandLine
	oldArgs := os.Args
	oldBooleanFlags := booleanFlags
	oldEnvPrefix := envPrefix
	oldOverridden := overriddenEnvPrefixes
	t.Cleanup(func() {
		f.CommandLine = oldCommandLine
		os.Args = oldArgs
		booleanFlags = oldBooleanFlags
		envPrefix = oldEnvPrefix
		overriddenEnvPrefixes = oldOverridden
	})
	f.CommandLine = f.NewFlagSet("test", f.ContinueOnError)
	booleanFlags = nil
	envPrefix = ""
	overriddenEnvPrefixes = nil
}

type binding struct {
	IP   string `flag:"ip"`
	Port uint   `flag:"port"`
}

type flatConfig struct {
	Name  string  `flag:"name"`
	Count int     `flag:"count"`
	Ratio float64 `flag:"ratio"`
	Big   uint64  `flag:"big"`
	On    bool    `flag:"on"`
}

func TestAddFlagsForStruct_FlatFields(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &flatConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-name", "svc", "-app-count", "3", "-app-ratio", "1.5", "-app-big", "42", "-app-on"}
	Parse()

	got := c.Parse()
	if got.Name != "svc" || got.Count != 3 || got.Ratio != 1.5 || got.Big != 42 || !got.On {
		t.Fatalf("unexpected result: %#v", got)
	}
}

// anonymousNestedConfig mirrors the inline-struct shape from the request:
// two sibling branches whose leaf fields share the same tag ("ip"/"port").
// Prior to the rework this collided because the value maps were keyed by
// leaf tag alone.
type anonymousNestedConfig struct {
	Name     string `flag:"name"`
	BindHttp struct {
		IP   string `flag:"ip"`
		Port uint   `flag:"port"`
	} `flag:"http"`
	BindGrpc struct {
		IP   string `flag:"ip"`
		Port uint   `flag:"port"`
	} `flag:"grpc"`
}

func TestAddFlagsForStruct_AnonymousNestedStruct_NoSiblingCollision(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &anonymousNestedConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{
		"cmd",
		"-app-name", "svc",
		"-app-http-ip", "127.0.0.1", "-app-http-port", "8080",
		"-app-grpc-ip", "10.0.0.1", "-app-grpc-port", "9090",
	}
	Parse()

	got := c.Parse()
	if got.Name != "svc" {
		t.Errorf("Name = %q, want %q", got.Name, "svc")
	}
	if got.BindHttp.IP != "127.0.0.1" || got.BindHttp.Port != 8080 {
		t.Errorf("BindHttp = %+v, want IP=127.0.0.1 Port=8080", got.BindHttp)
	}
	if got.BindGrpc.IP != "10.0.0.1" || got.BindGrpc.Port != 9090 {
		t.Errorf("BindGrpc = %+v, want IP=10.0.0.1 Port=9090", got.BindGrpc)
	}
}

// namedNestedConfig mirrors the second shape from the request: a defined
// struct used by value for one branch and by pointer for the other.
type namedNestedConfig struct {
	Name     string   `flag:"name"`
	BindHttp binding  `flag:"http"`
	BindGrpc *binding `flag:"grpc"`
}

func TestAddFlagsForStruct_NamedNestedStruct_ValueAlwaysPopulated(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &namedNestedConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-http-ip", "127.0.0.1", "-app-http-port", "8080"}
	Parse()

	got := c.Parse()
	if got.BindHttp.IP != "127.0.0.1" || got.BindHttp.Port != 8080 {
		t.Errorf("BindHttp = %+v, want IP=127.0.0.1 Port=8080", got.BindHttp)
	}
}

func TestAddFlagsForStruct_NamedNestedStruct_PointerPresentViaFlag(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &namedNestedConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-grpc-ip", "10.0.0.1", "-app-grpc-port", "9090"}
	Parse()

	got := c.Parse()
	if got.BindGrpc == nil {
		t.Fatal("BindGrpc = nil, want non-nil")
	}
	if got.BindGrpc.IP != "10.0.0.1" || got.BindGrpc.Port != 9090 {
		t.Errorf("BindGrpc = %+v, want IP=10.0.0.1 Port=9090", *got.BindGrpc)
	}
}

func TestAddFlagsForStruct_NamedNestedStruct_PointerAbsentStaysNil(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &namedNestedConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-name", "svc"}
	Parse()

	got := c.Parse()
	if got.BindGrpc != nil {
		t.Errorf("BindGrpc = %+v, want nil", *got.BindGrpc)
	}
}

func TestAddFlagsForStruct_NamedNestedStruct_PointerPresentViaEnv(t *testing.T) {
	resetForStructFlagTest(t)
	setEnv(t, "APP_GRPC_IP", "10.0.0.5")

	cfg := &namedNestedConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd"}
	Parse()

	got := c.Parse()
	if got.BindGrpc == nil {
		t.Fatal("BindGrpc = nil, want non-nil")
	}
	if got.BindGrpc.IP != "10.0.0.5" {
		t.Errorf("BindGrpc.IP = %q, want %q", got.BindGrpc.IP, "10.0.0.5")
	}
	if got.BindGrpc.Port != 0 {
		t.Errorf("BindGrpc.Port = %d, want 0 (never provided)", got.BindGrpc.Port)
	}
}

// Three levels: a value struct wrapping a pointer struct wrapping a leaf.
type deepInner struct {
	Value string `flag:"value"`
}
type deepMiddle struct {
	Inner *deepInner `flag:"inner"`
}
type deepOuter struct {
	Middle deepMiddle `flag:"mid"`
}

func TestAddFlagsForStruct_MultiLevelNesting(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &deepOuter{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-mid-inner-value", "hello"}
	Parse()

	got := c.Parse()
	if got.Middle.Inner == nil {
		t.Fatal("Middle.Inner = nil, want non-nil")
	}
	if got.Middle.Inner.Value != "hello" {
		t.Errorf("Middle.Inner.Value = %q, want %q", got.Middle.Inner.Value, "hello")
	}
}

// A pointer-to-struct branch that itself contains a pointer-to-struct
// branch: presence must propagate through both levels.
type ptrInPtrInner struct {
	Value string `flag:"value"`
}
type ptrInPtrOuter struct {
	Inner *ptrInPtrInner `flag:"inner"`
}
type ptrInPtrConfig struct {
	Outer *ptrInPtrOuter `flag:"outer"`
}

func TestAddFlagsForStruct_PointerInPointerPropagatesPresence(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &ptrInPtrConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-outer-inner-value", "hi"}
	Parse()

	got := c.Parse()
	if got.Outer == nil {
		t.Fatal("Outer = nil, want non-nil")
	}
	if got.Outer.Inner == nil {
		t.Fatal("Outer.Inner = nil, want non-nil")
	}
	if got.Outer.Inner.Value != "hi" {
		t.Errorf("Outer.Inner.Value = %q, want %q", got.Outer.Inner.Value, "hi")
	}
}

func TestAddFlagsForStruct_PointerInPointerAbsentStaysNil(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &ptrInPtrConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd"}
	Parse()

	got := c.Parse()
	if got.Outer != nil {
		t.Errorf("Outer = %+v, want nil", *got.Outer)
	}
}
