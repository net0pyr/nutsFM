package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/net0pyr/nutsFM/tviewApp"
)

var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
}

func main() {
	defer logFile.Close()

	var rootPath string
	var err error

	if len(os.Args) > 1 {
		rootPath, err = filepath.Abs(os.Args[1])
		if err != nil {
			log.Fatalf("Failed to get absolute path: %v", err)
		}
	} else {
		rootPath, err = filepath.Abs(".")
		if err != nil {
			log.Fatalf("Failed to get absolute path: %v", err)
		}
	}

	tviewApp.Create(rootPath)
}
