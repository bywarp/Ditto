package commands

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"wrp.sh/ditto/project"
	"wrp.sh/ditto/utils"
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

	var jobs project.Jobs = []project.Job{
		{
			Name: "Check Java installation",
			Action: func(context *project.JobContext) error {
				return context.Command("java --version")
			},
		},
		{
			Name: "Check Go installation",
			Action: func(context *project.JobContext) error {
				return context.Command("go version")
			},
		},
		{
			Name: "Check git installation",
			Action: func(context *project.JobContext) error {
				return context.Command("git version")
			},
		},
		{
			Name: "Check protobufs installation",
			Action: func(context *project.JobContext) error {
				go_install_dir, found := os.LookupEnv("GOPATH")
				if !found {
					return errors.New("Failed to find GOPATH in environment.")
				}
				context.Log("Found GOPATH: " + go_install_dir)

				installs, err := os.ReadDir(go_install_dir + "/bin")
				if err != nil {
					return err
				}

				has_protoc := slices.ContainsFunc(installs, func(install fs.DirEntry) bool {
					return strings.Contains(install.Name(), "protoc")
				})
				context.Log("Has protoc: " + utils.Ternary(has_protoc, "yes", "no"))

				has_protoc_gen := slices.ContainsFunc(installs, func(install fs.DirEntry) bool {
					return strings.Contains(install.Name(), "protoc-gen-go")
				})
				context.Log("Has protoc_gen: " + utils.Ternary(has_protoc, "yes", "no"))

				if !has_protoc || !has_protoc_gen {
					return errors.New("You're missing protoc-gen-go or protoc. Did you install them correctly?")
				}
				return nil
			},
		},
	}

	if err := jobs.RunAllJobs(); err != nil {
		return err
	}

	log.Println(color.GreenString("Done!"))
	return nil
}
