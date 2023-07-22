package server

import (
	"fmt"
	"net"
)

func server() {
	tcpConn("127.0.0.1", 22123)
	udpConn("127.0.0.1", 23456)
}

func tcpConn(ip string, port int32) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer conn.Close()

	fmt.Println(conn.LocalAddr().String())
	fmt.Println(conn.RemoteAddr().String())

	return nil
}

func udpConn(ip string, port int32) error {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer conn.Close()

	fmt.Println(conn.LocalAddr().String())
	fmt.Println(conn.RemoteAddr().String())

	if _, err := conn.Write([]byte("hello")); err != nil {
		fmt.Println("could not write payload to server:", err)
	}

	return nil
}
