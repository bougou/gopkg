package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bougou/gopkg/server/websocket"
)

func main() {

	mux := http.NewServeMux()
	mux.Handle("/", websocket.NewWebSocketServerHandler(generateMsg, handleMsg))

	// mux.Handle("/", websocket.Handler(websocketServerHandle))

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "127.0.0.1", 8040),
		Handler: mux,
	}
	fmt.Printf("start listening, %s\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		fmt.Printf("listen failed")
	}
}

func generateMsg(ctx context.Context, sentCh chan<- []byte) error {
	ticker := time.NewTicker(8 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			fmt.Println("generate server msg")

		}
	}

}

func handleMsg(msg []byte, sentCh chan<- []byte) error {
	fmt.Printf("server handle msg byte: %s\n", msg)

	fmt.Println("server handle msg succeed")

	ackData := []byte("server ack msg")
	sentCh <- ackData
	return nil
}
