package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type Pipeline struct {
	Name  string `yaml:"name"`
	Tasks []Task `yaml:"tasks"`
}

type Task struct {
	Name    string   `yaml:"name"`
	Command []string `yaml:"command"`
}

func readPipeline(path string) (pipeline Pipeline) {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("cannot read file: ", err)
	}
	err = yaml.Unmarshal(b, &pipeline)
	if err != nil {
		log.Fatal("cannot unmarshal yaml: ", err)
	}
	return pipeline
}

func readTask(path string) (task Task) {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(b, &task)
	if err != nil {
		log.Fatal(err)
	}
	return task
}

// run single task
func runTask(task Task) (err error) {
	fmt.Println(fmt.Sprintf("----- task | %s -----", task.Name))
	for i, cmd := range task.Command {
		scriptCmd := exec.Command("sh", "-c", cmd)
		var scriptStdErr bytes.Buffer
		scriptCmd.Stderr = &scriptStdErr
		output, err := scriptCmd.Output()
		fmt.Println(fmt.Sprintf("%s\n    output: ", cmd), string(output))
		if err != nil {
			fmt.Println("    stderr: ", scriptStdErr.String())
			fmt.Println("    err: ", err)
			// log.Println("executed command: ", scriptCmd.String())
			fmt.Printf("remaining command: %s\n", task.Command[i:])
			return err
		}
	}
	fmt.Printf("task (%s) completed\n", task.Name)
	return nil
}
