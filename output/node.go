package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/dewebprotocol/IBS/node"
)

type Node struct {
	Id     string `json:"id"`
	Region string `json:"region"`
	//DownloadBandwidth int    `json:"downloadBandwidth"` // byte/s
	UploadBandwidth int  `json:"uploadBandwidth"`
	Running         bool `json:"running"`
	CrashFactor     int  `json:"crashFactor"`
	CrashTimes      int  `json:"crashTimes"`
}

func newNode(n node.Node) *Node {
	return &Node{
		strconv.FormatUint(n.Id(), 10),
		n.Region(),
		//n.DownloadBandwidth(),
		n.UploadBandwidth(),
		n.Running(),
		n.CrashFactor(),
		n.CrashTimes(),
	}
}

type NodeOutput []*Node

func NewNodeOutput() NodeOutput {
	return NodeOutput{}
}

func (o *NodeOutput) Append(n node.Node) {
	*o = append(*o, newNode(n))
}

func (o *NodeOutput) WriteNodes(folder string) {
	b, err := json.Marshal(o)
	if err != nil {
		fmt.Println(err)
	}
	filename := fmt.Sprintf("%s/output_nodes.json", folder)
	err = os.WriteFile(filename, b, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
