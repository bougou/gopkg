package bench

import "testing"

func BenchmarkCountSerial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		countFunc()
	}
}

func BenchmarkCountParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			countFunc()
		}
	})
}

func BenchmarkCountParallelMultiGoroutines(b *testing.B) {
	b.SetParallelism(1000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			countFunc()
		}
	})
}
