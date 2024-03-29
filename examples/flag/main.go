package main

import (
	"fmt"

	"github.com/sascha-andres/reuse/flag"
)

func main() {
	var boolFlag bool
	var stringFlag string

	flag.BoolVar(&boolFlag, "bool", false, "a boolean flag")
	flag.StringVar(&stringFlag, "string", "default", "a string flag")
	flag.SetSeparated()

	flag.Parse()

	verbs := flag.GetVerbs()

	fmt.Println("Verbs:")
	for _, verb := range verbs {
		fmt.Printf(" found verb %q\n", verb)
	}
	fmt.Printf("separated: %s\n", flag.GetSeparated())
	fmt.Printf("bool: %t\n", boolFlag)
	fmt.Printf("string: %s\n", stringFlag)
}
