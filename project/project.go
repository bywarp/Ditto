package project

import (
	"encoding/json"
	"os"
)

const PROJECT_FILE_NAME = "project.ditto"

type Project struct {
	Name string         `json:"name"`
	Jobs map[string]Job `json:"jobs"`
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
