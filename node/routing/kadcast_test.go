package routing

import (
	"fmt"
	"sort"
	"testing"
)

func TestXOR(t *testing.T) {
	a := 7
	b := 1
	fmt.Println(a ^ b)
}

func TestPeerInfoSort(t *testing.T) {
	var peerInfos PeerInfos
	info1 := NewBasicPeerInfo(1)
	info1.SetLastSeen(int64(10))
	peerInfos = append(peerInfos, info1)
	info2 := NewBasicPeerInfo(2)
	info2.SetLastSeen(int64(15))
	peerInfos = append(peerInfos, info2)
	info3 := NewBasicPeerInfo(3)
	info3.SetLastSeen(int64(5))
	peerInfos = append(peerInfos, info3)
	sort.Sort(peerInfos)
	for _, info := range peerInfos {
		fmt.Println(info.LastSeen())
	}
}

func TestBucketLocating(t *testing.T) {
	b := 64
	xor := uint64(1 << 63)
	for xor != 0 {
		xor = xor >> 1
		b--
	}
	fmt.Println(b)
}

func TestRandomNFrom(t *testing.T) {
	ids := []uint64{1, 2, 3, 4, 5, 6, 7}
	fmt.Println(*randomNFrom(&ids, 3))
}
