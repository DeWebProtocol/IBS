package cmd

import "github.com/urfave/cli/v2"

var Root = []*cli.Command{
	Kademlia,
	Flood,
	Version,
}
