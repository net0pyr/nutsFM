package fstree

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Create(node *tview.TreeNode, path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Failed to read directory %s: %v", path, err)
		return
	}

	parentNode := tview.NewTreeNode("..").
		SetReference(filepath.Join(path, ".."))
	parentNode.SetColor(tcell.ColorPink)
	node.AddChild(parentNode)

	for _, file := range files {
		childPath := filepath.Join(path, file.Name())
		childNode := tview.NewTreeNode(file.Name()).
			SetReference(childPath)

		info, err := os.Stat(childPath)
		if err != nil {
			log.Printf("Failed to stat %s: %v", childPath, err)
			continue
		}

		if info.IsDir() {
			childNode.SetColor(tcell.ColorPink)
		} else {
			childNode.SetColor(tcell.ColorGreen)
		}

		childNode.Collapse()
		node.AddChild(childNode)
	}
}
