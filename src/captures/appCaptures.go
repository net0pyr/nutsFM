package captures

import (
	"log"
	"os"
	"os/exec"

	"github.com/net0pyr/nutsFM/fstree"
	"github.com/net0pyr/nutsFM/nodes"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var isSearchMode, isSortMode = false, false

var searchField = tview.NewInputField().
	SetLabel("Search: ").
	SetFieldWidth(30)

func AppCaptures(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune && event.Rune() == '/' {
		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(TreeView, 0, 1, true).
			AddItem(SortField, 1, 1, false)
		App.SetRoot(flex, true).SetFocus(SortField)
		isSortMode = true
		return nil
	} else if event.Key() == tcell.KeyCtrlF {
		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(TreeView, 0, 1, true).
			AddItem(searchField, 1, 1, false)
		App.SetRoot(flex, true).SetFocus(searchField)
		isSearchMode = true
		return nil
	} else if event.Key() == tcell.KeyEsc {
		SortField.SetText("")
		searchField.SetText("")
		nodes.ResetNodeStyles()
		App.SetRoot(TreeView, true).SetFocus(TreeView)
		isSortMode = false
		isSearchMode = false
		return nil
	} else if event.Key() == tcell.KeyUp && isSortMode {
		if len(nodes.SortResults) > 0 {
			nodes.CurrentIndex = (nodes.CurrentIndex - 1 + len(nodes.SortResults)) % len(nodes.SortResults)
			nodes.UpdateNodeStyles()
			TreeView.SetCurrentNode(nodes.SortResults[nodes.CurrentIndex])
		}
		return nil
	} else if event.Key() == tcell.KeyDown && isSortMode {
		if len(nodes.SortResults) > 0 {
			nodes.CurrentIndex = (nodes.CurrentIndex + 1) % len(nodes.SortResults)
			nodes.UpdateNodeStyles()
			TreeView.SetCurrentNode(nodes.SortResults[nodes.CurrentIndex])
		}
		return nil
	} else if event.Key() == tcell.KeyEnter && isSearchMode {
		files, err := fstree.Find(*RootPath, searchField.GetText())
		if err != nil {
			log.Printf("Failed to search: %v", err)
			return event
		}
		headNode := tview.NewTreeNode("Results").
			SetReference(*RootPath)

		for _, file := range files {
			childNode := tview.NewTreeNode(file).
				SetReference(file)

			info, err := os.Stat(file)
			if err != nil {
				log.Printf("Failed to stat %s: %v", file, err)
				continue
			}

			if info.IsDir() {
				childNode.SetColor(tcell.ColorPink)
			} else {
				childNode.SetColor(tcell.ColorGreen)
			}

			childNode.Collapse()
			headNode.AddChild(childNode)
		}

		findTreeView := tview.NewTreeView().
			SetRoot(headNode).
			SetCurrentNode(headNode)

		SortField.SetText("")
		searchField.SetText("")
		App.Stop()
		CreateSearchApp(findTreeView, *RootPath)
		isSortMode = false
		isSearchMode = false

		return nil
	} else if event.Key() == tcell.KeyEnter {
		node := TreeView.GetCurrentNode()
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
			newRootNode := tview.NewTreeNode(path).
				SetReference(path)
			fstree.Create(newRootNode, path)
			TreeView.SetRoot(newRootNode).
				SetCurrentNode(newRootNode)
			nodes.ResetNodeStyles()
		} else {
			cmd := exec.Command("xdg-open", path)
			err := cmd.Start()
			if err != nil {
				log.Printf("Failed to open file %s: %v", path, err)
			}
		}
		return nil
	} else if event.Key() == tcell.KeyCtrlT {
		node := TreeView.GetCurrentNode()
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
}
