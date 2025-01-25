package routing

import (
	"math"
	"math/rand"
	"sort"
)

type KadcastTable struct {
	*kademlia
	Beta      int
	peerCount int
}

func (t *KadcastTable) IsNeighbour(ID uint64) bool {
	if ID == 0 {
		return false
	}
	_, i := t.Locate(ID)
	if i != -1 {
		return true
	}
	return false
}

func NewKadcastTable(nodeID uint64, k, beta int) Table {
	return &KadcastTable{
		NewKademlia(nodeID, k),
		beta,
		0,
	}
}

func (t *KadcastTable) Length() int {
	return t.peerCount
}

func (t *KadcastTable) NoRoomForNewPeer(peerID uint64) bool {
	b, _ := t.Locate(peerID)
	if len(t.buckets[b]) >= t.k {
		return true
	}
	return false
}

func (t *KadcastTable) SetTableSize(k int) {
	t.k = k
}

func (t *KadcastTable) TableSize() int {
	return math.MaxInt
}

func (t *KadcastTable) AddPeer(peerInfo PeerInfo) error {
	return t.kademlia.AddPeer(peerInfo)
}

func (t *KadcastTable) RemovePeer(peerID uint64) {
	t.kademlia.RemovePeer(peerID)
}

func (t *KadcastTable) PeersToBroadcast(from uint64) []uint64 {
	b := 0
	if from != 0 {
		_b, _ := t.kademlia.Locate(from)
		b = _b + 1
	}
	//fmt.Println("start from bucket", b)
	//t.PrintTable()
	var peers []uint64
	// broadcast to all peers in buckets of subtree that height less than b
	for i := b; i < KeySpaceBits; i++ {
		var idsInBucket []uint64
		for _, info := range t.buckets[i] {
			idsInBucket = append(idsInBucket, info.PeerID())
		}
		ids := randomNFrom(&idsInBucket, t.Beta)
		peers = append(peers, *ids...)
	}
	return peers
}

func randomNFrom(ids *[]uint64, beta int) *[]uint64 {
	if len(*ids) < beta {
		return ids
	}
	var rids []uint64
	for i := 0; i < beta; i++ {
		r := rand.Intn(len(*ids))
		rids = append(rids, (*ids)[r])
		for j := r + 1; j < len(*ids); j++ {
			(*ids)[j-1] = (*ids)[j]
		}
		*ids = (*ids)[:len(*ids)-1]
	}
	return &rids
}

func (t *KadcastTable) SetLastSeen(id uint64, timestamp int64) error {
	return t.kademlia.SetLastSeen(id, timestamp)
}

func PeersInBucket(t *KadcastTable, i int) *[]uint64 {
	return t.kademlia.PeersInBucket(i)
}

func (t *KadcastTable) PrintTable() {
	t.kademlia.PrintBuckets()
}
func (t *KadcastTable) SortPeers() {
	for _, bucket := range t.buckets {
		sort.Sort(bucket)
	}
}
