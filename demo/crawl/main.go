package main

import (
	"fmt"
	"flag"
)

var command = flag.String("command", "simple", "command type")

func main() {
	flag.Parse()

	if *command == "simple" {
		parseSimple()
		<-(chan int)(nil)
	} else if *command == "reader" {
		parseReader()
		<-(chan int)(nil)
	} else if *command == "writer" {
		parseWriter()
		<-(chan int)(nil)
	}

	fmt.Println("Done!")
}
