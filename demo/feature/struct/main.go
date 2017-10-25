package main

import (
	"fmt"
)

type MyCls struct {
	channel chan *int;
	num int;
}

func (cls* MyCls) connect () {
	message := <- cls.channel
	fmt.Printf("welcome, %d\n", *message)
}

func main () {
	defer func(){
		println("hold on")
	}()

	channel := make(chan *int)
	mycls := &MyCls{channel, 5}

	go (func (){
		defer func (){
			println("inner defer")
		}()

		num := 6
		println("outer defer")
		mycls.channel <- &num
	})()

	mycls.connect()
	println("Finish")
}
