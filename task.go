package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"gopkg.in/yaml.v2"
)

type Pipeline struct {
	Name  string `yaml:"name"`
	Tasks []Task `yaml:"tasks"`
}

type Task struct {
	Name         string `yaml:"name"`
	Command      string `yaml:"command"`
	ExportOutput string `yaml:"export_output,omitempty"`
	Timeout      int64  `yaml:"timeout,omitempty"`
	Result       TaskResult
}

type TaskResult struct {
	Stdout string
	Stderr string
	Err    error
}

// unmarshal pipeline object from filepath
func readPipelines(path string) (pipelines []Pipeline) {
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("cannot read file: ", err)
		log.Fatal("cannot read file: ", err)
	}
	err = yaml.UnmarshalStrict(b, &pipelines)
	if err != nil {
		fmt.Println("cannot unmarshal yaml: ", err)
		log.Fatal("cannot unmarshal yaml: ", err)
	}
	return pipelines
}

// run single pipeline
func runPipeline(p Pipeline) (err error) {
	variables := os.Environ()
	for _, task := range p.Tasks {
		// prepare timeout
		var scriptCmd *exec.Cmd
		timeout := time.Duration(defaultTimeout)
		if task.Timeout != 0 {
			timeout = time.Duration(task.Timeout)
		}

		// set context&command
		ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
		defer cancel()
		scriptCmd = exec.Command("sh", "-c", task.Command)

		// set env var
		scriptCmd.Env = variables

		// prepare results
		var scriptStdErr bytes.Buffer
		scriptCmd.Stderr = &scriptStdErr
		var scriptStdOut bytes.Buffer
		scriptCmd.Stdout = &scriptStdOut

		// exec command
		ch := make(chan error)
		defer close(ch)
		go func() {
			ch <- scriptCmd.Run()
		}()
		select {
		case err = <-ch:
		case <-ctx.Done():
			err = fmt.Errorf("context deadline exeeded: timeout is set to %s", timeout*time.Second)
		}

		// output results
		if task.ExportOutput != "" {
			variables = append(variables, fmt.Sprintf("%s=%s", task.ExportOutput, scriptStdOut.String()))
		}
		task.Result.Stdout = scriptStdOut.String()
		task.Result.Stderr = scriptStdErr.String()
		task.Result.Err = err

		// print result
		// fmt.Println(concludeTask(task))
		if err != nil {
			fmt.Println(concludePipeline(p))
			return err
		}
	}
	fmt.Println(concludePipeline(p))
	return nil
}

// get all of the tasks name inputed as slice
func getTasksName(tasks []Task) (names []string) {
	for _, t := range tasks {
		names = append(names, t.Name)
	}
	return names
}

func getRemainingTasks(pipeline Pipeline) (tasks []Task) {
	for i, t := range pipeline.Tasks {
		if t.Result.Err != nil {
			return pipeline.Tasks[i:]
		}
	}
	return
}

func concludeTask(task Task) (summary string) {
	summary = fmt.Sprintf("----- task | %s -----\n", task.Name)
	summary += fmt.Sprintf("    command: %s\n    output: %s\n", task.Command, task.Result.Stdout)
	if task.Result.Err != nil {
		summary += fmt.Sprintf("    stderr: %s\n", task.Result.Stderr)
		summary += fmt.Sprintf("    err: %v\n", task.Result.Err)
		summary += fmt.Sprintf("executed command: sh -c %s\n", task.Command)
		// summary += fmt.Sprintf("task exit with error: %v\n", task.Result.Err)
		// summary += fmt.Sprintf("remaining tasks: %v\n", filepath.Join(getTasksName(task.Result.RemainingTasks)...))
	}
	return summary
}

func concludePipeline(pipeline Pipeline) (summary string) {
	summary = fmt.Sprintf("----- pipeline | %s -----\n", pipeline.Name)
	for _, task := range pipeline.Tasks {
		summary += concludeTask(task)
	}
	return summary
}
