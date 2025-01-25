package routing

import "sort"

type NeFloodTable struct {
	FloodTable
}

func NewNeFloodTable(tableSize, degree int) Table {
	return &NeFloodTable{
		FloodTable{
			map[uint64]PeerInfo{},
			tableSize,
			degree,
		},
	}
}

func (t *NeFloodTable) PeersToBroadcast(from uint64) []uint64 {
	ps := t.PeersExcept(from)
	sort.Sort(ps)
	return randomPeersBasedOnScore(ps, t.degree)
}

func (t *NeFloodTable) necastPeerInfo(ID uint64) *NePeerInfo {
	pi, ok := t.table[ID]
	if !ok {
		return nil
	}
	return pi.(*NePeerInfo)
}

func (t *NeFloodTable) IncrementNewMsg(ID uint64) {
	pi := t.necastPeerInfo(ID)
	if pi == nil {
		return
	}
	pi.NewMsg()
}
func (t *NeFloodTable) IncrementConfirmation(ID uint64) {
	pi := t.necastPeerInfo(ID)
	if pi == nil {
		return
	}
	pi.Confirmation()
}
func (t *NeFloodTable) IncrementReceivedConfirmation(ID uint64) {
	//fmt.Println("confirm")
	pi := t.necastPeerInfo(ID)
	if pi == nil {
		return
	}
	pi.ReceivedConfirmation()
}
