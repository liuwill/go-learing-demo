package main

import (
	"flag"
)

// go run main.go --addr :8000
var addr = flag.String("addr", ":8080", "http service address")

func testDefer(result int) int {
	defer func() {
		result++
	}()
	return 0
}

func main() {
	flag.Parse()
	println(*addr)

	println(testDefer(3))
}
