package snippet

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

// do things regularly until ctx done, then exit and notify waitgroup.
func tick(id int, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Println("hi")
		case <-ctx.Done():
			fmt.Println("exit", id)
			return
		}
	}
}

func main() {
	// var stopCh = make(chan struct{})
	// We use ctx to substitute stopCh
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(3)

	go tick(1, &wg, ctx)
	go tick(2, &wg, ctx)
	go tick(3, &wg, ctx)

	// listen for C-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("main: received C-c - shutting down")

	// close(stopCh)
	// cancel the ctx instead of close the stop channel
	cancel()

	wg.Wait()
	fmt.Println("main exit")

}
