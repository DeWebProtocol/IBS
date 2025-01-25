package node

import (
	"github.com/dewebprotocol/IBS/node/routing"
)

type BasicNode struct {
	id              uint64
	region          string
	routingTable    routing.Table
	uploadBandwidth int // byte/s
	crashFactor     int
	crashTimes      int

	//TsLastReceived int64 // the time(μs) when last packet was received
	tsLastSending int64 // the time(μs) when last packet was sent

	receivedPackets map[int]int64 // id -> delay
	running         bool
	malicious       bool
}

func NewBasicNode(id uint64, uploadBandwidth, crashFactor int, region string, table routing.Table) *BasicNode {
	return &BasicNode{
		id,
		region,
		table,
		//downloadBandwidth,
		uploadBandwidth,
		crashFactor,
		0,
		0,
		map[int]int64{},
		true,
		false,
	}
}

func (n *BasicNode) ResetRoutingTable(table routing.Table) {
	n.routingTable = table
}

func (n *BasicNode) ResetStates() {
	n.crashTimes = 0
	n.tsLastSending = 0
	n.receivedPackets = map[int]int64{}
}

func (n *BasicNode) Id() uint64 {
	return n.id
}
func (n *BasicNode) Region() string {
	return n.region
}

//	func (n *BasicNode) DownloadBandwidth() int {
//		return n.downloadBandwidth
//	}
func (n *BasicNode) UploadBandwidth() int {
	return n.uploadBandwidth
}
func (n *BasicNode) TsLastSending() int64 {
	return n.tsLastSending
}
func (n *BasicNode) SetTsLastSending(t int64) {
	n.tsLastSending = t
}
func (n *BasicNode) SetLastSeen(id uint64, timestamp int64) error {
	return n.routingTable.SetLastSeen(id, timestamp)
}

//func (n *BasicNode) LastDelay() int64 {
//	return n.lastDelay
//}
//func (n *BasicNode) SetLastDelay(delay int64) bool {
//	if delay < n.lastDelay {
//		return false
//	}
//	n.lastDelay = delay
//	return true
//}

func (n *BasicNode) NoRoomForNewPeer(peerID uint64) bool {
	return n.routingTable.NoRoomForNewPeer(peerID)
}

func (n *BasicNode) RoutingTableLength() int {
	return n.routingTable.Length()
}

func (n *BasicNode) AddPeer(peerInfo routing.PeerInfo) bool {
	err := n.routingTable.AddPeer(peerInfo)
	if err != nil {
		//fmt.Println(err)
		return false
	}
	return true
}

func (n *BasicNode) RemovePeer(peerID uint64) {
	n.routingTable.RemovePeer(peerID)
}

// return id of peers
func (n *BasicNode) PeersToBroadCast(from Node) []uint64 {
	peerIDs := n.routingTable.PeersToBroadcast(from.Id())
	return peerIDs
}

func (n *BasicNode) Received(msgId int, timestamp int64) bool {
	_, ok := n.receivedPackets[msgId]
	if timestamp == -1 { // just return whether this msg was received
		return ok
	}
	if !ok {
		n.receivedPackets[msgId] = timestamp
	}
	return ok
}
func (n *BasicNode) NumReceivedPackets() int {
	return len(n.receivedPackets)
}

func (n *BasicNode) Running() bool {
	return n.running
}
func (n *BasicNode) PrintTable() {
	n.routingTable.PrintTable()
}

func (n *BasicNode) Run() {
	n.running = true
}
func (n *BasicNode) Stop() {
	n.crashTimes++
	n.running = false
}
func (n *BasicNode) CrashFactor() int {
	return n.crashFactor
}
func (n *BasicNode) CrashTimes() int {
	return n.crashTimes
}
func (n *BasicNode) Infest() {
	n.malicious = true
}
func (n *BasicNode) Malicious() bool {
	return n.malicious
}
