package network

import (
	"github.com/dewebprotocol/IBS/information"
	"github.com/dewebprotocol/IBS/network/num_set"
	"github.com/dewebprotocol/IBS/node"
	"github.com/dewebprotocol/IBS/node/hash"
	"github.com/dewebprotocol/IBS/node/routing"
)

type NecastNet struct {
	*KadcastNet
}

func NewNePeerInfo(n node.Node) routing.PeerInfo {
	return routing.NewNePeerInfo(n.Id())
}
func NewNecastNode(index int, uploadBandwidth int, region string, config map[string]int) node.Node {
	nodeID := hash.Hash64(uint64(index))
	return node.NewNeNode(
		nodeID,
		uploadBandwidth,
		index,
		region,
		routing.NewNecastTable(nodeID, config["k"], config["beta"]),
	)
}

func NewNecastNet(size, k, beta int) Network {
	// bootNode is used for message generation (from node) only here
	bootNode := node.NewBasicNode(BootNodeID, 0, 0, "", nil)
	net := NewBasicNetwork(bootNode)
	config := map[string]int{"k": k, "beta": beta}
	net.generateNodes(size, NewNecastNode, config)
	nNet := &NecastNet{
		&KadcastNet{
			k,
			beta,
			net,
			num_set.NewSet(),
		},
	}
	nNet.initConnections(NewNePeerInfo)
	return nNet
}

func (nNet *NecastNet) Churn(crashFrom int, once bool) int {
	return nNet.churn(crashFrom, once, routing.NewNecastTable, NewNePeerInfo)
}

//	func (nNet *NecastNet) PacketReplacement(p *information.BasicPacket) (information.Packets, int, int) {
//		packets, malicious, total := nNet.BaseNetwork.PacketReplacement(p)
func (nNet *NecastNet) PacketReplacement(p *information.BasicPacket) information.Packets {
	packets := nNet.BaseNetwork.PacketReplacement(p)
	neNode := p.To().(*node.NeNode)
	if neNode.Id() != p.Origin().Id() && neNode.IsNeighbour(p.Origin().Id()) && neNode.Id() != p.Relay().Id() {
		packets = append(packets, p.ConfirmPacket())
	}
	return packets
}
