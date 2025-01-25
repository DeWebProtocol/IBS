package routing

import (
	"errors"
	"fmt"
	"math/rand"
)

type FloodTable struct {
	// id:int
	table     map[uint64]PeerInfo
	tableSize int
	degree    int
}

func (t *FloodTable) IsNeighbour(u uint64) bool {
	_, ok := t.table[u]
	return ok
}

func NewFloodTable(tableSize, degree int) Table {
	return &FloodTable{
		make(map[uint64]PeerInfo),
		tableSize,
		degree,
	}
}

func (t *FloodTable) Length() int {
	return len(t.table)
}

func (t *FloodTable) SetTableSize(n int) {
	t.tableSize = n
}

func (t *FloodTable) TableSize() int {
	return t.tableSize
}
func (t *FloodTable) NoRoomForNewPeer(peerID uint64) bool {
	return len(t.table) >= t.tableSize
}

func (t *FloodTable) AddPeer(peerInfo PeerInfo) error {
	if !t.NoRoomForNewPeer(peerInfo.PeerID()) {
		t.table[peerInfo.PeerID()] = peerInfo
	} else {
		s := fmt.Sprintf("adding peer into a full table, size:%d", t.tableSize)
		return errors.New(s)
	}
	return nil
}

func (t *FloodTable) RemovePeer(peerID uint64) {
	delete(t.table, peerID)
}

func randomFrom(degree int, all []uint64) []uint64 {
	if degree > len(all) {
		degree = len(all)
	}
	copi := make([]uint64, len(all))
	selected := make([]uint64, 0, degree)
	copy(copi, all)
	for i := 0; i < degree; i++ {
		r := rand.Intn(len(copi))
		selected = append(selected, copi[r])
		copy(copi[r:], copi[r+1:])
		copi = copi[:len(copi)-1]
	}
	return selected
}

func (t *FloodTable) PeersExcept(eid uint64) PeerInfos {
	ps := make([]PeerInfo, 0, len(t.table))
	for id, info := range t.table {
		if id != eid {
			ps = append(ps, info)
		}
	}
	return ps
}

// TODO: 从n个中随机选择k个
func (t *FloodTable) PeersToBroadcast(from uint64) []uint64 {
	var peers []uint64
	// broadcast to all peers except the sender
	for id, _ := range t.table {
		if id != from {
			peers = append(peers, id)
		}
	}
	return randomFrom(t.degree, peers)
}

func (t *FloodTable) SetLastSeen(id uint64, timestamp int64) error {
	peer, ok := t.table[id]
	if ok {
		peer.SetLastSeen(timestamp)
	}
	return errors.New("flood SetLastSeen: No such peer")
}
func (t *FloodTable) PrintTable() {
	pis := make(PeerInfos, 0, len(t.table))
	for _, info := range t.table {
		pis = append(pis, info)
		//fmt.Printf("%d(%f) ", info.PeerID(), info.Score())
	}
	fmt.Println(pis)
}
