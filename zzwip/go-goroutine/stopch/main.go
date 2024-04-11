package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func tick(id int, wg *sync.WaitGroup, stopCh chan struct{}) {
	defer wg.Done()
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Println("hi")
		case <-stopCh:
			fmt.Println("exit", id)
			return
		}
	}
}

func main() {
	var stopCh = make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(3)

	go tick(1, &wg, stopCh)
	go tick(2, &wg, stopCh)
	go tick(3, &wg, stopCh)

	// listen for C-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("main: received C-c - shutting down")

	close(stopCh)
	// Note, CAN NOT use the following method
	// which just send a value to the channel, thus only one goroutine received the value.
	// stopCh <- struct{}{}

	wg.Wait()
	fmt.Println("main exit")

}
