package snippet

import (
	"fmt"
	"sync"
	"testing"
)

func Test_test1(t *testing.T) {
	type T int

	// 保持结果的数组
	var res []T

	// 并发的协程把结果发送到这个队列中
	queue := make(chan T)

	// 控制并发的协程
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(j int) {
			defer wg.Done()

			fmt.Printf("Goroutine %d\n", j)

			a := []int{1, 2, 3, 4, 5}
			for _, v := range a {
				e := 10*j + v

				queue <- T(e)
			}
		}(i)
	}

	go func() {
		fmt.Println("Wait for all goroutines stop")
		wg.Wait()
		close(queue)
	}()

	for t := range queue {
		res = append(res, t)
	}

	fmt.Println(res)
}

func Test_test2(t *testing.T) {
	type T int

	var res []T
	var wg sync.WaitGroup
	queue := make(chan []T)

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(j int) {
			defer wg.Done()

			fmt.Printf("Goroutine %d\n", j)

			a := []int{1, 2, 3, 4, 5}
			var b []T
			for _, v := range a {
				e := 10*j + v
				b = append(b, T(e))
			}
			queue <- b
		}(i)
	}

	go func() {
		fmt.Println("Wait for all goroutines stop")
		wg.Wait()
		close(queue)
	}()

	for t := range queue {
		res = append(res, t...)
	}

	fmt.Println(res)
}
