package main

import (
	"fmt"

	"github.com/sascha-andres/reuse/flag"
)

type Config struct {
	Name      string  `flag:"name,set the name"`
	Age       int     `flag:"age,set the age"`
	Human     bool    `flag:"human,set if human"`
	Weight    float64 `flag:"weight,set weight"`
	Height    uint    `flag:"height,set height"`
	SubStruct struct {
		Value uint `flag:"val"`
	} `flag:"sub"`
}

func main() {
	flag.SetEnvPrefix("STRUCTFLAG")

	a := &Config{Age: 18}
	c, err := flag.AddFlagsForStruct("test", a)
	if err != nil {
		panic(err)
	}
	b := &Config{}
	_, err = flag.AddFlagsForStruct("test-2", b)
	if err != nil {
		panic(err)
	}
	flag.Parse()
	flag.PrintDefaults()
	c.Parse()
	fmt.Printf("struct: %#v\n", a)
	fmt.Printf("height: %d\n", a.Height)
}
