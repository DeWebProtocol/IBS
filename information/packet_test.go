package information

import (
	"fmt"
	"testing"

	"github.com/dewebprotocol/IBS/node"
)

func TestBasicPacket(t *testing.T) {
	n1 := node.NewBasicNode(1, 0, 0, "a", nil)
	n2 := node.NewBasicNode(2, 0, 0, "b", nil)
	n3 := node.NewBasicNode(3, 0, 0, "c", nil)
	bp := NewBasicPacket(0, 0, n1, n1, n2, nil, 0)
	bp2 := *bp
	bp2.relayNode = n2
	bp3 := *bp
	bp3.relayNode = n3

	n2.SetTsLastSending(1234)

	fmt.Println("bp", bp)
	fmt.Println("bp2", bp2.relayNode.TsLastSending())
	fmt.Println("bp3", bp3)
}
