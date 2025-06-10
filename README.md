# Ditto

A command line tool for managing Melon infrastructure

### Project file

Ditto can handle sub-commands specified from a project-to-project basis. You can define a [project.ditto](./project.ditto) file like this:
```JSON
{
    "name": "Ditto",
    "tasks": {
        "hi": {
            "name": "hi",
            "jobs": [
                {
                    "name": "Echo hi!",
                    "run": "echo 'hi'"
                }
            ]
        }
    }
}
```

It's basically a JSON file that you can specify workflows within. You can run the 'hi' subcommand by doing `ditto hi`