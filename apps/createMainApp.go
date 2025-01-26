package apps

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/net0pyr/nutsFM/captures"
	"github.com/net0pyr/nutsFM/fstree"
	"github.com/net0pyr/nutsFM/nodes"

	"github.com/rivo/tview"
)

// createMainApp создает и запускает основное TUI-приложение.
func CreateMainApp(rootDir string) {
	app := tview.NewApplication()
	captures.App = app

	var rootPath string
	captures.RootPath = &rootPath
	var err error

	if rootDir == "" {
		rootPath, err = filepath.Abs(".")
		if err != nil {
			log.Fatalf("Не удалось получить абсолютный путь: %v", err)
		}
	} else {
		rootPath, err = filepath.Abs(rootDir)
		if err != nil {
			log.Fatalf("Не удалось получить абсолютный путь: %v", err)
		}

		info, err := os.Stat(rootPath)
		if os.IsNotExist(err) {
			fmt.Printf("Директория не существует: %s\n", rootPath)
			return
		}
		if !info.IsDir() {
			fmt.Printf("Это не директория: %s\n", rootPath)
			return
		}
	}

	rootNode := tview.NewTreeNode(rootPath).
		SetReference(rootPath)

	fstree.Create(rootNode, rootPath)

	treeView := tview.NewTreeView().
		SetRoot(rootNode).
		SetCurrentNode(rootNode)
	captures.TreeView = treeView
	nodes.TreeView = treeView

	treeView.SetInputCapture(captures.ViewCaptures)

	Pages.AddPage("main", treeView, true, true)

	nodes.ResetNodeStyles()

	captures.SortField.SetChangedFunc(func(text string) {
		nodes.Sort(text)
	})
	captures.SortField.SetChangedFunc(func(text string) {
		nodes.Sort(text)
	})

	app.SetInputCapture(captures.AppCaptures)

	if err := app.SetRoot(Pages, true).Run(); err != nil {
		log.Fatal(err)
	}
}
