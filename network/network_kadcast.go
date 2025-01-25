package network

import (
	"log"

	"github.com/dewebprotocol/IBS/network/num_set"
	"github.com/dewebprotocol/IBS/node"
	"github.com/dewebprotocol/IBS/node/hash"
	"github.com/dewebprotocol/IBS/node/routing"
)

type KadcastNet struct {
	k, beta int
	*BaseNetwork
	idSet *num_set.Set
}

func NewKadcastNode(index int, uploadBandwidth int, region string, config map[string]int) node.Node {
	nodeID := hash.Hash64(uint64(index))
	//nodeID := uint64(index)
	return node.NewBasicNode(
		nodeID,
		uploadBandwidth,
		index,
		region,
		routing.NewKadcastTable(nodeID, config["k"], config["beta"]),
	)
}
func NewKadcastNet(size, k, beta int) Network {
	//fmt.Println("===== kademlia =====")
	//fmt.Println("beta:", beta, "bucket size:", k)
	// bootNode is used for message generation (from node) only here
	bootNode := node.NewBasicNode(BootNodeID, 0, 0, "", nil)
	net := NewBasicNetwork(bootNode)
	config := map[string]int{"k": k, "beta": beta}
	net.generateNodes(size, NewKadcastNode, config)
	kNet := &KadcastNet{
		k,
		beta,
		net,
		num_set.NewSet(),
	}
	kNet.initConnections(NewBasicPeerInfo)
	return kNet
}

// Introduce : return n nodes
func (kNet *KadcastNet) Introduce(id uint64, n int) []node.Node {
	var nodes []node.Node
	for b := 0; b < routing.KeySpaceBits; b++ {
		fakeID, err := routing.FakeIDForBucket(id, b)

		if err != nil {
			log.Fatal(err)
		}
		// 2*n+1 make sures the right keys for the bucket are covered
		// a crude node discovery implementation
		ids := kNet.idSet.Around(fakeID, 2*n+1)
		for _, u := range ids {
			if kNet.Node(u).Running() == true {
				nodes = append(nodes, kNet.Node(u))
			}
		}
		//fmt.Println("introduce to:", id, "peerIDS:", ids)
	}
	return nodes
}

func (kNet *KadcastNet) initConnections(f NewPeerInfo) {
	for _, n := range kNet.Nodes {
		kNet.idSet.Insert(n.Id())
	}
	for _, n := range kNet.Nodes {
		kNet.introduceAndConnect(n, f)
	}
}

func (kNet *KadcastNet) introduceAndConnect(n node.Node, f NewPeerInfo) {
	peers := kNet.Introduce(n.Id(), kNet.k)
	for _, peer := range peers {
		kNet.Connect(n, peer, f)
	}
}

func (kNet *KadcastNet) churn(crashFrom int, once bool,
	routing func(nodeID uint64, k, beta int) routing.Table,
	peerInfo NewPeerInfo) int {
	for _, n := range kNet.Nodes {
		if n.Running() == false {
			// it can be seen as the crashed nodes leave the network
			n.ResetRoutingTable(routing(n.Id(), kNet.k, kNet.beta))
			n.Run()
			kNet.introduceAndConnect(n, peerInfo)
		}
	}
	return kNet.NodeCrash(crashFrom, once)
}

func (kNet *KadcastNet) Churn(crashFrom int, once bool) int {
	return kNet.churn(crashFrom, once, routing.NewKadcastTable, NewBasicPeerInfo)
}

func (kNet *KadcastNet) Infest(crashFrom int) int {
	return kNet.NodeInfest(crashFrom)
}
