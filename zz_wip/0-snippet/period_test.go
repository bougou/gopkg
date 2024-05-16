package snippet

import (
	"fmt"
	"testing"
	"time"
)

func DoPeriodically() {

	fmt.Printf("app start: %s\n", time.Now().Format("2006 01 02 15:04:05"))

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	var i int

	// 第一次获得时间是在第一个周期之后，也就是说第一次 tick 是延后的
	// 如果想让第一次 loop 是立即的，可以用下面写法，但是第一次loop 是没有获得ticker 时间的
	// var t time.Time
	// for ; true; t = <-ticker.C {
	// }

	for t := range ticker.C {
		i++
		fmt.Println()
		fmt.Printf("ticker run at, %s\n", time.Now().Format("2006 01 02 15:04:05"))
		fmt.Printf("ticker receive time, %s\n", t.Format("2006 01 02 15:04:05"))

		// do something
		time.Sleep(5 * time.Second)

		if i >= 5 {
			break
		}
	}
	// 但是，取决于 loop 内程序逻辑的执行时间，下次 loop 的 「运行时间」是不确定的。
	// 但是「运行时间」只会晚于「接收时间」
	// 如果 loop 内逻辑很快完成(在 ticker 周期内），for 会 block 直到 ticker 的下一次 tick 发生。
	// 如果 loop 内运行很慢，
	// 注意，for loop 每次收到的时间不仅依赖于 ticker 周期，也依赖于 loop 运行时间

}

func DoIntervally() {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	var t time.Time
	for i := 0; true; t = <-ticker.C {
		i++
		if i >= 5 {
			break
		}
		fmt.Println()
		fmt.Printf("run time: %s\n", time.Now().Format("2006 01 02 15:04:05"))
		fmt.Printf("received time: %s\n", t.Format("2006 01 02 15:04:05"))

		time.Sleep(5 * time.Second)
		fmt.Printf("reset time: %s\n", time.Now().Format("2006 01 02 15:04:05"))
		ticker.Reset(time.Second * 2)
	}

	// run time: 2021 10 21 16:12:11
	// received time: 0001 01 01 00:00:00
	// reset time: 2021 10 21 16:12:16   //

	// run time: 2021 10 21 16:12:16
	// received time: 2021 10 21 16:12:13    //
	// reset time: 2021 10 21 16:12:21

	// run time: 2021 10 21 16:12:21
	// received time: 2021 10 21 16:12:18
	// reset time: 2021 10 21 16:12:26

	// run time: 2021 10 21 16:12:26
	// received time: 2021 10 21 16:12:23
	// reset time: 2021 10 21 16:12:31
}

func Test_DoPeriodically(t *testing.T) {
	DoPeriodically()
}

func Test_DoIntervally(t *testing.T) {
	DoIntervally()
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func Test_fibonacci(t *testing.T) {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}
