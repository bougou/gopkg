package websocket

import (
	"context"
	"fmt"
	"path"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type WebSocketClient struct {
	scheme        string
	addr          string
	path          string
	genMsgFunc    GenMsgFunc
	handleMsgFunc HandleMsgFunc

	conn   *websocket.Conn
	sendCh chan []byte
	recvCh chan []byte

	peerClosedCh chan int
	retryCh      chan int

	retryIntervalSec int
}

func NewWebSocketClient(scheme, addr, path string, genMsgFunc GenMsgFunc, handleMsgFunc HandleMsgFunc) *WebSocketClient {
	return &WebSocketClient{
		scheme:        scheme,
		addr:          addr,
		path:          path,
		genMsgFunc:    genMsgFunc,
		handleMsgFunc: handleMsgFunc,

		sendCh: make(chan []byte),
		recvCh: make(chan []byte),

		peerClosedCh:     make(chan int, 1),
		retryCh:          make(chan int),
		retryIntervalSec: 10,
	}
}

func (c *WebSocketClient) Close() error {
	return c.conn.Close()
}

func (c *WebSocketClient) Run(ctx context.Context) error {
	if c.conn == nil {
		fmt.Println("conn not exists, connecting...")
		if err := c._connect(ctx); err != nil {
			return fmt.Errorf("connect failed, err: %s", err)
		}
		fmt.Println("connected, :)")
	}

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		c.connect(ctx)
	}()

	go func() {
		defer wg.Done()
		c.gen(ctx)
	}()

	go func() {
		defer wg.Done()
		c.handle(ctx)
	}()

	go func() {
		defer wg.Done()
		c.send(ctx)
	}()

	go func() {
		defer wg.Done()
		c.recv(ctx)
	}()

	wg.Wait()
	fmt.Println("wait finish, run exit")

	close(c.recvCh)
	close(c.sendCh)
	close(c.peerClosedCh)

	return nil
}

func (c *WebSocketClient) gen(ctx context.Context) error {
	if c.genMsgFunc != nil {
		// block to gen msg without exit until ctx is done
		c.genMsgFunc(ctx, c.sendCh)
		fmt.Println("gen exit")
		return nil
	}
	return nil
}

func (c *WebSocketClient) handle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			c.conn.Close()

			fmt.Println("==> ctx done, handle exit")
			return nil

		case msg := <-c.recvCh:
			fmt.Println("handle got msg from recv channle, handle it")
			if err := c.handleMsgFunc(msg, c.sendCh); err != nil {
				fmt.Printf("hanle msg failed: %s\n", err)
			} else {
				fmt.Println("hanle msg succeed")
			}
		}
	}
}

// fetch msg from sendCh and sent it.
func (c *WebSocketClient) send(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("==> ctx done, send exit")
			return nil

		case msg := <-c.sendCh:
			if err := websocket.Message.Send(c.conn, msg); err != nil {
				fmt.Printf("send failed, err: %s\n", err)
			}
		}
	}

}

// Recv Msg and check whether the peer closed
func (c *WebSocketClient) recv(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("recv exit")
			return nil

		default:

			recvMsg := make([]byte, 1024)
			// second parameter must be pointer
			fmt.Println("block on receive")
			if err := websocket.Message.Receive(c.conn, &recvMsg); err != nil {
				if err == websocket.ErrFrameTooLarge {
					fmt.Printf("warn: %s\n", err)
					continue
				}

				fmt.Printf("receive failed, err: %s\n", err)
				// the connect goroutine should read this and rebuild the conn
				c.peerClosedCh <- 1
				continue
			} else {
				fmt.Printf("receive succed, send to recv channel, %s\n", recvMsg)
				c.recvCh <- recvMsg
				fmt.Println("receive succed, send to recv channel finished")
			}

		}

	}
}

func (c *WebSocketClient) _c(ctx context.Context) error {
	url := path.Join(c.addr, c.path)
	url = fmt.Sprintf("%s://%s", c.scheme, url)
	fmt.Printf("connect to %s\n", url)
	conn, err := websocket.Dial(url, "", "http://www.baidu.com")
	if err != nil {
		fmt.Printf("connect server (%s) failed, err: %s\n", url, err)
		fmt.Printf("retry after %d seconds\n", c.retryIntervalSec)

		return fmt.Errorf("connect failed, err: %s", err)
	}

	fmt.Println("connect server succeed.")
	fmt.Printf("local addr: %s\n", conn.LocalAddr())
	fmt.Printf("remote addr: %s\n", conn.RemoteAddr())
	c.conn = conn
	return nil
}

func (c *WebSocketClient) _connect(ctx context.Context) error {
	var retryCh = make(chan struct{})
	go func() {
		retryCh <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("ctx done, stop connect")
			return nil

		case <-retryCh:
			if err := c._c(ctx); err != nil {
				go func() {
					time.Sleep(time.Duration(c.retryIntervalSec) * time.Second)
					retryCh <- struct{}{}
				}()
				continue
			}
			return nil
		}
	}
}

// Connect and re-connect
func (c *WebSocketClient) connect(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("==> ctx done, connect exit")
			return nil

		case <-c.peerClosedCh:
			fmt.Println("error, peer closed, re-connect...")
			c._connect(ctx)
		}
	}
}
