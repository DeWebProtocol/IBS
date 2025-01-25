package information

import (
	"container/heap"
	"fmt"
	"testing"
)

func TestInfoSorter(t *testing.T) {
	s := []int64{1, 3, 9, 5, 3, 6, 8, 1, 0, 1, 3, 6, 100, 34, 36, 67, 87, 33, 10, 4, 5, 8, 4, 9, 0, 7, 5, 4, 3, 2, 3, 5, 6, 7, 3, 1}
	//var infos []*information.BasicPacket
	sorter := NewInfoSorter()
	for i, _ := range s {
		p := NewBasicPacket(0, 0, nil, nil, nil, nil, s[i])
		heap.Push(sorter, p)
		fmt.Println(i, "append", p.Timestamp(), sorter.Len())
	}
	sorter.print()
	for sorter.Len() > 0 {
		info := heap.Pop(sorter).(Packet)
		fmt.Println("take", info.Timestamp())
	}
}
