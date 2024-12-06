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

func CreateMainApp(rootDir string) {
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

	isSortMode := false

	sortField := tview.NewInputField().
		SetLabel("/").
		SetFieldWidth(30)

	var sortResults []*tview.TreeNode
	var currentIndex int

	setDefaultNodeStyles := func() {
		for _, node := range treeView.GetRoot().GetChildren() {
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

	updateNodeStyles := func() {
		setDefaultNodeStyles()
		for _, node := range sortResults {
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
		if len(sortResults) > 0 {
			sortResults[currentIndex].SetColor(tcell.ColorRed)
		}
	}

	sort := func(query string) {
		if query == "" {
			setDefaultNodeStyles()
			return
		}
		sortResults = []*tview.TreeNode{}
		currentIndex = 0
		treeView.GetRoot().Walk(func(node, parent *tview.TreeNode) bool {
			if strings.Contains(strings.ToLower(node.GetText()), strings.ToLower(query)) {
				sortResults = append(sortResults, node)
			}
			return true
		})
		updateNodeStyles()
		if len(sortResults) > 0 {
			treeView.SetCurrentNode(sortResults[currentIndex])
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

	sortField.SetChangedFunc(func(text string) {
		sort(text)
	})

	isSearchMode := false

	searchField := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldWidth(30)

	sortField.SetChangedFunc(func(text string) {
		sort(text)
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == '/' {
			flex := tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(treeView, 0, 1, true).
				AddItem(sortField, 1, 1, false)
			app.SetRoot(flex, true).SetFocus(sortField)
			isSortMode = true
			return nil
		} else if event.Key() == tcell.KeyCtrlF {
			flex := tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(treeView, 0, 1, true).
				AddItem(searchField, 1, 1, false)
			app.SetRoot(flex, true).SetFocus(searchField)
			isSearchMode = true
			return nil
		} else if event.Key() == tcell.KeyEsc {
			sortField.SetText("")
			searchField.SetText("")
			resetNodeStyles()
			app.SetRoot(treeView, true).SetFocus(treeView)
			isSortMode = false
			isSearchMode = false
			return nil
		} else if event.Key() == tcell.KeyUp && isSortMode {
			if len(sortResults) > 0 {
				currentIndex = (currentIndex - 1 + len(sortResults)) % len(sortResults)
				updateNodeStyles()
				treeView.SetCurrentNode(sortResults[currentIndex])
			}
			return nil
		} else if event.Key() == tcell.KeyDown && isSortMode {
			if len(sortResults) > 0 {
				currentIndex = (currentIndex + 1) % len(sortResults)
				updateNodeStyles()
				treeView.SetCurrentNode(sortResults[currentIndex])
			}
			return nil
		} else if event.Key() == tcell.KeyEnter && isSearchMode {
			files, err := fstree.Find(rootPath, searchField.GetText())
			if err != nil {
				log.Printf("Failed to search: %v", err)
				return event
			}
			headNode := tview.NewTreeNode("Results").
				SetReference(rootPath)

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

			sortField.SetText("")
			searchField.SetText("")
			app.Stop()
			createSearchApp(app, findTreeView, rootPath)
			isSortMode = false
			isSearchMode = false

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
