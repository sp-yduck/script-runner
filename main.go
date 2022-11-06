package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var defaultTimeout int64
var defaultLogdir string
var defaultLogVerbosity int8

func init() {
	defaultTimeout = 5
	defaultLogdir = "./log"
	defaultLogVerbosity = 0
}

func main() {
	logInit(defaultLogdir, defaultLogVerbosity)

	baseDir := "./pipelines"
	files, err := os.ReadDir(baseDir)
	if err != nil {
		fmt.Printf("cannot read directory '%s': %v\n", baseDir, err)
		log.Fatalf("cannot read directory '%s': %v\n", baseDir, err)
	}

	// read pipelines
	var pipelines []Pipeline
	for _, f := range files {
		log.Println(f.Name())
		relScriptPath := filepath.Join(baseDir, f.Name())
		absScriptPath, err := filepath.Abs(relScriptPath)
		if err != nil {
			fmt.Println("cannot get absolution path: ", err)
			log.Fatal("cannot get absolution path: ", err)
		}
		pipeline := readPipelines(absScriptPath)
		pipelines = append(pipelines, pipeline...)
	}

	// run pipelines
	ch := make(chan error, len(pipelines))
	defer close(ch)
	for _, p := range pipelines {
		go func(p Pipeline) {
			ch <- p.Run()
		}(p)
	}

	// to do: get result in order for finished pipeline
	for _, p := range pipelines {
		if err := <-ch; err != nil {
			log.Println(fmt.Sprintf("pipeline(%s) failed: ", p.Name), err)
		}
	}
}
