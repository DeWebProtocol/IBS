package information

import (
	"fmt"
)

type PacketSorter struct {
	minHeap []Packet
}

func NewInfoSorter() *PacketSorter {
	s := &PacketSorter{
		make([]Packet, 0, 1<<16),
	}
	return s
}

func (s *PacketSorter) Len() int {
	return len(s.minHeap)
}
func (s *PacketSorter) Less(i, j int) bool {
	return s.minHeap[i].Timestamp() < s.minHeap[j].Timestamp()
}
func (s *PacketSorter) Swap(i, j int) {
	s.minHeap[i], s.minHeap[j] = s.minHeap[j], s.minHeap[i]
}
func (s *PacketSorter) Push(x interface{}) {
	s.minHeap = append(s.minHeap, x.(Packet))
}
func (s *PacketSorter) Pop() (v interface{}) {
	v = s.minHeap[len(s.minHeap)-1]
	s.minHeap = s.minHeap[:len(s.minHeap)-1]
	return
}

func (s *PacketSorter) print() {
	for i := 0; i < s.Len(); i++ {
		fmt.Print(s.minHeap[i].Timestamp(), " ")
	}
	fmt.Println()
}
