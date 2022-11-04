package main

import (
	"log"
	"os"
	"path/filepath"
)

var defaultTimeout int64

func init() {
	defaultTimeout = 5
}

func main() {
	baseDir := "./pipelines"
	files, err := os.ReadDir(baseDir)
	if err != nil {
		log.Fatalf("cannot read directory '%s': ", baseDir)
	}

	// read pipelines
	var pipelines []Pipeline
	for _, f := range files {
		log.Println(f.Name())
		relScriptPath := filepath.Join(baseDir, f.Name())
		absScriptPath, err := filepath.Abs(relScriptPath)
		if err != nil {
			log.Fatal("cannot get absolution path: ", err)
		}
		pipeline := readPipelines(absScriptPath)
		pipelines = append(pipelines, pipeline...)
	}

	// run pipelines
	for _, p := range pipelines {
		runPipeline(p)
	}
}
