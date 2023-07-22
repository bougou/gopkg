package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	var serverAddr string
	flag.StringVar(&serverAddr, "addr", "127.0.0.1:8040", "the server addr")
	flag.Parse()

	url := fmt.Sprintf("ws://%s/api/v1/agents/websocket", serverAddr)
	conn, err := websocket.Dial(url, "", "http://www.baidu.com")
	if err != nil {
		fmt.Printf("websocket dail failed, err: %s\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("websocket dail succeed.")

	person := &Person{
		Name: "tom",
		Age:  20,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// recv server sent messages and handle it
	go func() {
		defer wg.Done()
		for {
			if err := websocket.JSON.Receive(conn, person); err != nil {
				fmt.Printf("receive failed, err: %s\n", err)
			} else {
				now := time.Now().Format("2006-01-02 15:04:05")
				fmt.Printf("%s: receive succeed, %v\n", now, *person)
			}
		}
	}()

	// send messages to server
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for ; true; <-ticker.C {
			if err := websocket.JSON.Send(conn, person); err != nil {
				fmt.Printf("send failed, err: %s\n", err)
			} else {
				now := time.Now().Format("2006-01-02 15:04:05")
				fmt.Printf("%s: send succeed, %v\n", now, *person)
			}

		}
	}()

	wg.Wait()

	fmt.Println("exit")
}
