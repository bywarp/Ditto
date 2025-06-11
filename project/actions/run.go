package actions

import (
	"bufio"
	"fmt"
	"os/exec"
	"syscall"
)

// Action to run a Bash command
// - command: the command to run
var Run = Action{
	Inputs: []string{"command"},
	Function: func(inputs map[string]string) error {
		cmd := exec.Command("bash", "-c", inputs["command"])
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
				fmt.Printf(" | " + str)
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
	},
}
