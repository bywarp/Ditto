package main

import (
	"context"
	"log"
	"os"

	"wrp.sh/ditto/commands"
	"wrp.sh/ditto/project"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

func main() {

	log.SetPrefix(color.CyanString("Ditto: "))
	log.SetFlags(log.Lmsgprefix)

	cmd := &cli.Command{
		Name:  "ditto",
		Usage: "A command-line tool for Melon development",
		Commands: []*cli.Command{
			commands.Init{}.Command(),
			commands.Dev{}.Command(),
			commands.Build{}.Command(),
		},
		CommandNotFound: func(ctx context.Context, c *cli.Command, s string) {
			project, err := project.ReadProjectFile()
			if err != nil {
				log.Println(color.RedString("project.ditto could not be found!"))
				return
			}

			job, found := project.Jobs[s]
			if !found {
				log.Printf(color.RedString("There was no job named '%s' in the Ditto file!"), s)
				return
			}

			job.Run()
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Println(color.RedString(err.Error()))
	}
}
