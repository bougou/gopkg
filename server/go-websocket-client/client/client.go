package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/gopkg/server/websocket"
	"github.com/kr/pretty"
)

type Client struct {
	scheme string
	addr   string
	path   string

	client *websocket.WebSocketClient
}

func NewClient(scheme, addr string, path string) *Client {

	c := websocket.NewWebSocketClient(scheme, addr, path, generateMsg, handleMsg)

	return &Client{
		scheme: scheme,
		addr:   addr,
		path:   path,

		client: c,
	}
}

func (c *Client) Run(ctx context.Context) error {
	return c.client.Run(ctx)
}

func (c *Client) Stop() error {
	return c.client.Close()

}

type Heartbeat struct {
	Time time.Time
}

type MsgType string

const (
	MsgTypeAck       MsgType = "ack"
	MsgTypeHeartbeat MsgType = "heartbeat"
)

type ClientMsg struct {
	MsgType   MsgType
	ACK       *string
	Heartbeat *Heartbeat
}

func NewClientAckMsg() *ClientMsg {
	now := time.Now().Format("2006-01-02 15:04:05")

	return &ClientMsg{
		MsgType: MsgTypeAck,
		ACK:     common.StringPtr(fmt.Sprintf("%s: client receive succeed", now)),
	}
}

func NewHeartbeatMsg() *ClientMsg {
	now := time.Now().Format("2006-01-02 15:04:05")

	return &ClientMsg{
		MsgType: MsgTypeHeartbeat,
		ACK:     common.StringPtr(fmt.Sprintf("%s: client heartbeat msg", now)),
	}
}

type ServerMsgType string

const (
	ServerMsgTypeConfChange ServerMsgType = "conf-change"
)

type ServerMsg struct {
	MsgType    ServerMsgType
	ConfChange []byte
}

func handleMsg(msg []byte, sentCh chan<- []byte) error {
	fmt.Printf("client handle msg byte: %s\n", msg)
	serverMsg := &ServerMsg{}

	if err := json.Unmarshal(msg, serverMsg); err != nil {
		return fmt.Errorf("decode msg failed, err: %s", err)
	}
	now := time.Now().Format("2006-01-02 15:04:05")

	pretty.Printf("%s: receive succeed, server msg: %v\n", now, serverMsg)

	ackMsg := NewClientAckMsg()

	ackData, err := json.Marshal(ackMsg)
	if err != nil {
		return fmt.Errorf("encode ack msg failed, err: %s", err)
	}

	sentCh <- ackData
	return nil
}

func generateMsg(ctx context.Context, sentCh chan<- []byte) error {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		generatePeriodHearbeat(ctx, sentCh)
		fmt.Println("gen msg1 exit")
	}()

	go func() {
		defer wg.Done()
		generatePeriodHearbeat2(ctx, sentCh)
		fmt.Println("gen msg2 exit")

	}()

	wg.Wait()

	fmt.Println("wait finish, gen msg exit")

	return nil
}

func generatePeriodHearbeat(ctx context.Context, sentCh chan<- []byte) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			heartbeatMsg := NewHeartbeatMsg()
			data, err := json.Marshal(heartbeatMsg)
			if err != nil {
				return fmt.Errorf("encode heartbeat msg failed, err: %s", err)
			}
			now := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("client begin send hearbeat at %s\n", now)

			sentCh <- data
			fmt.Printf("client send hearbeat at %s\n", now)
		}
	}
}

func generatePeriodHearbeat2(ctx context.Context, sentCh chan<- []byte) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			heartbeatMsg := NewHeartbeatMsg()
			data, err := json.Marshal(heartbeatMsg)
			if err != nil {
				return fmt.Errorf("encode heartbeat msg 2 failed, err: %s", err)
			}
			now := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("client begin send hearbeat 2 at %s\n", now)

			sentCh <- data
			fmt.Printf("client send hearbeat 2 at %s\n", now)
		}
	}

}
