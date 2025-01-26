package main

import (
	"log"

	"github.com/net0pyr/nutsFM/captures"

	"github.com/rivo/tview"
)

// createSearchApp создает и запускает приложение поиска.
func createSearchApp(searchTreeView *tview.TreeView, rootDir string) {
	searchApp := tview.NewApplication()
	captures.SearchApp = searchApp
	captures.RootDir = &rootDir

	searchTreeView.SetInputCapture(captures.SearchAppCaptures)
	captures.SearchTreeView = searchTreeView

	if err := searchApp.SetRoot(searchTreeView, true).Run(); err != nil {
		log.Fatal(err)
	}
}
