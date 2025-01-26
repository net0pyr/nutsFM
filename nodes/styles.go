package nodes

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SetDefaultNodeStyles устанавливает стили по умолчанию для узлов дерева файлов.
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

// UpdateNodeStyles обновляет стили узлов дерева файлов в соответствии с результатами сортировки.
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

// ResetNodeStyles сбрасывает стили узлов дерева файлов к значениям по умолчанию.
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
