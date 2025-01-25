package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

const version = "v0.1.2"

var Version = &cli.Command{
	Name: "version",
	Action: func(ctx *cli.Context) error {
		fmt.Println("IBS version:", version)
		return nil
	},
}
