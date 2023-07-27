package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"golang.org/x/net/proxy"
)

type UDPClient struct {
	Host string // Host can be literal IP address or DNS name.
	Port int

	proxy proxy.Dialer
	conn  net.Conn

	bufferSize int
	timeout    time.Duration
}

func NewUDPClient(host string, port int) *UDPClient {
	return &UDPClient{
		Host: host,
		Port: port,
	}
}

func (c *UDPClient) WithProxyDailer(p proxy.Dialer) *UDPClient {
	c.proxy = p
	return c
}

func (c *UDPClient) InitConn() error {
	if c.conn != nil {
		return nil
	}

	if c.proxy != nil {
		conn, err := c.proxy.Dial("udp", fmt.Sprintf("%s:%d", c.Host, c.Port))
		if err != nil {
			return fmt.Errorf("proxy dail failed, err: %s", err)
		}
		c.conn = conn
		return nil
	}

	remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		return fmt.Errorf("resolve UDP addr failed, err: %s", err)
	}

	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		return fmt.Errorf("dial UDP failed, err: %s", err)
	}
	c.conn = conn

	return nil
}

func (c *UDPClient) ReInitConn() error {
	c.conn = nil
	return c.InitConn()
}

// If you have some bytes to send, use `bytes.NewReader(b []byte)` to convert
// the byte slice to bytes.Reader struct which implemented the io.Reader interface.
func (c *UDPClient) Exchange(reader io.Reader) ([]byte, error) {
	ctx := context.Background()
	return c.ExchangeContext(ctx, reader)
}

// If you have some bytes to send, use `bytes.NewReader(b []byte)` to convert
// the byte slice to bytes.Reader struct which implemented the io.Reader interface.
func (c *UDPClient) ExchangeContext(ctx context.Context, reader io.Reader) ([]byte, error) {
	if err := c.InitConn(); err != nil {
		return nil, fmt.Errorf("init udp conn failed, err: %w", err)
	}
	recvBuffer := make([]byte, c.bufferSize)

	doneChan := make(chan error, 1)
	// recvChan stores the integer number which indicates how many bytes
	recvChan := make(chan int, 1)

	go func() {
		// Send data to remote through UDP conn
		_, err := io.Copy(c.conn, reader)
		if err != nil {
			doneChan <- fmt.Errorf("write to conn failed, err: %w", err)
			return
		}

		// Set a deadline for the ReadOperation so that we don't
		// wait forever for a server that might not respond on
		// a resonable amount of time.
		deadline := time.Now().Add(c.timeout)
		err = c.conn.SetReadDeadline(deadline)
		if err != nil {
			doneChan <- fmt.Errorf("set conn read deadline failed, err: %w", err)
			return
		}

		// Recv data from the UDP conn
		nRead, err := c.conn.Read(recvBuffer)
		if err != nil {
			doneChan <- fmt.Errorf("read from conn failed, err: %w", err)
			return
		}

		doneChan <- nil
		recvChan <- nRead
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("canceled from caller")

	case err := <-doneChan:
		if err != nil {
			return nil, err
		}
		recvCount := <-recvChan
		return recvBuffer[:recvCount], nil
	}
}

func (c *UDPClient) Start() error {
	return c.InitConn()
}

func (c *UDPClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// RemoteIP returns the parsed ip address of the target.
func (c *UDPClient) RemoteIPPort() (string, int) {
	if net.ParseIP(c.Host) == nil {
		addrs, err := net.LookupHost(c.Host)
		if err == nil && len(addrs) > 0 {
			return addrs[0], c.Port
		}
	}
	return c.Host, c.Port
}

func (c *UDPClient) LocalIPPort() (string, int) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		return "", 0
	}
	defer conn.Close()
	host, port, _ := net.SplitHostPort(conn.LocalAddr().String())
	p, _ := strconv.Atoi(port)
	return host, p
}
