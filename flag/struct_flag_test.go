package flag

import (
	f "flag"
	"os"
	"testing"
	"time"
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

type pointerScalarConfig struct {
	Timeout *int   `flag:"timeout"`
	Label   string `flag:"label"`
}

func TestAddFlagsForStruct_PointerScalarLeaf_PresentViaFlag(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &pointerScalarConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-timeout", "30", "-app-label", "svc"}
	Parse()

	got := c.Parse()
	if got.Timeout == nil {
		t.Fatal("Timeout = nil, want non-nil")
	}
	if *got.Timeout != 30 {
		t.Errorf("*Timeout = %d, want 30", *got.Timeout)
	}
}

func TestAddFlagsForStruct_PointerScalarLeaf_PresentViaEnv(t *testing.T) {
	resetForStructFlagTest(t)
	setEnv(t, "APP_TIMEOUT", "45")

	cfg := &pointerScalarConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd"}
	Parse()

	got := c.Parse()
	if got.Timeout == nil {
		t.Fatal("Timeout = nil, want non-nil")
	}
	if *got.Timeout != 45 {
		t.Errorf("*Timeout = %d, want 45", *got.Timeout)
	}
}

func TestAddFlagsForStruct_PointerScalarLeaf_AbsentStaysNil(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &pointerScalarConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-label", "svc"}
	Parse()

	got := c.Parse()
	if got.Timeout != nil {
		t.Errorf("Timeout = %d, want nil", *got.Timeout)
	}
	if got.Label != "svc" {
		t.Errorf("Label = %q, want %q", got.Label, "svc")
	}
}

// A pointer-scalar leaf nested inside a value-typed branch and a
// pointer-typed branch: its nil-ness is independent of the branch's.
type pointerScalarInValueBranch struct {
	Retries *int `flag:"retries"`
}
type pointerScalarInPointerBranch struct {
	Retries *int `flag:"retries"`
}
type nestedPointerScalarConfig struct {
	ValueBranch   pointerScalarInValueBranch    `flag:"value"`
	PointerBranch *pointerScalarInPointerBranch `flag:"ptr"`
}

func TestAddFlagsForStruct_PointerScalarLeaf_NestedInValueBranch(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &nestedPointerScalarConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-value-retries", "3"}
	Parse()

	got := c.Parse()
	if got.ValueBranch.Retries == nil || *got.ValueBranch.Retries != 3 {
		t.Errorf("ValueBranch.Retries = %+v, want *3", got.ValueBranch.Retries)
	}
	if got.PointerBranch != nil {
		t.Errorf("PointerBranch = %+v, want nil (its own leaf never provided)", *got.PointerBranch)
	}
}

func TestAddFlagsForStruct_PointerScalarLeaf_NestedInPointerBranch(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &nestedPointerScalarConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-ptr-retries", "5"}
	Parse()

	got := c.Parse()
	if got.ValueBranch.Retries != nil {
		t.Errorf("ValueBranch.Retries = %d, want nil", *got.ValueBranch.Retries)
	}
	if got.PointerBranch == nil {
		t.Fatal("PointerBranch = nil, want non-nil")
	}
	if got.PointerBranch.Retries == nil || *got.PointerBranch.Retries != 5 {
		t.Errorf("PointerBranch.Retries = %+v, want *5", got.PointerBranch.Retries)
	}
}

type durationConfig struct {
	Timeout  time.Duration  `flag:"timeout"`
	Deadline *time.Duration `flag:"deadline"`
}

func TestAddFlagsForStruct_DurationValueLeaf_ParsesDurationString(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &durationConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-timeout", "5s"}
	Parse()

	got := c.Parse()
	if got.Timeout != 5*time.Second {
		t.Errorf("Timeout = %v, want 5s", got.Timeout)
	}
}

func TestAddFlagsForStruct_DurationPointerLeaf_PresentViaFlag(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &durationConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-deadline", "1h30m"}
	Parse()

	got := c.Parse()
	if got.Deadline == nil {
		t.Fatal("Deadline = nil, want non-nil")
	}
	if *got.Deadline != 90*time.Minute {
		t.Errorf("*Deadline = %v, want 1h30m", *got.Deadline)
	}
}

func TestAddFlagsForStruct_DurationPointerLeaf_PresentViaEnv(t *testing.T) {
	resetForStructFlagTest(t)
	setEnv(t, "APP_DEADLINE", "10s")

	cfg := &durationConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd"}
	Parse()

	got := c.Parse()
	if got.Deadline == nil {
		t.Fatal("Deadline = nil, want non-nil")
	}
	if *got.Deadline != 10*time.Second {
		t.Errorf("*Deadline = %v, want 10s", *got.Deadline)
	}
}

func TestAddFlagsForStruct_DurationPointerLeaf_AbsentStaysNil(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &durationConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-timeout", "1s"}
	Parse()

	got := c.Parse()
	if got.Deadline != nil {
		t.Errorf("Deadline = %v, want nil", *got.Deadline)
	}
}

// A duration leaf nested inside a value-typed branch and a pointer-typed
// branch, mirroring the pointer-scalar nesting tests.
type durationInValueBranch struct {
	Interval *time.Duration `flag:"interval"`
}
type durationInPointerBranch struct {
	Interval *time.Duration `flag:"interval"`
}
type nestedDurationConfig struct {
	ValueBranch   durationInValueBranch    `flag:"value"`
	PointerBranch *durationInPointerBranch `flag:"ptr"`
}

func TestAddFlagsForStruct_DurationLeaf_NestedInValueBranch(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &nestedDurationConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-value-interval", "3s"}
	Parse()

	got := c.Parse()
	if got.ValueBranch.Interval == nil || *got.ValueBranch.Interval != 3*time.Second {
		t.Errorf("ValueBranch.Interval = %+v, want *3s", got.ValueBranch.Interval)
	}
	if got.PointerBranch != nil {
		t.Errorf("PointerBranch = %+v, want nil (its own leaf never provided)", *got.PointerBranch)
	}
}

func TestAddFlagsForStruct_DurationLeaf_NestedInPointerBranch(t *testing.T) {
	resetForStructFlagTest(t)

	cfg := &nestedDurationConfig{}
	c, err := AddFlagsForStruct("app", cfg)
	if err != nil {
		t.Fatalf("AddFlagsForStruct: %v", err)
	}

	os.Args = []string{"cmd", "-app-ptr-interval", "7s"}
	Parse()

	got := c.Parse()
	if got.ValueBranch.Interval != nil {
		t.Errorf("ValueBranch.Interval = %v, want nil", *got.ValueBranch.Interval)
	}
	if got.PointerBranch == nil {
		t.Fatal("PointerBranch = nil, want non-nil")
	}
	if got.PointerBranch.Interval == nil || *got.PointerBranch.Interval != 7*time.Second {
		t.Errorf("PointerBranch.Interval = %+v, want *7s", got.PointerBranch.Interval)
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
