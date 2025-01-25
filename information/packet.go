package information

import (
	"fmt"

	"github.com/dewebprotocol/IBS/node"
)

// meta is the basic information of a packet.
type meta struct {
	id         int
	timestamp  int64
	dataSize   int // in Byte
	originNode node.Node
}

// Information is the basic information of a packet.
// should not be changed after creation.
type Information struct {
	*meta
	relayNode node.Node
}

func (i *Information) ID() int {
	return i.id
}

func (i *Information) DataSize() int {
	return i.dataSize
}

type Packet interface {
	ID() int
	Timestamp() int64
	From() node.Node
	To() node.Node
}

type Packets []Packet

func (ps Packets) Len() int {
	return len(ps)
}
func (ps Packets) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}
func (ps Packets) Less(i, j int) bool {
	return ps[i].Timestamp() < ps[j].Timestamp()
}

func (p *BasicPacket) Print() {
	fmt.Printf(
		"pacekt: %d %d->%d  size: %dB timestamp: %d propagationDelay: %d transmissionDelay: %d queuingDelaySending: %d\n",
		p.id, p.from.Id(), p.to.Id(), p.dataSize, p.timestamp, p.propagationDelay, p.transmissionDelay, p.queuingDelaySending)
}
