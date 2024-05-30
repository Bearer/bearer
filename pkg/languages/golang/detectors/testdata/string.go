package main

import (
	"fmt"
	"os"
)

var Greeting = "Hello World"

func main() {
	s := Greeting + "!"
	s += "!!"
	fmt.Println(s)

	s2 := "hey "
	s2 += os.Args[0]
	s2 += " there"
	fmt.Println(s2)

	s3 := "foo " + os.Args[0] + " bar"
	fmt.Println(s3)
}
