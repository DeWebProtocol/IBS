package output

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type info struct {
	Id     int    `json:"id"`
	Delay  int    `json:"delay"`
	Region string `json:"region"`
}

type LatencyOutput []*info // packetID -> received

func NewLatencyOutput() LatencyOutput {
	return LatencyOutput{}
}
func (o *LatencyOutput) Append(id, delay int, region string) {
	*o = append(*o, &info{id, delay, region})
}

func (o *LatencyOutput) WriteLatency(folder string) {
	b, err := json.Marshal(o)
	if err != nil {
		fmt.Println(err)
	}
	fileName := fmt.Sprintf("%s/output_delay.json", folder)
	err = ioutil.WriteFile(fileName, b, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
