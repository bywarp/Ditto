package project

import (
	"log"
)

type Job struct {
	Description string   `json:"description"`
	DependsOn   []string `json:"depends_on"`
	Tasks       Tasks    `json:"tasks"`
}

func (j *Job) Run() error {
	log.Println(j.Description)
	if err := j.Tasks.Run(); err != nil {
		return err
	}
	return nil
}
