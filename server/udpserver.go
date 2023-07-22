package server

import (
	"fmt"
	"net"
)

type UDPServer struct {
	ListenAddr string
	ListenPort int
}

func NewUDPServer(addr string, port int) *UDPServer {
	return &UDPServer{
		ListenAddr: addr,
		ListenPort: port,
	}
}

func (s *UDPServer) Start() error {
	addr := &net.UDPAddr{
		Port: s.ListenPort,
		IP:   net.ParseIP(s.ListenAddr),
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("listen udp failed, err: %s", err)
	}

	maxBufferSize := 10
	buf := make([]byte, maxBufferSize)

	for {
		_, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("ERR: read from %s failed, err: %s", addr.String(), err)
			continue
		}
		go response(conn, remoteAddr)
	}
}

func response(conn *net.UDPConn, addr *net.UDPAddr) {
	conn.WriteToUDP([]byte("server recved"), addr)
}
