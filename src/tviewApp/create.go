package tviewApp

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

	isSearchMode := false

	inputField := tview.NewInputField().
		SetLabel("/").
		SetFieldWidth(30)

	var searchResults []*tview.TreeNode
	var currentIndex int

	updateNodeStyles := func() {
		for _, node := range searchResults {
			node.SetColor(tcell.ColorYellow)
		}
		if len(searchResults) > 0 {
			searchResults[currentIndex].SetColor(tcell.ColorRed)
		}
	}

	search := func(query string) {
		searchResults = []*tview.TreeNode{}
		currentIndex = 0
		treeView.GetRoot().Walk(func(node, parent *tview.TreeNode) bool {
			if strings.Contains(strings.ToLower(node.GetText()), strings.ToLower(query)) {
				searchResults = append(searchResults, node)
			}
			return true
		})
		updateNodeStyles()
		if len(searchResults) > 0 {
			treeView.SetCurrentNode(searchResults[currentIndex])
		}
	}

	resetNodeStyles := func() {
		treeView.GetRoot().Walk(func(node, parent *tview.TreeNode) bool {
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

	inputField.SetChangedFunc(func(text string) {
		search(text)
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == '/' {
			flex := tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(treeView, 0, 1, true).
				AddItem(inputField, 1, 1, false)
			app.SetRoot(flex, true).SetFocus(inputField)
			isSearchMode = true
			return nil
		} else if event.Key() == tcell.KeyEsc {
			inputField.SetText("")
			resetNodeStyles()
			flex := tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(treeView, 0, 1, true)
			app.SetRoot(flex, true).SetFocus(treeView)
			isSearchMode = false
			return nil
		} else if event.Key() == tcell.KeyUp && isSearchMode {
			if len(searchResults) > 0 {
				currentIndex = (currentIndex - 1 + len(searchResults)) % len(searchResults)
				updateNodeStyles()
				treeView.SetCurrentNode(searchResults[currentIndex])
			}
			return nil
		} else if event.Key() == tcell.KeyDown && isSearchMode {
			if len(searchResults) > 0 {
				currentIndex = (currentIndex + 1) % len(searchResults)
				updateNodeStyles()
				treeView.SetCurrentNode(searchResults[currentIndex])
			}
			return nil
		} else if event.Key() == tcell.KeyEnter {
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
			return nil
		} else if event.Key() == tcell.KeyCtrlT {
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

	if err := app.SetRoot(treeView, true).Run(); err != nil {
		log.Fatal(err)
	}
}
