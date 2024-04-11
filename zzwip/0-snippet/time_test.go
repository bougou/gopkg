package snippet

import (
	"fmt"
	"testing"
	"time"
)

func Test_Time(tt *testing.T) {
	fmt.Println("Hello, playground")
	a := []int{1, 2, 3}
	fmt.Println(a)
	fmt.Printf("%T\n", a)
	fmt.Println(a[:2])

	fmt.Println("now:", time.Now())
	fmt.Println("now: ", time.Now().Format("2006.01.02"))
	fmt.Println("now sec: ", time.Now().Unix())
	fmt.Println("now nanosec: ", time.Now().UnixNano())

	tStr := "07-29-2020 17:15:26"
	var t time.Time

	tLayout := "01-02-2006 15:04:05"
	// time.Parse always interpret the time str as UTC
	t, _ = time.Parse(tLayout, tStr)
	fmt.Println(t) // 2020-07-29 17:15:26 +0000 UTC

	// 将 time.Time 表示成其它timezone
	loc, _ := time.LoadLocation("Local") // Asia/Shanghai
	t1 := t.In(loc)
	fmt.Println(t1) // 2020-07-30 01:15:26 +0800 CST

	t1, _ = time.ParseInLocation(tLayout, tStr, loc)
	fmt.Println(t1) // 2020-07-29 17:15:26 +0800 CST

	loc, _ = time.LoadLocation("America/New_York")
	t1 = t.In(loc)
	fmt.Println(t1) // 2020-07-29 13:15:26 -0400 EDT

	// Specify timezone
	t, _ = time.Parse("01-02-2006 15:04:05 -0700", tStr+" +0800")
	fmt.Println(t) // 2020-07-29 17:15:26 +0800 CST

}
