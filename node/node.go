package node

import "github.com/dewebprotocol/IBS/node/routing"

type Node interface {
	ResetRoutingTable(routing.Table)
	ResetStates()
	Id() uint64
	Region() string
	//DownloadBandwidth() int
	UploadBandwidth() int
	TsLastSending() int64
	SetTsLastSending(int64)
	SetLastSeen(id uint64, timestamp int64) error
	NoRoomForNewPeer(id uint64) bool
	RoutingTableLength() int
	AddPeer(routing.PeerInfo) bool
	RemovePeer(uint64)
	PeersToBroadCast(Node) []uint64
	Received(int, int64) bool
	NumReceivedPackets() int
	Running() bool
	PrintTable()
	Run()
	Stop()
	CrashFactor() int
	CrashTimes() int
	Infest()
	Malicious() bool
}
