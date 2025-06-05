package commands

import (
	"context"
	"log"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

type Init struct{}

func (cmd Init) Command() *cli.Command {
	return &cli.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "Initialize a Ditto project",
		Action:  cmd.Action,
	}
}

func (Init) Action(ctx context.Context, cmd *cli.Command) error {
	log.Println("Initializing Melon project..")
	log.Println(color.GreenString("Done!"))
	return nil
}
