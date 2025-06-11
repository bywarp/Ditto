package project

import (
	"log"

	"wrp.sh/ditto/project/actions"
)

type Tasks []Task
type Task struct {
	Description string            `json:"description"`
	Action      string            `json:"action"`
	Inputs      map[string]string `json:"inputs"`
}

func (t Tasks) Run() error {
	for _, task := range t {
		if task.Description != "" {
			log.Println(task.Description)
		}
		if err := actions.ActionList.Run(task.Action, task.Inputs); err != nil {
			return err
		}
	}
	return nil
}
