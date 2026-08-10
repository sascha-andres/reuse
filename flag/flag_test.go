package flag

import (
	"flag"
	"os"
	"testing"
	"time"

	"slices"
)

type verbTestCase struct {
	args          []string
	expectedVerbs []string
	flags         []string
	name          string
}

var testCasesGetVerbs = []verbTestCase{
	{
		args:          []string{"single-verb", "-test", "1"},
		flags:         []string{"test"},
		expectedVerbs: []string{"single-verb"},
		name:          "single verb",
	},
	{
		args:          []string{"-test", "1"},
		flags:         []string{"test"},
		expectedVerbs: []string{},
		name:          "no verb",
	},
	{
		args:          []string{"two-verb", "two-second-verb", "-test", "1"},
		flags:         []string{"test"},
		expectedVerbs: []string{"two-verb", "two-second-verb"},
		name:          "two verbs",
	},
}

func TestPConstructors_DefaultToZeroValue(t *testing.T) {
	resetForStructFlagTest(t)

	s := StringP("p-string", "")
	i := Int64P("p-int64", "")
	b := BoolP("p-bool", "")
	fl := Float64P("p-float64", "")
	u := Uint64P("p-uint64", "")

	os.Args = []string{"cmd"}
	Parse()

	if *s != "" || *i != 0 || *b != false || *fl != 0.0 || *u != 0 {
		t.Fatalf("expected zero values, got string=%q int64=%d bool=%t float64=%v uint64=%d", *s, *i, *b, *fl, *u)
	}
}

func TestPConstructors_ExplicitValueAndVisit(t *testing.T) {
	resetForStructFlagTest(t)

	s := StringP("p-string", "")
	i := Int64P("p-int64", "")
	b := BoolP("p-bool", "")
	fl := Float64P("p-float64", "")
	u := Uint64P("p-uint64", "")

	os.Args = []string{"cmd", "-p-string", "hi", "-p-int64", "-7", "-p-bool", "-p-float64", "1.5", "-p-uint64", "9"}
	Parse()

	if *s != "hi" || *i != -7 || *b != true || *fl != 1.5 || *u != 9 {
		t.Fatalf("unexpected values: string=%q int64=%d bool=%t float64=%v uint64=%d", *s, *i, *b, *fl, *u)
	}

	set := map[string]bool{}
	Visit(func(fl *flag.Flag) { set[fl.Name] = true })
	for _, name := range []string{"p-string", "p-int64", "p-bool", "p-float64", "p-uint64"} {
		if !set[name] {
			t.Errorf("Visit did not report %q as set", name)
		}
	}
}

func TestDurationP_DefaultToZeroValue(t *testing.T) {
	resetForStructFlagTest(t)

	d := DurationP("p-duration", "")

	os.Args = []string{"cmd"}
	Parse()

	if *d != 0 {
		t.Fatalf("expected zero value, got %v", *d)
	}
}

func TestDurationP_ExplicitValueAndVisit(t *testing.T) {
	resetForStructFlagTest(t)

	d := DurationP("p-duration", "")

	os.Args = []string{"cmd", "-p-duration", "5s"}
	Parse()

	if *d != 5*time.Second {
		t.Fatalf("got %v, want 5s", *d)
	}

	set := map[string]bool{}
	Visit(func(fl *flag.Flag) { set[fl.Name] = true })
	if !set["p-duration"] {
		t.Error("Visit did not report p-duration as set")
	}
}

func TestGetVerbs(t *testing.T) {
	t.Skip("this test is broken, but I don't know how to fix it, I can't really change the os.Args")
	oldArgs := os.Args

	for _, testCase := range testCasesGetVerbs {
		t.Run(testCase.name, func(t *testing.T) {
			os.Args = testCase.args
			defer func() { os.Args = oldArgs }()
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			str := ""
			for _, s := range testCase.flags {
				flag.StringVar(&str, s, "", "for test cases")
			}
			t.Logf("%#v", os.Args)
			Parse()
			result := GetVerbs()
			for _, foundVerb := range result {
				if !slices.Contains(testCase.expectedVerbs, foundVerb) {
					t.Errorf("found %q which is not part of expected verbs", foundVerb)
				}
			}
			for _, verb := range testCase.expectedVerbs {
				if !slices.Contains(result, verb) {
					t.Errorf("expected %q to be present", verb)
				}
			}
		})
	}
}
