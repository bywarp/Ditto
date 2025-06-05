package commands

import (
	"context"
	"log"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

type Build struct{}

func (cmd Build) Command() *cli.Command {
	return &cli.Command{
		Name:    "build",
		Aliases: []string{"b"},
		Usage:   "Build a production Melon package",
		Action:  cmd.Action,
	}
}

func (Build) Action(ctx context.Context, cmd *cli.Command) error {
	log.Println(color.RedString("This command is not yet implemented!"))
	return nil
}
