package actions

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strings"

	"github.com/fatih/color"
	"wrp.sh/ditto/utils"
)

// Action to check if a Go program is installed
// - name: name of the program
var CheckGoInstall = Action{
	Inputs: []string{"name"},
	Function: func(inputs map[string]string) error {
		go_install_dir, found := os.LookupEnv("GOPATH")
		if !found {
			return errors.New("failed to find GOPATH in environment")
		}
		fmt.Println(" | Found GOPATH: " + go_install_dir)

		installs, err := os.ReadDir(go_install_dir + "/bin")
		if err != nil {
			return err
		}

		name := inputs["name"]
		has_exec := slices.ContainsFunc(installs, func(install fs.DirEntry) bool {
			return strings.Contains(install.Name(), inputs["name"])
		})
		fmt.Println(" | Has " + name + ": " + utils.Ternary(has_exec, color.GreenString("yes"), color.RedString("no")))

		if !has_exec {
			return errors.New("'" + name + "' was not found!")
		}
		return nil
	},
}
