package bench

import "testing"

func BenchmarkByArchiver(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ByArchiver()
	}
}

func BenchmarkByUnzip(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ByUnzip()
	}
}
