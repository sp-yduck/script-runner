package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Pipeline struct {
	Name  string `yaml:"name"`
	Tasks []Task `yaml:"tasks"`
}

type Task struct {
	Name         string `yaml:"name"`
	Command      string `yaml:"command"`
	ExportOutput string `yaml:"export_output,omitempty"`
}

// unmarshal pipeline object from filepath
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

// run single pipeline
func runPipeline(p Pipeline) (err error) {
	fmt.Printf("----- pipeline | %s -----", p.Name)
	variables := os.Environ()
	for i, task := range p.Tasks {
		fmt.Println(fmt.Sprintf("\n----- task | %s -----", task.Name))
		scriptCmd := exec.Command("sh", "-c", task.Command)
		scriptCmd.Env = variables
		var scriptStdErr bytes.Buffer
		scriptCmd.Stderr = &scriptStdErr

		output, err := scriptCmd.Output()
		if task.ExportOutput != "" {
			variables = append(variables, fmt.Sprintf("%s=%s", task.ExportOutput, string(output)))
		}

		fmt.Print(fmt.Sprintf("    command: %s\n    output: ", task.Command), string(output))
		if err != nil {
			fmt.Print("    stderr: ", scriptStdErr.String())
			fmt.Println("    err: ", err)
			log.Println("executed command: ", scriptCmd.String())
			log.Println("task exit with error: ", err)
			log.Printf("remaining tasks: %v\n", filepath.Join(getTasksName(p.Tasks[i:])...))
			return err
		}
	}
	return nil
}

// get all of the tasks name inputed as slice
func getTasksName(tasks []Task) (names []string) {
	for _, t := range tasks {
		names = append(names, t.Name)
	}
	return names
}
