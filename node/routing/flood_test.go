package routing

import (
	"fmt"
	"testing"
)

func TestRandomFrom(t *testing.T) {
	ints := []uint64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	randoms := randomFrom(12, ints)
	fmt.Println(randoms)
}

func BenchmarkRandFrom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ints := []uint64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
		randomFrom(12, ints)
	}
}
