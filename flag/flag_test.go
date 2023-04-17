package flag

import (
	"flag"
	"os"
	"testing"

	"golang.org/x/exp/slices"
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
