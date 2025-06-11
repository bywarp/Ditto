# Ditto

A command line tool for managing Melon infrastructure

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