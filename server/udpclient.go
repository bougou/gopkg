package server

import (
	"fmt"
	"net"
)

type UDPClient struct {
	Addr string
	Port int

	conn *net.UDPConn
}

func NewUDPClient(addr string, port int) *UDPClient {
	return &UDPClient{
		Addr: addr,
		Port: port,
	}
}

func (c *UDPClient) Start() error {
	addr := &net.UDPAddr{
		IP:   net.ParseIP("10.9.80.194"),
		Port: 623,
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("connect failed, err: %s", err)
	}
	c.conn = conn
	return nil
}

func (c *UDPClient) SendThenRecv(data []byte) error {
	if c.conn == nil {
		return fmt.Errorf("conn not started")
	}

	fmt.Fprint(c.conn, string(data))

	buf := make([]byte, 2048)
	_, remoteAddr, err := c.conn.ReadFromUDP(buf)
	if err != nil {
		return fmt.Errorf("read from server failed, err: %s", err)
	}
	fmt.Printf("read: %s from %s", buf, remoteAddr.String())

	return nil
}

func (c *UDPClient) Close() error {
	defer c.conn.Close()
	return nil
}
