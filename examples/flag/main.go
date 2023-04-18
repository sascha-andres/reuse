package main

import "github.com/sascha-andres/reuse/flag"
import "fmt"

func main() {
  flag.Parse()

  verbs := flag.GetVerbs()

  for _, verb := range verbs {
    fmt.Printf(" found verb %q\n", verb)
  }
}
