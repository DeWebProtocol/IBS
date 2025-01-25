package routing

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const KeySpaceBits = 64

type kademlia struct {
	nodeID  uint64
	buckets [KeySpaceBits]PeerInfos
	k       int // bucket size
}

func NewKademlia(nodeID uint64, k int) *kademlia {
	return &kademlia{
		nodeID,
		[KeySpaceBits]PeerInfos{},
		k,
	}
}

func fakeIDForBucket(num uint64, b int) uint64 {
	var base uint64 = 1 << b
	rd := rand.New(rand.NewSource(time.Now().Unix()))
	if b == KeySpaceBits-1 {
		// random 63bit positive number
		r := uint64(rd.Int63n(KeySpaceBits-2))<<1 + uint64(rd.Intn(2))
		return num ^ (base + r)
	}
	return num ^ (uint64(rd.Int63n(int64(base))) + base)
}

func FakeIDForBucket(nodeID uint64, b int) (uint64, error) {
	b = KeySpaceBits - 1 - b
	if b < 0 || b > KeySpaceBits-1 {
		return nodeID, errors.New("FakeIDForBucket: out of range")
	}
	if b == 0 {
		return nodeID ^ 1, nil
	}
	return fakeIDForBucket(nodeID, b), nil
}

func (k *kademlia) SetLastSeen(id uint64, timestamp int64) error {
	b, i := k.Locate(id)
	if i != -1 {
		k.buckets[b][i].SetLastSeen(timestamp)
		return nil
	}
	str := fmt.Sprintf("kademlia SetLastSeen: %d failed to find peer %d", k.nodeID, id)
	return errors.New(str)
}

func locate(k1, k2 uint64) int {

	xor := k1 ^ k2
	b := KeySpaceBits
	for xor != 0 {
		xor = xor >> 1
		b--
	}
	//fmt.Println("locate", k1, k2, b)
	return b
}

func (k *kademlia) Locate(peerID uint64) (int, int) {
	b := locate(k.nodeID, peerID)
	for i, info := range k.buckets[b] {
		if info.PeerID() == peerID {
			return b, i
		}
	}
	return b, -1
}

func (k *kademlia) AddPeer(info PeerInfo) error {
	b, i := k.Locate(info.PeerID())
	if i != -1 {
		// already included
		return errors.New("kademlia AddPeer: peer already included")
	}
	if k.buckets[b].Len() < k.k {
		k.buckets[b] = append(k.buckets[b], info)
		//sort.Sort(k.buckets[b])
		return nil
	}

	//last := k.buckets[b][k.buckets[b].Len()-1]
	//if last.Score() < info.Score() {
	//	k.buckets[b][k.buckets[b].Len()-1] = info
	//	sort.Sort(k.buckets[b])
	//	return nil
	//}
	return errors.New("kademlia AddPeer: older peers are preferred")
}

func (k *kademlia) RemovePeer(peerID uint64) {
	b, i := k.Locate(peerID)
	if i == -1 {
		return
	}
	for ; i < len(k.buckets[b])-1; i++ {
		k.buckets[b][i] = k.buckets[b][i+1]
	}
	k.buckets[b] = k.buckets[b][:i]
}

func (k *kademlia) PeersInBucket(b int) *[]uint64 {
	ids := []uint64{}
	for _, info := range k.buckets[b] {
		ids = append(ids, info.PeerID())
	}
	return &ids
}

func (k *kademlia) PrintBuckets() {
	for i, bucket := range k.buckets {
		if len(bucket) > 0 {
			fmt.Print("bucket", i, ": ")
			for _, info := range bucket {
				//fmt.Printf("distance=%d(%d), ", info.PeerID()^k.nodeID, info.PeerID())
				fmt.Printf("score=%f(%d), ", info.Score(), info.PeerID())
			}
			fmt.Println()
		}
	}
}
