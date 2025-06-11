package commands

import (
	"context"
	"log"
	"strconv"

	"github.com/urfave/cli/v3"
	"wrp.sh/ditto/project"
)

type List struct{}

func (cmd List) Command() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"l", "ls"},
		Usage:   "List all jobs in the project file",
		Action:  cmd.Action,
	}
}

func (List) Action(ctx context.Context, cmd *cli.Command) error {
	project, err := project.ReadProjectFile()
	if err != nil {
		return err
	}

	jobs := project.Jobs
	log.Println(project.Name + " - " + strconv.Itoa(len(jobs)) + " job(s):")
	for name, job := range jobs {
		log.Println(" - " + name + ": " + job.Description)
	}

	return nil
}
