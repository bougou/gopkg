package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("mutual wait")

	if err := mutualWait(); err != nil {
		fmt.Printf("mutualWait failed, err: %s\n", err)
	}
}

func mutualWait() error {

	var wg sync.WaitGroup
	wg.Add(2)

	exitCh := make(chan error)

	go func(_wg *sync.WaitGroup, _exitCh chan error) {
		defer _wg.Done()

		var i int
		for {
			select {
			case err := <-_exitCh:
				// doHeavyTask exit
				fmt.Printf("goroutine-1 recv exit of doHeavyTask, exitErr: %s, goroutine-1 exit\n", err)
				return

			case <-time.After(1 * time.Second):
				// goroutine-1 do its own task untils recv exit from parent doHeavyTask
				i++
				fmt.Printf("goroutine-1 run #%d\n", i)
			}
		}

	}(&wg, exitCh)

	go func(_wg *sync.WaitGroup, _exitCh chan error) {
		defer wg.Done()

		var i int
		for {
			select {
			case err := <-_exitCh:
				// doHeavyTask exit
				fmt.Printf("goroutine-2 recv exit of doHeavyTask, exitErr: %s, goroutine-2 exit\n", err)
				return

			case <-time.After(1 * time.Second):
				// goroutine-2 do its own task at most 5 seconds.
				i++
				fmt.Printf("goroutine-2 run #%d\n", i)
				if i > 5 {
					fmt.Printf("goroutine-2 exit after #%d\n", i)
					return
				}
			}

		}

	}(&wg, exitCh)

	fmt.Println("doHeavyTask start")
	err := doHeavyTask()
	fmt.Println("doHeavyTask end")
	if err != nil {
		fmt.Printf("doHeavyTask failed, err: %s\n", err)
	}
	exitCh <- err

	wg.Wait()
	fmt.Printf("all sub-goroutines exit")

	return err
}

func doHeavyTask() error {
	for i := 1; i < 10; i++ {
		time.Sleep(2 * time.Second)
	}

	return nil
}
