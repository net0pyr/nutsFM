package captures

import (
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/net0pyr/nutsFM/fstree"
	"github.com/net0pyr/nutsFM/history"
	"github.com/net0pyr/nutsFM/modalWindows"
	"github.com/net0pyr/nutsFM/nodes"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/term"
)

var isSearchMode, isSortMode = false, false

var searchField = tview.NewInputField().
	SetLabel("Search: ").
	SetFieldWidth(30)

func AppCaptures(event *tcell.EventKey) *tcell.EventKey {
	switch {
	case event.Key() == tcell.KeyRune && event.Rune() == '/':
		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(TreeView, 0, 1, true).
			AddItem(SortField, 1, 1, false)
		App.SetRoot(flex, true).SetFocus(SortField)
		isSortMode = true
		return nil

	case event.Key() == tcell.KeyCtrlF:
		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(TreeView, 0, 1, true).
			AddItem(searchField, 1, 1, false)
		App.SetRoot(flex, true).SetFocus(searchField)
		isSearchMode = true
		return nil

	case event.Key() == tcell.KeyEsc:
		SortField.SetText("")
		searchField.SetText("")
		nodes.ResetNodeStyles()
		App.SetRoot(TreeView, true).SetFocus(TreeView)
		isSortMode = false
		isSearchMode = false
		return nil

	case event.Key() == tcell.KeyUp && isSortMode:
		if len(nodes.SortResults) > 0 {
			nodes.CurrentIndex = (nodes.CurrentIndex - 1 + len(nodes.SortResults)) % len(nodes.SortResults)
			nodes.UpdateNodeStyles()
			TreeView.SetCurrentNode(nodes.SortResults[nodes.CurrentIndex])
		}
		return nil

	case event.Key() == tcell.KeyDown && isSortMode:
		if len(nodes.SortResults) > 0 {
			nodes.CurrentIndex = (nodes.CurrentIndex + 1) % len(nodes.SortResults)
			nodes.UpdateNodeStyles()
			TreeView.SetCurrentNode(nodes.SortResults[nodes.CurrentIndex])
		}
		return nil

	case event.Key() == tcell.KeyEnter && isSearchMode:
		files, err := fstree.Find(*RootPath, searchField.GetText())
		if err != nil {
			log.Printf("Failed to search: %v", err)
			modalWindows.ShowErrorModal("Ошибка поиска: "+err.Error(), App, Pages, TreeView)
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

	case event.Key() == tcell.KeyEnter:
		if App.GetFocus() != TreeView {
			return event
		}

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
			modalWindows.ShowErrorModal("Ошибка поиска: "+err.Error(), App, Pages, TreeView)
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
				modalWindows.ShowErrorModal("Ошибка поиска: "+err.Error(), App, Pages, TreeView)
			}
		}
		return nil

	case event.Key() == tcell.KeyCtrlT:
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
			modalWindows.ShowErrorModal("Ошибка поиска: "+err.Error(), App, Pages, TreeView)
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

	case event.Key() == tcell.KeyDelete:
		log.Println("Нажата клавиша Delete")
		node := TreeView.GetCurrentNode()
		if node == nil {
			log.Println("Текущий узел не найден")
			return event
		}

		ref := node.GetReference()
		if ref == nil {
			log.Println("Ссылка на текущий узел не найдена")
			return event
		}

		path := ref.(string)
		modalWindows.ShowConfirmDeleteModal(path, node, App, Pages, TreeView)
		return nil

	case event.Key() == tcell.KeyRune && event.Rune() == ':':
		currentUser, err := user.Current()
		var userLabel string
		if err != nil {
			userLabel = "I don't know your name: "
		} else {
			userLabel = currentUser.Username + ": "
		}

		cmdInput := tview.NewInputField().
			SetLabel(userLabel)

		historyCommands := history.GetAll()
		historyIndex := len(historyCommands)

		outputView := tview.NewTextView().
			SetWrap(true).
			SetDynamicColors(true).
			SetScrollable(true)

		cmdInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch {
			case event.Key() == tcell.KeyTab:
				App.SetFocus(outputView)
				return nil

			case event.Key() == tcell.KeyUp:
				if historyIndex > 0 {
					historyIndex--
					cmdInput.SetText(historyCommands[historyIndex])
				}
				return nil

			case event.Key() == tcell.KeyDown:
				if historyIndex < len(historyCommands)-1 {
					historyIndex++
					cmdInput.SetText(historyCommands[historyIndex])
				} else {
					historyIndex = len(historyCommands)
					cmdInput.SetText("")
				}
				return nil

			default:
				return event
			}
		})

		cmdInput.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				root := TreeView.GetRoot()
				if root == nil || root.GetReference() == nil {
					return
				}
				dir := root.GetReference().(string)
				userCmd := cmdInput.GetText()

				history.AppendHistory(userCmd)

				cmd := exec.Command("bash", "-c", userCmd)
				cmd.Dir = dir
				out, err := cmd.CombinedOutput()
				output := string(out)
				if err != nil {
					output += "\nОшибка: " + err.Error()
				}

				outputView.SetText(output)

				screenWidth, _, termErr := term.GetSize(int(os.Stdout.Fd()))
				if termErr != nil || screenWidth < 1 {
					screenWidth = 80
				}

				lines := 0
				maxWidth := screenWidth
				if maxWidth < 1 {
					maxWidth = 1
				}
				for _, line := range strings.Split(output, "\n") {
					lineLen := len(line)
					if lineLen == 0 {
						lines++
						continue
					}
					linesNeeded := lineLen / maxWidth
					if lineLen%maxWidth != 0 {
						linesNeeded++
					}
					lines += linesNeeded
				}

				if lines < 1 {
					lines = 1
				}
				if lines > 10 {
					lines = 10
				}

				outputView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					if event.Key() == tcell.KeyTab {
						App.SetFocus(cmdInput)
						return nil
					}
					return event
				})

				newRootNode := tview.NewTreeNode(dir).
					SetReference(dir)
				fstree.Create(newRootNode, dir)
				TreeView.SetRoot(newRootNode).
					SetCurrentNode(newRootNode)
				nodes.ResetNodeStyles()

				layout := tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(TreeView, 0, 1, false).
					AddItem(cmdInput, 1, 1, false).
					AddItem(outputView, lines, 1, false)

				App.SetRoot(layout, true)
				App.SetFocus(cmdInput)
			}
		})

		bashFlex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(TreeView, 0, 1, false).
			AddItem(cmdInput, 1, 1, true)

		App.SetRoot(bashFlex, true)
		App.SetFocus(cmdInput)
		return nil

	default:
		return event
	}
}
