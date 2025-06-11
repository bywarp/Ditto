package actions

import (
	"errors"
	"strings"
)

type Action struct {
	Inputs   []string
	Function func(map[string]string) error
}

func (action *Action) Run(inputs map[string]string) error {
	for _, input_name := range action.Inputs {
		_, found := inputs[input_name]
		if !found {
			return errors.New("'" + input_name + "' is required for action to run")
		}
	}

	return action.Function(inputs)
}

type Actions map[string]Action

func (actions Actions) Run(name string, inputs map[string]string) error {
	items := strings.Split(name, "/")
	if len(items) < 2 ||
		items[0] != "@ditto" {
		return errors.New("'" + name + "' is not a valid action")
	}

	action_name := items[1]
	action, found := actions[items[1]]
	if !found {
		return errors.New("'" + action_name + "' is not a valid action name")
	}

	return action.Run(inputs)
}

var ActionList Actions = map[string]Action{
	"check_go_install": CheckGoInstall,
	"run":              Run,
	"write":            Write,
}
