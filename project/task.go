package project

import "log"

type Task struct {
	Name string `json:"name"`
	Jobs Jobs   `json:"jobs"`
}

func (t *Task) Run() {
	log.Printf("Running task: %s", t.Name)
	t.Jobs.RunAllJobs()
}
