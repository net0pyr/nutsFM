package main

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

func createMainApp(rootDir string) {
	app := tview.NewApplication()
	captures.App = app

	var rootPath string
	captures.RootPath = &rootPath
	var err error

	if rootDir == "" {
		rootPath, err = filepath.Abs(".")
		if err != nil {
			log.Fatalf("Failed to get absolute path: %v", err)
		}
	} else {
		rootPath, err = filepath.Abs(rootDir)
		if err != nil {
			log.Fatalf("Failed to get absolute path: %v", err)
		}

		info, err := os.Stat(rootPath)
		if os.IsNotExist(err) {
			fmt.Printf("Directory does not exist: %s\n", rootPath)
			return
		}
		if !info.IsDir() {
			fmt.Printf("Not a directory: %s\n", rootPath)
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

	nodes.ResetNodeStyles()

	captures.SortField.SetChangedFunc(func(text string) {
		nodes.Sort(text)
	})
	captures.SortField.SetChangedFunc(func(text string) {
		nodes.Sort(text)
	})

	app.SetInputCapture(captures.AppCaptures)

	if err := app.SetRoot(treeView, true).Run(); err != nil {
		log.Fatal(err)
	}
}
