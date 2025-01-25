package routing

import (
	"fmt"
	"testing"
)

func TestLocate(t *testing.T) {
	for i := 1; i < 8; i++ {
		fmt.Printf("peer %d in bucket %d\n", i, locate(uint64(0), uint64(i)))
	}
}

func TestFakeID(t *testing.T) {
	id, err := FakeIDForBucket(5, 63)
	fmt.Println("fakeID: ", id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(locate(5, id))
}
func TestRemove(t *testing.T) {
	arr := []int{0, 1, 2, 3, 4, 5}
	i := 2
	for ; i < len(arr)-1; i++ {
		arr[i] = arr[i+1]
	}
	arr = arr[:i]
	fmt.Println(arr)
}
