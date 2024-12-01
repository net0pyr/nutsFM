package fstree

import (
	"log"
	"os"
	"path/filepath"

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
	parentNode.SetColor(tview.Styles.SecondaryTextColor)
	node.AddChild(parentNode)

	for _, file := range files {
		childNode := tview.NewTreeNode(file.Name()).
			SetReference(filepath.Join(path, file.Name()))

		if file.IsDir() {
			childNode.SetColor(tview.Styles.PrimaryTextColor)
			childNode.Collapse()
		} else {
			childNode.SetColor(tview.Styles.SecondaryTextColor)
		}

		node.AddChild(childNode)
	}
}
