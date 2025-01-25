package hash

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	int1 := Hash64(uint64(1))
	int2 := Hash64(uint64(2))
	fmt.Println(int1, int2, int1^int2)
}
