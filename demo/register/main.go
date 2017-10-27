package main

import (
	"net"
	"fmt"
)

var strChan = make(chan string)

type Connect struct {
	conn net.Conn
}

func (theConn *Connect) send(message string) {
	theConn.conn.Write([]byte(message))
}

type Register struct {
	conns []Connect
	message chan string
}

func (register *Register) addConn(conn net.Conn){
	register.conns = append(register.conns, Connect{conn: conn})
}

func (register *Register) bootstrap() {
	for {
		message := <- register.message
		for _, conn := range register.conns {
			conn.send(message)
		}
	}
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
	register := &Register{message: make(chan string), conns: []Connect{}}

	go register.bootstrap()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // 终止程序
		}
		// connList = append(connList, conn)
		register.addConn(conn)
		go doServerStuff(register, conn)
	}
}

func doServerStuff(register *Register, conn net.Conn) {
	for {
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return //终止程序
		}
		fmt.Printf("Received data: %v", string(buf[:len]))

		(*register).message <- string(buf[:len])
		// strChan <- string(buf[:len])
		// _, err = conn.Write(buf[:len])
	}
}
