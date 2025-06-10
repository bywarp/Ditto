package project

type Job struct {
	Description string `json:"description"`
	DependsOn   string `json:"depends_on"`
	Tasks       Tasks  `json:"tasks"`
}

func (j *Job) Run() {
	j.Tasks.RunAllJobs()
}
