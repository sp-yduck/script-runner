package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

func copyFile(srcPath string, dstPath string) {
	src, err := os.Open(srcPath)
	if err != nil {
		log.Println("cannot open file: ", err)
	}
	defer src.Close()
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Println("cannot create new file: ", err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		log.Println("cannot copy file: ", err)
	}
}

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
