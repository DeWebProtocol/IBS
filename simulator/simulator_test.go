package simulator

import (
	"fmt"
	"testing"
)

func TestCrashRm(t *testing.T) {
	var arr = []int{0, 1, 3, 5, 6, 4, 8, 9, 0, 7, 6, 3, 1, 6}
	cnt := 0
	for i := 0; i < len(arr); i++ {
		arr[i-cnt] = arr[i]
		if arr[i]%2 == 0 {
			cnt++
		}
	}
	arr = arr[:len(arr)-cnt]
	fmt.Println(arr)
}

func TestSlice(t *testing.T) {
	var s []int = []int{1, 2, 3}
	writeSlice(s)
	fmt.Println(s)
}
func writeSlice(s []int) {
	s[0] = 0
}
