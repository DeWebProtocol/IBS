package network

import (
	"github.com/dewebprotocol/IBS/node"
	"github.com/dewebprotocol/IBS/node/routing"
)

type NeFloodNet struct {
	*FloodNet
}

func NewNeFloodNode(index int, uploadBandwidth int, region string, config map[string]int) node.Node {
	return node.NewNeNode(
		uint64(index),
		uploadBandwidth,
		index,
		region,
		routing.NewNeFloodTable(config["tableSize"], config["degree"]),
	)
}

func NewNeFloodNet(size, tableSize, degree int) Network {
	// bootNode is used for message generation (from node) only here
	bootNode := node.NewNeNode(BootNodeID, 0, 0, "", nil)
	net := NewBasicNetwork(bootNode)
	config := map[string]int{"tableSize": tableSize, "degree": degree}
	net.generateNodes(size, NewNeFloodNode, config)
	nNet := &NeFloodNet{
		&FloodNet{
			tableSize,
			degree,
			net,
		},
	}
	nNet.initConnections(NewNePeerInfo)
	//for u, n := range net.Nodes {
	//	fmt.Println(u)
	//	n.PrintTable()
	//}
	return nNet
}

func (nNet *NeFloodNet) Churn(crashFrom int, once bool) int {
	return nNet.churn(crashFrom, once, routing.NewNeFloodTable, NewNePeerInfo)
}

//func (nNet *NeFloodNet) PacketReplacement(p *information.BasicPacket) information.Packets {
//	packets := nNet.BaseNetwork.PacketReplacement(p)
//	neNode := p.To().(*node.NeNode)
//	if neNode.Id() != p.Origin().Id() && neNode.IsNeighbour(p.Origin().Id()) && neNode.Id() != p.Relay().Id() {
//		packets = append(packets, p.ConfirmPacket())
//	}
//	return packets
//}
