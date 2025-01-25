package main

import (
	"log"
	"os"

	"github.com/dewebprotocol/IBS/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:     "IBS",
		Commands: cmd.Root,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
