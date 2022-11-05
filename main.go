package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var defaultTimeout int64
var defaultLogdir string

func init() {
	defaultTimeout = 5
	defaultLogdir = "./log"
}

func main() {
	logInit(defaultLogdir)

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
	// to do: pipelineごとに出力ファイル分ける
	// stdout,stderrが複数のpipelineのものでごっちゃにならないようにする
	ch := make(chan error, len(pipelines))
	defer close(ch)
	for _, p := range pipelines {
		go func(p Pipeline) {
			ch <- runPipeline(p)
		}(p)
	}

	for i := 0; i < len(pipelines); i++ {
		fmt.Println(<-ch)
	}
}
