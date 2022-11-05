package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func logInit(logdir string) {
	// prepare log directory
	if err := os.Mkdir(logdir, 0666); err != nil {
		log.Println("cannot make new dir: ", err)
	}

	// configure log file name format
	now := time.Now()
	path := filepath.Join(logdir, now.Format("20060102-150405")+".log")

	// prepare log file
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("cannot open file: ", err)
	}
	logfile := io.MultiWriter(os.Stdout, file)

	// configure log
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(logfile)
}
