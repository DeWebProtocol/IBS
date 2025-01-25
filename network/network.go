package network

import (
	"github.com/dewebprotocol/IBS/information"
	"github.com/dewebprotocol/IBS/node"
	"github.com/dewebprotocol/IBS/node/routing"
)

const BootNodeID = 0

type NewPeerInfo func(node.Node) routing.PeerInfo

func NewBasicPeerInfo(n node.Node) routing.PeerInfo {
	return routing.NewBasicPeerInfo(n.Id())
}

type Network interface {
	ClearState()
	BootNode() node.Node
	Node(id uint64) node.Node
	NodeID(i int) uint64
	Connect(a, b node.Node, f NewPeerInfo) bool
	Add(n node.Node, i int)
	Size() int
	NodeCrash(i int, once bool) int
	NodeInfest(i int) int
	ResetNodesReceived()
	NewPacketGeneration(timestamp int64) information.Packet
	succeedingPackets(p *information.BasicPacket, IDs *[]uint64) information.Packets
	//PacketReplacement(p *information.BasicPacket) (information.Packets, int, int)
	PacketReplacement(p *information.BasicPacket) information.Packets
	Churn(crashFrom int, once bool) int
	Infest(infestFrom int) int
	OutputNodes(string)
}
