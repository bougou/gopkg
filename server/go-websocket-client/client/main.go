package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"
)

func main() {
	var serverAddr string
	flag.StringVar(&serverAddr, "addr", "127.0.0.1:8040", "the server addr")
	flag.Parse()

	// stopCh := make(chan struct{})

	reloadCh := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-signals

		fmt.Println("caught signal")
		for {
			switch sig {
			case os.Interrupt, syscall.SIGTERM, syscall.SIGINT:
				fmt.Println("caught Interrupt signal, stopping")
				cancel()

				return
			case syscall.SIGHUP:
				fmt.Println("caught SIGHUP signal, reloading")
				reloadCh <- struct{}{}
			default:
				fmt.Println("not recognized signal")
			}
		}

	}()

	scheme := "ws"
	path := "/api/v1/agents/websocket"

	client := NewClient(scheme, serverAddr, path)
	client.Run(ctx)
}
