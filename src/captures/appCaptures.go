package captures

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/net0pyr/nutsFM/fstree"
	"github.com/net0pyr/nutsFM/nodes"

	"github.com/gdamore/tcell/v2"
	"github.com/msteinert/pam"
	"github.com/rivo/tview"
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
		// if App.GetFocus() != TreeView {
		// 	return event
		// }
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
			showErrorModal("Ошибка поиска: " + err.Error())
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
			showErrorModal("Ошибка доступа: " + err.Error())
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
				showErrorModal("Не удалось открыть файл: " + err.Error())
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
			showErrorModal("Ошибка доступа: " + err.Error())
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
		showConfirmDeleteModal(path, node)
		return nil

	default:
		return event
	}
}

func showErrorModal(message string) {
	errorModal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Закрыть"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			log.Println("Закрытие модального окна и возврат фокуса на TreeView")
			Pages.RemovePage("errorModal")
			App.SetRoot(TreeView, true)
			App.SetFocus(TreeView)
		})

	Pages.AddPage("errorModal", errorModal, true, true)
	App.SetRoot(Pages, true)
	App.SetFocus(errorModal)
	log.Println("Модальное окно отображено и фокус установлен на него")
}

func showConfirmDeleteModal(path string, node *tview.TreeNode) {
	log.Println("Вызов showConfirmDeleteModal")
	confirmModal := tview.NewModal().
		SetText("Вы уверены, что хотите удалить:\n" + path + "?").
		AddButtons([]string{"Да", "Нет"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			log.Println("Закрытие модального окна подтверждения удаления")
			Pages.RemovePage("confirmDeleteModal")
			if buttonLabel == "Да" {
				err := os.RemoveAll(path)
				if err != nil {
					if strings.Contains(err.Error(), "permission denied") {
						promptSudoPassword(path, node)
						return
					}
					log.Printf("Failed to delete %s: %v", path, err)
					showErrorModal("Ошибка удаления: " + err.Error())
					return
				}
				removeNodeFromTree(node)
			}
			App.SetRoot(TreeView, true)
			App.SetFocus(TreeView)
		})

	Pages.AddPage("confirmDeleteModal", confirmModal, true, true)
	App.SetRoot(Pages, true)
	App.SetFocus(confirmModal)
	log.Println("Модальное окно подтверждения удаления отображено и фокус установлен на него")
}

func getParentNode(root, target *tview.TreeNode) *tview.TreeNode {
	for _, child := range root.GetChildren() {
		if child == target {
			return root
		}
		if p := getParentNode(child, target); p != nil {
			return p
		}
	}
	return nil
}

func authenticateWithPam(user, password string) error {
	t, err := pam.StartFunc("sudo", user, func(s pam.Style, msg string) (string, error) {
		switch s {
		case pam.PromptEchoOff, pam.PromptEchoOn:
			return password, nil
		default:
			return "", nil
		}
	})
	if err != nil {
		return err
	}
	return t.Authenticate(0)
}

func promptSudoPassword(path string, node *tview.TreeNode) {
	log.Println("Вызов promptSudoPassword")
	passwordField := tview.NewInputField().
		SetLabel("Введите пароль sudo (PAM): ").
		SetMaskCharacter('*')
	passwordField.
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				pass := passwordField.GetText()
				// Проверяем пароль через PAM
				if err := authenticateWithPam("net0pyr", pass); err != nil {
					log.Printf("Ошибка аутентификации PAM: %v", err)
					showWrongPasswordModal(path, node)
					return
				}
				// Если PAM-проверка прошла, удаляем файл через sudo
				cmd := exec.Command("sudo", "-S", "rm", "-rf", path)
				cmd.Stdin = strings.NewReader(pass + "\n")
				if err := cmd.Run(); err != nil {
					log.Printf("Ошибка удаления через sudo %s: %v", path, err.Error())
					showWrongPasswordModal(path, node)
					return
				}
				removeNodeFromTree(node)
				App.SetRoot(TreeView, true)
				App.SetFocus(TreeView)
			}
		})
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(TreeView, 0, 3, false).
		AddItem(passwordField, 1, 1, true)

	App.SetRoot(layout, true)
	App.SetFocus(passwordField)
	log.Println("Форма ввода пароля отображена и фокус установлен на нее")
}

func showWrongPasswordModal(path string, node *tview.TreeNode) {
	log.Println("Вызов showWrongPasswordModal")
	modal := tview.NewModal().
		SetText("Пароль неверный").
		AddButtons([]string{"Отмена", "Ввести заново"}).
		SetDoneFunc(func(i int, label string) {
			log.Println("Закрытие модального окна ошибки ввода пароля")
			Pages.RemovePage("wrongPassModal")
			if label == "Ввести заново" {
				promptSudoPassword(path, node)
				return
			}
			App.SetRoot(TreeView, true)
			App.SetFocus(TreeView)
		})

	Pages.AddPage("wrongPassModal", modal, true, true)
	App.SetRoot(Pages, true)
	App.SetFocus(modal)
	log.Println("Модальное окно ошибки ввода пароля отображено и фокус установлен на него")
}
func removeNodeFromTree(node *tview.TreeNode) {
	parent := getParentNode(TreeView.GetRoot(), node)
	if parent != nil {
		parent.RemoveChild(node)
	}
}
