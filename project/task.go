package project

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"syscall"
)

type Tasks []Task
type Task struct {
	Name             string                   `json:"name"`
	Run              string                   `json:"run"`
	WorkingDirectory string                   `json:"working_directory"`
	Action           func(*TaskContext) error `json:"-"`
}

func (t Tasks) RunAllJobs() error {
	context := TaskContext{}
	for _, task := range t {
		log.Printf("%s", task.Name)
		if task.Action != nil {
			if err := task.Action(&context); err != nil {
				return err
			}
		} else if task.Run != "" {
			if err := context.Command(task.Run); err != nil {
				return err
			}
		} else {
			return errors.New("Job must have a 'run' specified")
		}
	}
	return nil
}

type TaskContext struct {
	Progress int `json:"progress"`
}

func (t *TaskContext) Log(message string) {
	fmt.Println(" | " + message)
}

func (t *TaskContext) LogNoNewline(message string) {
	fmt.Print(" | " + message)
}

func (t *TaskContext) Command(command string) error {
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
			t.LogNoNewline(str)
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
