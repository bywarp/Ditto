package main

import (
	"context"
	"log"
	"os"

	"wrp.sh/ditto/commands"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

func main() {

	log.SetPrefix(color.CyanString("Ditto: "))
	log.SetFlags(log.Lmsgprefix)

	init := commands.Init{}
	cmd := &cli.Command{
		Name:  "ditto",
		Usage: "A command-line tool for Melon development",
		Commands: []*cli.Command{
			init.Command(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
