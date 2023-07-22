package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/", websocket.Handler(handleFunc))

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "127.0.0.1", 8040),
		Handler: mux,
	}
	fmt.Printf("start listening, %s\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		fmt.Printf("listen failed")
	}
}

func handleFunc(conn *websocket.Conn) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	go func() {
		ticker := time.NewTicker(8 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			heartbeatMsg := "server heartbeat 1"
			data, err := json.Marshal(heartbeatMsg)
			if err != nil {
				fmt.Printf("encode heartbeat msg failed, err: %s\n", err)
				continue
			}
			now := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("server begin send hearbeat at %s\n", now)
			if err := websocket.Message.Send(conn, []byte(data)); err != nil {
				fmt.Printf("websocket send failed, err: %s\n", err)
			}
			fmt.Printf("server send hearbeat at %s\n", now)
		}
	}()

	for {
		msg := make([]byte, 1024)
		if err := websocket.Message.Receive(conn, &msg); err != nil {
			fmt.Printf("websocket receive failed, err: %s\n", err)
			return
		}
		fmt.Printf("websocket receive succeed, %v\n", msg)

		// processing
		// person.Age++
		sentMsg := "server receive succeed"
		if err := websocket.Message.Send(conn, []byte(sentMsg)); err != nil {
			fmt.Printf("websocket send failed, err: %s\n", err)
		}
		fmt.Printf("websocket send succeed, %v\n", sentMsg)
	}

}
