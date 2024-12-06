package tviewApp

import (
	"log"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createSearchApp(searchTreeView *tview.TreeView, rootDir string) {
	searchApp := tview.NewApplication()

	searchTreeView.SetSelectedFunc(func(node *tview.TreeNode) {
		log.Printf("Selected node: %s", node.GetText())
		ref := node.GetReference()
		if ref == nil {
			return
		}

		path := ref.(string)
		info, err := os.Stat(path)
		if err != nil {
			log.Printf("Failed to access %s: %v", path, err)
			return
		}

		if info.IsDir() {
			searchApp.Stop()
			CreateMainApp(path)
		} else {
			cmd := exec.Command("xdg-open", path)
			err := cmd.Start()
			if err != nil {
				log.Printf("Failed to open file %s: %v", path, err)
			}
		}
	})

	searchTreeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		log.Printf("Event: %v", event)
		if event.Key() == tcell.KeyEsc {
			searchApp.Stop()
			CreateMainApp(rootDir)
			return nil
		}
		return event
	})

	if err := searchApp.SetRoot(searchTreeView, true).Run(); err != nil {
		log.Fatal(err)
	}

}
