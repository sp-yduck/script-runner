# ScriptRunner
ScriptRunner is a golang tool for running a tons of scripts easily/simple/fast.

## how to use
### define your Tasks

A Task consists of a name, commands and export_output.

```sample Task
# Task has a commands

name: name of Task
command: which kubectl
export_output: KUBECTL
```

### define your Pipelines

A Pipeline consists of a series of Tasks.
If you register a output with `export_output` field, you can use that output as env variable in only same Pipeline.
If you are going to run multiple Pipelines, your Pipelines need to be independent from other Pipelines.
<!-- so that ScriptRunner can run multiple Pipelines in parallel -->

```sample Pipeline
# Pipeline has a name and series of Tasks

name: name of Pipeline
tasks:
  - name: name of Task1
    command: which kubectl
    export_output: KUBECTL

  - name: name of Task2
    command: echo $KUBECTL

  - name: name of Task3
    command: echo 'hello world'
```

### execute ScriptRunner
```
go run main.go
```

## Roadmap
- [x] run multiple pipelines in parallel
- [ ] save output in files & reduce std output
- [x] add timeout feature
- [ ] be able to configure running parameters by config file (e.g. log level, timeout)
- [x] output log to files
- [ ] enrich log information (e.g. struct, log level)
- [ ] add more feature to Tasks/Pipelines (e.g. ignore_error)
- [x] register output to variables
- [ ] become a CLI