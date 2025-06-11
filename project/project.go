package project

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
)

const PROJECT_FILE_NAME = "project.ditto"

type Project struct {
	Name string         `json:"name"`
	Jobs map[string]Job `json:"jobs"`
}

func (project *Project) Run(name string) error {
	job, found := project.Jobs[name]
	if !found {
		return errors.New("there was no job named '" + name + "' in the project.ditto file")
	}

	depends_on := job.DependsOn
	if len(depends_on) > 0 {
		log.Println(name + " depends on: " + strings.Join(depends_on, ", "))
		for _, depends := range depends_on {
			project.Run(depends)
		}
	}

	if err := job.Run(); err != nil {
		return err
	}

	return nil
}

func ReadProjectFile() (*Project, error) {
	content, err := os.ReadFile(PROJECT_FILE_NAME)
	if err != nil {
		return nil, err
	}

	project := Project{}
	if err := json.Unmarshal(content, &project); err != nil {
		return nil, err
	}

	return &project, nil
}
