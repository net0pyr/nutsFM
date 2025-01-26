// Пакет main является точкой входа для приложения nutsFM.
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/net0pyr/nutsFM/captures"
	"github.com/rivo/tview"
)

// logFile представляет файл журнала для записи логов.
var logFile *os.File

// pages представляет страницы приложения.
var pages = tview.NewPages()

// init инициализирует приложение, настраивая логирование и страницы.
func init() {
	captures.Pages = pages
	var err error
	logFile, err = os.OpenFile(os.Getenv("HOME")+"/.nutsfm/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
}

// main является точкой входа приложения. Он настраивает корневой путь и запускает основное приложение.
func main() {
	defer logFile.Close()

	var rootPath string
	var err error

	if len(os.Args) > 1 {
		rootPath, err = filepath.Abs(os.Args[1])
		if err != nil {
			log.Fatalf("Не удалось получить абсолютный путь: %v", err)
		}
	} else {
		rootPath, err = filepath.Abs(".")
		if err != nil {
			log.Fatalf("Не удалось получить абсолютный путь: %v", err)
		}
	}

	captures.CreateSearchApp = createSearchApp
	captures.CreateMainApp = createMainApp
	createMainApp(rootPath)
}
