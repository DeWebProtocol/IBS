package flags

import "github.com/urfave/cli/v2"

var WithNE = &cli.BoolFlag{
	Name:  "with_ne",
	Usage: "using NE mechanism",
}

var MaliciousNode = &cli.BoolFlag{
	Name: "malicious",
	Usage: "true: half of nodes refuse to relay messages\n" +
		"false: half of nodes disconnect",
}

var OutputPacket = &cli.BoolFlag{
	Name: "output_packet",
}

var NeLearn = &cli.BoolFlag{
	Name: "ne_learn",
}
