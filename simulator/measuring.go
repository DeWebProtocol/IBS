package simulator

import (
	"github.com/dewebprotocol/IBS/node"
)

type MessageRec struct {
	From       node.Node
	Received   int
	MaxHop     int
	LastTS     int64
	Timestamps map[int]int64
}

func NewPacketStatistic(from node.Node, timestamp int64) *MessageRec {
	return &MessageRec{
		from,
		0,
		0,
		0,
		map[int]int64{0: timestamp},
	}
}

func (ps *MessageRec) Delay(last int) int {
	return int(ps.Timestamps[last] - ps.Timestamps[0])
}

func (ps *MessageRec) Duration() int {
	return int(ps.LastTS)
}
