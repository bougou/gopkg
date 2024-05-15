package timeutil

import "time"

// Ticker interface 报时器
// 1. 周期性发送时间
// 2. 不会超时，除非主动停止
type Ticker interface {
	C() <-chan time.Time
	Stop()
}

type realTicker struct {
	ticker *time.Ticker
}

func (t *realTicker) C() <-chan time.Time {
	return t.ticker.C
}

func (t *realTicker) Stop() {
	t.ticker.Stop()
}

type fakeTicker struct {
	c <-chan time.Time
}

func (t *fakeTicker) C() <-chan time.Time {
	return t.c
}

func (t *fakeTicker) Stop() {
}
