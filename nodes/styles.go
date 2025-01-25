package nodes

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SetDefaultNodeStyles() {
	for _, node := range TreeView.GetRoot().GetChildren() {
		ref := node.GetReference()
		if ref != nil {
			path := ref.(string)
			info, err := os.Stat(path)
			if err == nil {
				if info.IsDir() {
					node.SetColor(tcell.ColorPink)
				} else {
					node.SetColor(tcell.ColorGreen)
				}
			}
		}
	}
}

func UpdateNodeStyles() {
	SetDefaultNodeStyles()
	for _, node := range SortResults {
		ref := node.GetReference()
		if ref != nil {
			path := ref.(string)
			info, err := os.Stat(path)
			if err == nil {
				if info.IsDir() {
					node.SetColor(tcell.ColorRed)
				} else {
					node.SetColor(tcell.ColorSpringGreen)
				}
			}
		}
	}
	if len(SortResults) > 0 {
		SortResults[CurrentIndex].SetColor(tcell.ColorRed)
	}
}

func ResetNodeStyles() {
	TreeView.GetRoot().Walk(func(node, parent *tview.TreeNode) bool {
		ref := node.GetReference()
		if ref != nil {
			path := ref.(string)
			info, err := os.Stat(path)
			if err == nil {
				if info.IsDir() {
					node.SetColor(tcell.ColorPink)
				} else {
					node.SetColor(tcell.ColorGreen)
				}
			}
		}
		return true
	})
}
