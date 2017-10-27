package main

import (
	"net"
	"fmt"
)

var strChan = make(chan string)

type Connect struct {
	conn net.Conn
	position int
	register *Register
}

func (theConn *Connect) send(message string) {
	theConn.conn.Write([]byte(message))
}

type Register struct {
	conns []Connect
	message chan string
}

func (register *Register) addConn(conn net.Conn) Connect{
	connect := Connect{conn: conn, position: len(register.conns)}
	register.conns = append(register.conns, connect)
	connect.register = register

	return connect
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

		connect := register.addConn(conn)
		go doServerStuff(register, connect)
	}
}

func doServerStuff(register *Register, connect Connect) {
	conn := connect.conn
	for {
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return //终止程序
		}

		message := string(buf[:len])
		fmt.Printf("Received data: %v", message)

		(*register).message <- message
	}
}
