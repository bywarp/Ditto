package job

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"syscall"
)

type Jobs []Job
type Job struct {
	Name   string                  `json:"name"`
	Action func(*JobContext) error `json:"-"`
}

func (j Jobs) RunAllJobs() error {
	context := JobContext{}
	for _, job := range j {
		log.Printf("%s", job.Name)
		if err := job.Action(&context); err != nil {
			return err
		}
	}
	return nil
}

type JobContext struct {
	Progress int `json:"progress"`
}

func (j *JobContext) Log(message string) {
	fmt.Println(" | " + message)
}

func (j *JobContext) LogNoNewline(message string) {
	fmt.Print(" | " + message)
}

func (j *JobContext) Command(command string) error {
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	defer stdout.Close()
	stdoutReader := bufio.NewReader(stdout)

	go (func() {
		for {
			str, err := stdoutReader.ReadString('\n')
			if err != nil {
				break
			}
			j.LogNoNewline(str)
		}
	})()

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if _, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return nil
			}
		}
		return err
	}

	return nil
}
