package actions

import (
	"fmt"
	"os"
)

// Action to write data to a file
// - data: the data to write
// - file: the file to write to
var Write = Action{
	Inputs: []string{"data", "file"},
	Function: func(inputs map[string]string) error {

		data := []byte(inputs["data"])
		file := inputs["file"]

		fmt.Println(" | Writing to file " + file)
		if err := os.WriteFile(file, data, 0644); err != nil {
			return err
		}

		return nil
	},
}
