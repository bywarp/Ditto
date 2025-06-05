package commands

import (
	"context"
	"log"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

type Dev struct{}

func (cmd Dev) Command() *cli.Command {
	return &cli.Command{
		Name:    "dev",
		Aliases: []string{"d"},
		Usage:   "Manage a local Melon deployment",
		Action:  cmd.Action,
	}
}

func (Dev) Action(ctx context.Context, cmd *cli.Command) error {
	log.Println(color.RedString("This command is not yet implemented!"))
	return nil
}
