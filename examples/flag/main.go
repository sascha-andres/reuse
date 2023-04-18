package main

import "github.com/sascha-andres/reuse/flag"
import "fmt"

func main() {
	var boolFlag bool
	var stringFlag string

	flag.BoolVar(&boolFlag, "bool", false, "a boolean flag")
	flag.StringVar(&stringFlag, "string", "default", "a string flag")

	flag.Parse()

	verbs := flag.GetVerbs()

	for _, verb := range verbs {
		fmt.Printf(" found verb %q\n", verb)
	}
}
