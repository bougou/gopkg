package websocket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

func NewWebSocketServerHandler(genMsgFunc GenMsgFunc, handleMsgFunc HandleMsgFunc) websocket.Handler {
	s := NewWebSocketServer(genMsgFunc, handleMsgFunc)

	return websocket.Handler(func(conn *websocket.Conn) {
		s.Run()
	})
}

type WebSocketServer struct {
	genMsgFunc    GenMsgFunc
	handleMsgFunc HandleMsgFunc

	conn   *websocket.Conn
	sendCh chan []byte
	recvCh chan []byte

	peerClosedCh chan int
	retryCh      chan int

	retryIntervalSec int
}

type GenMsgFunc func(cxt context.Context, sendCh chan<- []byte) error
type HandleMsgFunc func(msg []byte, sentCh chan<- []byte) error

func NewWebSocketServer(genMsgFunc GenMsgFunc, handleMsgFunc HandleMsgFunc) *WebSocketServer {
	return &WebSocketServer{

		genMsgFunc:    genMsgFunc,
		handleMsgFunc: handleMsgFunc,

		sendCh: make(chan []byte),
		recvCh: make(chan []byte),

		peerClosedCh:     make(chan int, 1),
		retryCh:          make(chan int),
		retryIntervalSec: 10,
	}
}

func (c *WebSocketServer) Close() error {
	return c.conn.Close()
}

func (c *WebSocketServer) Run() error {

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
	}()

	go func() {
		defer wg.Done()
	}()

	go func() {
		defer wg.Done()
	}()

	go func() {
		defer wg.Done()
	}()

	wg.Wait()
	fmt.Println("wait finish, run exit")

	close(c.recvCh)
	close(c.sendCh)
	close(c.peerClosedCh)

	return nil
}

func (c *WebSocketServer) gen(ctx context.Context) error {
	if c.genMsgFunc != nil {
		// block to gen msg without exit until ctx is done
		c.genMsgFunc(ctx, c.sendCh)
		fmt.Println("gen exit")
		return nil
	}
	return nil
}

func (c *WebSocketServer) handle(ctx context.Context) error {
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

func (c *WebSocketServer) sent(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("==> ctx done, send exit")
			return nil

		case msg := <-c.sendCh:
			if err := websocket.Message.Send(c.conn, msg); err != nil {
				fmt.Printf("send failed, err: %s\n", err)
			} else {
				now := time.Now().Format("2006-01-02 15:04:05")
				fmt.Printf("send succeed, %s\n", now)
			}
		}
	}
}

func (c *WebSocketServer) recv(ctx context.Context) error {
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
