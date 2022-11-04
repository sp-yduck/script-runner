package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
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
