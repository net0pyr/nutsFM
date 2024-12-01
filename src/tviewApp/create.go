package tviewApp

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/net0pyr/nutsFM/tviewApp/fstree"
	"github.com/rivo/tview"
)

func Create(rootDir string) {
	app := tview.NewApplication()

	var rootPath string
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

	treeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlT {
			node := treeView.GetCurrentNode()
			if node == nil {
				return event
			}

			ref := node.GetReference()
			if ref == nil {
				return event
			}

			path := ref.(string)
			info, err := os.Stat(path)
			if err != nil {
				log.Printf("Failed to access %s: %v", path, err)
				return event
			}

			if info.IsDir() {
				if node.IsExpanded() {
					node.Collapse()
				} else {
					node.ClearChildren()
					fstree.Create(node, path)
					node.Expand()
				}
			}
			return nil
		}
		return event
	})

	treeView.SetSelectedFunc(func(node *tview.TreeNode) {
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
			newRootNode := tview.NewTreeNode(path).
				SetReference(path)
			fstree.Create(newRootNode, path)
			treeView.SetRoot(newRootNode).
				SetCurrentNode(newRootNode)
		} else {
			cmd := exec.Command("xdg-open", path)
			err := cmd.Start()
			if err != nil {
				log.Printf("Failed to open file %s: %v", path, err)
			}
		}
	})

	if err := app.SetRoot(treeView, true).Run(); err != nil {
		log.Fatal(err)
	}
}
