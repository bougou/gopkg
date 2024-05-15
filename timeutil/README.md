# clock

- clock 时钟
- timer 定时器，计时器（有到期时间，会超时）
  - `Timer` interface
- ticker 滴答报时器（周期性报告时间，不会超时，除非主动 Stop）
  - `Ticker` interface（标准库 `time.Ticker` 是 struct）
