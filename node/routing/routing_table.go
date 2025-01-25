package routing

import (
	"fmt"
	"sort"
)

type Table interface {
	Length() int
	SetTableSize(int)
	TableSize() int
	NoRoomForNewPeer(peerID uint64) bool
	AddPeer(PeerInfo) error
	RemovePeer(uint64)
	PeersToBroadcast(from uint64) []uint64
	SetLastSeen(uint64, int64) error // peerID, timestamp
	PrintTable()
	IsNeighbour(uint64) bool
}

type NeTable interface {
	Table
	IncrementNewMsg(uint64)
	IncrementConfirmation(uint64)
	IncrementReceivedConfirmation(uint64)
}

// PeerInfo : rank from the peer
//type PeerInfo struct {
//	NodeID  int
//}

type PeerInfo interface {
	PeerID() uint64
	Score() float64 // higher score, higher priority
	SetLastSeen(int64)
	LastSeen() int64
}

type PeerInfos []PeerInfo

func (ps PeerInfos) String() string {
	sort.Sort(ps)
	str := ""
	for _, p := range ps {
		str += fmt.Sprintf("%d(%f) ", p.PeerID(), p.Score())
	}
	return str
}

func (ps PeerInfos) Len() int {
	return len(ps)
}
func (ps PeerInfos) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}
func (ps PeerInfos) Less(i, j int) bool {
	return ps[i].Score() > ps[j].Score()
}

func (ps PeerInfos) Includes(p PeerInfo) bool {
	for i := 0; i < ps.Len(); i++ {
		if ps[i].PeerID() == p.PeerID() {
			return true
		}
	}
	return false
}

type BasicPeerInfo struct {
	id       uint64
	lastSeen int64
}

func NewBasicPeerInfo(id uint64) *BasicPeerInfo {
	return &BasicPeerInfo{id, 0}
}

func (i *BasicPeerInfo) PeerID() uint64 {
	return i.id
}

func (i *BasicPeerInfo) Score() float64 {
	return float64(i.LastSeen())
}

func (i *BasicPeerInfo) SetLastSeen(lastSeen int64) {
	i.lastSeen = lastSeen
}
func (i *BasicPeerInfo) LastSeen() int64 {
	return i.lastSeen
}

type NePeerInfo struct {
	*BasicPeerInfo
	newMsg               int
	confirmation         int
	receivedConfirmation int
}

func NewNePeerInfo(id uint64) *NePeerInfo {
	return &NePeerInfo{
		NewBasicPeerInfo(id),
		0,
		0,
		0,
	}
}
func (n *NePeerInfo) NewMsg() {
	n.newMsg += 1
}
func (n *NePeerInfo) Confirmation() {
	n.confirmation += 1
}
func (n *NePeerInfo) ReceivedConfirmation() {
	n.receivedConfirmation += 1
}
func (n *NePeerInfo) Score() float64 {
	return float64(n.newMsg + n.confirmation + n.receivedConfirmation)
}
