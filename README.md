# Ditto

A command line tool for managing Melon infrastructure

### Commands

- `ditto init` - will create a ditto.project file in the current directory
- `ditto list` - will list all of the jobs that can be executed with Ditto
- `ditto <job>` - will run a job specified in the project.ditto file

### Project file

Ditto can handle sub-commands specified from a project-to-project basis. You can define a [project.ditto](./project.ditto) file like this:
```JSON
{
    "name": "Ditto",
    "jobs": {
        "build": {
            "description": "Build Ditto executable",
            "tasks": [
                {
                    "action": "@ditto/run",
                    "inputs": {
                        "command": "go build ."
                    }
                }
            ]
        }
    }
}
```

It's basically a JSON file that you can specify workflows within. You can run the 'build' subcommand by doing `ditto build`

### Actions

Ditto currently has a few supported actions:

| Action | Description | Inputs |  
| - | - | - |
| `@ditto/run` | run a command | command
| `@ditto/write` | write to a file | file, data
| `@ditt/check_go_install` | check if a go program is installed | name

There will be more added in the future, this is just currently what our projects need