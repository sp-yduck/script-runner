package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	files, err := os.ReadDir("./pipelines")
	if err != nil {
		log.Fatal("cannot read directory './pipelines': ")
	}

	// read pipelines
	var pipelines []Pipeline
	for _, f := range files {
		log.Println(f.Name())
		relScriptPath := filepath.Join("./pipelines", f.Name())
		absScriptPath, err := filepath.Abs(relScriptPath)
		if err != nil {
			log.Fatal("cannot get absolution path: ", err)
		}

		pipeline := readPipeline(absScriptPath)
		pipelines = append(pipelines, pipeline)
	}

	// run pipelines
	for _, p := range pipelines {
		fmt.Printf("----- pipeline | %s -----\n", p.Name)
		for i, task := range p.Tasks {
			err = runTask(task)
			if err != nil {
				log.Println("task exit with error: ", err)
				fmt.Printf("remaining tasks: %v\n", p.Tasks[i:])
				break
			}
		}
	}
}
