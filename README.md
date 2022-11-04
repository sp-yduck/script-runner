# ScriptRunner
ScriptRunner is a golang tool for running a tons of scripts easily/simple/fast.

## how to use
### define your Tasks

A Task consists of a series of commands.

```sample Task
# Task has a name and series of commands

name: name of Task
command:
- which kubectl
- kubectl version
```

### define your Pipelines

A Pipeline consists of a series of Tasks.
If you are going to run multiple Pipelines, your Pipelines need to be independent from other Pipelines.
<!-- so that ScriptRunner can run multiple Pipelines in parallel -->

```sample Pipeline
# Pipeline has a name and series of Tasks

name: name of Pipeline
tasks:
  - name: name of Task1
    command:
    - which kubectl
    - kubectl version

  - name: name of Task2
    command:
    - which helm
    - helm version

  - name: name of Task3
    command:
    - echo 'hello world'
```

### execute ScriptRunner
```
go run main.go
```