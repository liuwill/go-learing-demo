package main

import (
	"net"
	"fmt"
)

var strChan = make(chan string)

type Connect struct {
	conn net.Conn
	message chan string
}

func main() {
	fmt.Println("Starting the server ...")
	// 创建 listener
	listener, err := net.Listen("tcp", "localhost:50000")
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return //终止程序
	}
	// 监听并接受来自客户端的连接
	connList := []net.Conn{}

	go func (){
		for {
			message := <- strChan

			println(message)
			for _, conn := range connList {
				// println(*conn.Write)

				conn.Write([]byte(message))
			}
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // 终止程序
		}
		connList = append(connList, conn)
		go doServerStuff(conn)
	}
}

func doServerStuff(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return //终止程序
		}
		fmt.Printf("Received data: %v", string(buf[:len]))

		strChan <- string(buf[:len])
		// _, err = conn.Write(buf[:len])
	}
}
