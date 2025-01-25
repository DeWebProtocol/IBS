package node

import (
	"github.com/dewebprotocol/IBS/node/routing"
)

type NeNode struct {
	*BasicNode
	//Tasks map[int]*BroadcastTask // metaInfo id -> task
}

func NewNeNode(id uint64, uploadBandwidth, crashFactor int, region string, table routing.Table) *NeNode {
	n := &NeNode{
		NewBasicNode(id, uploadBandwidth, crashFactor, region, table),
		//make(map[int]*BroadcastTask),
	}
	return n
}

func (n *NeNode) IsNeighbour(ID uint64) bool {
	//return n.routingTable.(*routing.NeCastTable).IsNeighbour(ID)
	return n.routingTable.IsNeighbour(ID)
}
func (n *NeNode) NewMsg(peerID uint64) {
	table := n.routingTable.(routing.NeTable)
	table.IncrementNewMsg(peerID)
}
func (n *NeNode) Confirmation(peerID, relayID uint64) {
	table := n.routingTable.(routing.NeTable)
	table.IncrementConfirmation(peerID)
	table.IncrementReceivedConfirmation(relayID)
}

//func (n *NeNode) SortPeers() {
//	table := n.routingTable.(*routing.NeCastTable)
//	table.SortPeers()
//}
