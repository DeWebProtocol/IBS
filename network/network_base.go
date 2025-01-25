package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/dewebprotocol/IBS/information"
	"github.com/dewebprotocol/IBS/node"
	"github.com/dewebprotocol/IBS/node/routing"
	"github.com/dewebprotocol/IBS/output"
)

type Region struct {
	Region       string  `json:"region"`
	Distribution float32 `json:"distribution"`
}

type Bandwidth struct {
	UploadBandwidth int     `json:"uploadBandwidth"`
	Distribution    float32 `json:"distribution"`
}

type Delays [][]int32

type BaseNetwork struct {
	newPeerInfo func(node.Node) routing.PeerInfo

	bootNode       node.Node
	Nodes          map[uint64]node.Node
	indexes        map[int]uint64 // order in the network -> id
	DelayOfRegions Delays

	RegionId                    map[string]int
	regions                     []string
	nodeDistribution            []float32
	uploadBandwidths            []int
	uploadBandwidthDistribution []float32

	//lastPacketGeneratedAt int64
	lastPacketIndex     int
	lastOriginNodeIndex int
}

func NewBasicNetwork(bootNode node.Node) *BaseNetwork {
	// unit: μs (0.000,001s)
	net := &BaseNetwork{
		NewBasicPeerInfo,

		bootNode,
		map[uint64]node.Node{},
		map[int]uint64{},
		Delays{},

		make(map[string]int),
		[]string{},
		[]float32{},
		[]int{},
		[]float32{},
		//0,
		0,
		0,
		//[]int{},
	}
	net.loadConf()
	return net
}

func (net *BaseNetwork) ClearState() {
	net.lastPacketIndex = 0
	net.lastOriginNodeIndex = 0
}

func (net *BaseNetwork) loadConf() {
	// DelayOfRegions
	delay, err := os.ReadFile("conf/delay.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(delay, &net.DelayOfRegions)
	if err != nil {
		panic(err)
	}

	// regions
	var regions []Region
	region, err := ioutil.ReadFile("conf/region.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(region, &regions)
	if err != nil {
		panic(err)
	}

	for i, r := range regions {
		net.RegionId[r.Region] = i
		net.regions = append(net.regions, r.Region)
		net.nodeDistribution = append(net.nodeDistribution, r.Distribution)
	}

	// bandwidth
	var bandwidths []Bandwidth
	bandwidth, err := os.ReadFile("conf/bandwidth.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bandwidth, &bandwidths)
	if err != nil {
		panic(err)
	}
	for _, b := range bandwidths {
		net.uploadBandwidthDistribution = append(net.uploadBandwidthDistribution, b.Distribution)
		net.uploadBandwidths = append(net.uploadBandwidths, 1<<b.UploadBandwidth)
	}
}

// generateNodes generate nodes by given newNode function, its region and bandwidth
// is randomly assigned according to configuration files. And add the node into network
func (net *BaseNetwork) generateNodes(n int, newNode func(int, int, string, map[string]int) node.Node, config map[string]int) {
	rd := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 1; i <= n; i++ {
		regionIndex := 0
		r := rd.Float32()
		acc := float32(0)
		for index, f := range net.nodeDistribution {
			if r > acc && r < acc+f {
				regionIndex = index
			}
			acc += f
		}
		bandwidthIndex := 0
		r = rd.Float32()
		acc = float32(0)
		for index, f := range net.uploadBandwidthDistribution {
			if r > acc && r < acc+f {
				bandwidthIndex = index
			}
			acc += f
		}
		_node := newNode(
			i,
			//net.downloadBandwidth[regionIndex],
			net.uploadBandwidths[bandwidthIndex],
			net.regions[regionIndex],
			config,
		)
		net.Add(_node, i)
	}
}

func (net *BaseNetwork) BootNode() node.Node {
	return net.bootNode
}

func (net *BaseNetwork) Node(id uint64) node.Node {
	return net.Nodes[id]
}

// NodeID (index in the network ->nodeID)
func (net *BaseNetwork) NodeID(id int) uint64 {
	return net.indexes[id]
}

// Connect two peers
func (net *BaseNetwork) Connect(a, b node.Node, f NewPeerInfo) bool {
	//fmt.Println("connect", a.Id(), b.Id())
	if a.Id() == b.Id() {
		return false
	}
	//if a.NoRoomForNewPeer(b.Id()) || b.NoRoomForNewPeer(a.Id()) {
	//	return false
	//}
	//fmt.Printf("connect %d to %d\n", a.Id(), b.Id())
	bInfo := f(b)
	aInfo := f(a)
	a.AddPeer(bInfo)
	b.AddPeer(aInfo)
	//fmt.Println("connect", a.Id(), b.Id())
	return true
}

func (net *BaseNetwork) Add(n node.Node, i int) {
	net.Nodes[n.Id()] = n
	net.indexes[i] = n.Id()
}

func (net *BaseNetwork) Size() int {
	return len(net.indexes)
}

// NewPacketGeneration select next node to init a \broadcast at time timestamp
func (net *BaseNetwork) NewPacketGeneration(timestamp int64) information.Packet {
	var origin node.Node
	for i := 0; i <= net.Size(); i++ {
		net.lastOriginNodeIndex = (net.lastOriginNodeIndex)%net.Size() + 1
		origin = net.Node(net.NodeID(net.lastOriginNodeIndex))
		if origin.Running() == true && origin.Malicious() == false {
			break
		}
	}
	p := information.NewBasicPacket(net.lastPacketIndex, 1<<7, origin, net.BootNode(), origin, nil, timestamp)
	// fmt.Printf("node index: %d, timestamp: %d\n", net.lastOriginNodeIndex, timestamp)
	net.lastPacketIndex++
	return p
}

// NodeCrash makes some nodes from i to netSize offline
func (net *BaseNetwork) NodeCrash(i int, once bool) int {
	if once {
		return net.NodeCrashHalf(i)
	}
	return net.NodeCrashRandom(i)
}

func (net *BaseNetwork) NodeCrashRandom(i int) int {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	cnt := 0
	if i < 1 {
		i = 1
	}
	for ; i <= net.Size(); i++ {
		id := net.NodeID(i)
		n := net.Node(id)
		r := rd.Intn(net.Size())
		if n.CrashFactor() >= r {
			cnt++
			n.Stop()
		}
	}
	return cnt
}

func (net *BaseNetwork) NodeCrashHalf(i int) int {
	cnt := 0
	if i < 1 {
		i = 1
	}
	for ; i <= net.Size(); i++ {
		if i%2 == 0 {
			id := net.NodeID(i)
			n := net.Node(id)
			cnt++
			n.Stop()
		}
	}
	return cnt
}

// NodeInfest makes nodes from i to netSize refuse to relay messages
//
//	func (net *BaseNetwork) NodeInfest(i int) int {
//		rd := rand.New(rand.NewSource(time.Now().Unix()))
//		cnt := 0
//		if i < 1 {
//			i = 1
//		}
//		for ; i <= net.Size(); i++ {
//			id := net.NodeID(i)
//			n := net.Node(id)
//			r := rd.Intn(net.Size())
//			if n.CrashFactor() >= r {
//				cnt++
//				n.Infest()
//			}
//		}
//		return cnt
//	}
func (net *BaseNetwork) NodeInfest(i int) int {
	cnt := 0
	if i < 1 {
		i = 1
	}
	for ; i <= net.Size(); i++ {
		if i%2 == 0 {
			id := net.NodeID(i)
			n := net.Node(id)
			cnt++
			n.Infest()
		}
	}
	return cnt
}

func (net *BaseNetwork) ResetNodesReceived() {
	fmt.Println("reset nodes received, len:", len(net.Nodes))
	for _, n := range net.Nodes {
		n.ResetStates()
	}
}

func (net *BaseNetwork) succeedingPackets(p *information.BasicPacket, IDs *[]uint64) information.Packets {
	var packets information.Packets
	sender := p.To()
	if sender.Running() == false {
		return packets
	}
	//if sender.Malicious() == true {
	//	p.SetRedundancy(true)
	//	return packets
	//}
	receivedAt := p.Timestamp()
	//received := sender.Received(p.ID(), p.Timestamp())
	//if received == true {
	//	p.SetRedundancy(true)
	//	//fmt.Printf("%d->%d info=%d hop=%d t=%d μs (redundancy: %t)\n", p.from.Id(), sender.Id(), p.id, p.hop, p.timestamp, p.redundancy)
	//	return packets
	//}
	//fmt.Printf("%d->%d info=%d hop=%d t=%d μs (redundancy: %t)\n", p.from.Id(), sender.Id(), p.id, p.hop, p.timestamp, p.redundancy)
	//IDs := sender.PeersToBroadCast(p.from)
	regionID := net.RegionId
	for _, toID := range *IDs {
		to := net.Node(toID)
		if to.Running() == false {
			continue
		}
		// p.to: sender of next packets
		propagationDelay := net.DelayOfRegions[regionID[sender.Region()]][regionID[to.Region()]]
		bandwidth := sender.UploadBandwidth()
		transmissionDelay := p.DataSize() * 1_000_000 / bandwidth // μs
		var packet *information.BasicPacket

		//if p.from.Id() == p.net.BootNode().Id() {
		//	packet.relayer = to
		//}
		//log.Println("fromID:", p.From().Id())
		if p.From().Id() == BootNodeID {
			//log.Println("set relayNode", to.Id())
			packet = p.NextPacket(to, propagationDelay, int32(transmissionDelay), true)
		} else {
			packet = p.NextPacket(to, propagationDelay, int32(transmissionDelay), false)
		}
		packets = append(packets, packet)
	}
	// add sending queuing delay for each packet
	// sending the packet that is earliest to be received first
	sort.Sort(packets)
	base := int32(0)
	if receivedAt < sender.TsLastSending() {
		base = int32(sender.TsLastSending() - receivedAt)
	}
	for _, packet := range packets {
		bp := packet.(*information.BasicPacket)
		bp.SetAndAddQueuingDelay(base)
		base += bp.TransmissionDelay()
		//packet.to.TsLastReceived = packet.timestamp
	}
	sender.SetTsLastSending(receivedAt + int64(base))
	return packets
}

//	func (net *BaseNetwork) PacketReplacement(p *information.BasicPacket) (information.Packets, int, int) {
//		malicious, total := 0, 0
func (net *BaseNetwork) PacketReplacement(p *information.BasicPacket) information.Packets {
	received := p.To().Received(p.ID(), p.Timestamp())
	if received == true {
		p.SetRedundancy(true)
		//fmt.Printf("%d->%d info=%d hop=%d t=%d μs (redundancy: %t)\n", p.from.Id(), sender.Id(), p.id, p.hop, p.timestamp, p.redundancy)
		return information.Packets{}
	}
	switch sender := p.To().(type) {
	case *node.NeNode:
		if p.From().Malicious() == true {
			fmt.Println("new msg from malicious node")
		}
		sender.NewMsg(p.From().Id())
	}
	var peers = p.To().PeersToBroadCast(p.From())
	crashCnt := 0
	var n node.Node
	for i, peerID := range peers {
		peers[i-crashCnt] = peers[i]
		n = net.Node(peerID)
		if n.Running() == false {
			p.To().RemovePeer(peerID)
			crashCnt++
		}
		//if n.Malicious() {
		//	malicious++
		//}
		//total++
	}
	peers = peers[:len(peers)-crashCnt]
	//if p.ID() == 89 {
	//	fmt.Println(p.From().Id(), "->", p.To().Id())
	//	p.To().PrintTable()
	//	fmt.Println("relay to:", peers)
	//}
	//if total > 10 {
	//fmt.Printf("total:%d, malicious:%d\n", total, malicious)
	//}

	//return net.succeedingPackets(p, &peers), malicious, total
	return net.succeedingPackets(p, &peers)
}

func (net *BaseNetwork) OutputNodes(folder string) {
	outputNodes := output.NewNodeOutput()
	for _, n := range net.Nodes {
		//n.PrintTable()
		outputNodes.Append(n)
	}
	outputNodes.WriteNodes(folder)
}
