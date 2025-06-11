package commands

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"wrp.sh/ditto/project"
)

type Init struct{}

func (cmd Init) Command() *cli.Command {
	return &cli.Command{
		Name:      "init",
		Aliases:   []string{"i"},
		ArgsUsage: "The name of the project",
		Usage:     "Initialize a Ditto project",
		Action:    cmd.Action,
	}
}

func (Init) Action(ctx context.Context, cmd *cli.Command) error {

	name := cmd.Args().Get(0)

	if name == "" {
		return errors.New("please specify a project name")
	}

	new_project := project.Project{
		Name: name,
		Jobs: map[string]project.Job{
			"test": {
				Description: "Example job",
				Tasks: project.Tasks{
					{
						Action:      "@ditto/run",
						Description: "Echo 'Hello, World' to the console",
						Inputs: map[string]string{
							"command": "echo \"Hello, World!\"",
						},
					},
				},
			},
		},
	}

	data, err := json.MarshalIndent(new_project, "", "  ")
	if err != nil {
		return err
	}

	if _, err := os.Stat(project.PROJECT_FILE_NAME); !os.IsNotExist(err) {
		return errors.New("project file already exists in this directory")
	}

	if err := os.WriteFile(project.PROJECT_FILE_NAME, data, 0666); err != nil {
		return err
	}

	log.Println(color.GreenString("Created " + project.PROJECT_FILE_NAME + " file in directory!"))
	return nil
}
