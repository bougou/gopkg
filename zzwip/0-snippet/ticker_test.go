package snippet

import (
	"fmt"
	"testing"
	"time"
)

func Test_Ticker_ForWaitTick(t *testing.T) {
	var period = 5 * time.Second

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	fmt.Println(time.Now())

	// the first loop will wait the ticker fires
	var i int
	for tt := range ticker.C {
		i++
		fmt.Printf("loop %d, %s hello\n", i, tt)
	}
}

func Test_Ticker_ForImmediately(t *testing.T) {
	var period = 5 * time.Second

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	fmt.Println(time.Now())

	// the first loop will start immediately, then wait unit the ticker fires.
	var tt time.Time
	var i int
	for ; true; tt = <-ticker.C {
		i++
		fmt.Printf("loop %d, %s, hello\n", i, tt)
	}
}
