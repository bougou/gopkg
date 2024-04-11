package snippet

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	Channel()
}

func Channel1() {
	var wg sync.WaitGroup
	wg.Add(2)

	var queue = make(chan bool)

	fmt.Println("app start: ", time.Now())
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for t := range ticker.C {
			fmt.Printf("send at: %s\n", t)
			queue <- true
		}
	}()

	go func() {
		defer wg.Done()

		for {
			r := <-queue
			fmt.Println(r)
		}

	}()

	wg.Wait()
	fmt.Println("exit main")
}

func Channel() {
	var queue = make(chan bool)

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for t := range ticker.C {
			fmt.Println("send:", t)
			queue <- true
		}
	}()

	go func() {
		for {
			fmt.Println("loop1")
			r := <-queue
			fmt.Println("recv1:", r)
		}
	}()

	for {
		fmt.Println("loop")
		r := <-queue
		fmt.Println("recv:", r)
	}
}
