package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

// run single script
func runScript(name string) (err error) {
	scriptCmd := exec.Command(name)
	var scriptStdErr bytes.Buffer
	scriptCmd.Stderr = &scriptStdErr
	output, err := scriptCmd.Output()
	fmt.Println(fmt.Sprintf("----- %s -----\n", name), string(output), fmt.Sprintf("\n--- end of output (%s) ---", name))
	if err != nil {
		fmt.Println(fmt.Sprintf("Stderr of %s: ", name), scriptStdErr.String())
		log.Println("cannot execute script: ", err)
		// log.Println("executed command: ", scriptCmd.String())
		return err
	}
	return nil
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
