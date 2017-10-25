package main

import (
	"strings"
	"fmt"
)

func main() {
	rawStr := "name,title,job,salary"
	strList := strings.Split(rawStr, ",")
	fmt.Println(strings.Join(strList, "\t"))
}
