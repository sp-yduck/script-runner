package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func logInit(logdir string, verbosity int8) {
	// prepare log directory
	if _, err := os.Stat(logdir); err != nil {
		// make new dir
		if err := os.Mkdir(logdir, 0666); err != nil {
			log.Println("cannot make new dir: ", err)
		}
	}

	// configure log file name format
	now := time.Now()
	path := filepath.Join(logdir, now.Format("20060102-150405")+".log")

	// prepare log file
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("cannot open file: ", err)
	}

	// configure log output files
	var logfiles io.Writer
	switch verbosity {
	case 0:
		logfiles = io.MultiWriter(file)
	default:
		logfiles = io.MultiWriter(os.Stdout, file)
	}

	// configure log
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(logfiles)
}
