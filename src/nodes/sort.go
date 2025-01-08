package nodes

import (
	"strings"

	"github.com/rivo/tview"
)

var Sort = func(query string) {
	if query == "" {
		SetDefaultNodeStyles()
		return
	}
	SortResults = []*tview.TreeNode{}
	CurrentIndex = 0
	TreeView.GetRoot().Walk(func(node, parent *tview.TreeNode) bool {
		if strings.Contains(strings.ToLower(node.GetText()), strings.ToLower(query)) {
			SortResults = append(SortResults, node)
		}
		return true
	})
	UpdateNodeStyles()
	if len(SortResults) > 0 {
		TreeView.SetCurrentNode(SortResults[CurrentIndex])
	}
}
