package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/dewebprotocol/IBS/information"
)

type Packet struct {
	Id                  int
	Timestamp           int64 `json:"timestamp"` // delay(μs) from the generation(timestamp) of information
	PropagationDelay    int32 `json:"propagationDelay"`
	TransmissionDelay   int32 `json:"transmissionDelay"`
	QueuingDelaySending int32 `json:"queuingDelaySending"`
	//QueuingDelayReceiving int32  `json:"queuingDelayReceiving"`
	From       string `json:"from"`
	To         string `json:"to"`
	Hop        int    `json:"hop"`
	Redundancy bool   `json:"redundancy"`
}

func newPacket(p *information.BasicPacket) *Packet {
	return &Packet{
		p.ID(),
		p.Timestamp(),
		p.PropagationDelay(),
		p.TransmissionDelay(),
		p.QueuingDelaySending(),
		//p.QueuingDelayReceiving(),
		strconv.FormatUint(p.From().Id(), 10),
		strconv.FormatUint(p.To().Id(), 10),
		p.Hop(),
		p.Redundancy(),
	}
}

type PacketOutput []*Packet

func NewPacketOutput() PacketOutput {
	return PacketOutput{}
}
func (o *PacketOutput) Append(p *information.BasicPacket) {
	*o = append(*o, newPacket(p))
}

func (o *PacketOutput) WritePackets(folder string) {
	b, err := json.Marshal(o)
	if err != nil {
		fmt.Println(err)
	}
	filename := fmt.Sprintf("%s/output_packets.json", folder)
	err = os.WriteFile(filename, b, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
