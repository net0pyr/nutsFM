// Пакет modalWindows предоставляет инструменты для отображения модальных окон в приложении nutsFM.
package modalWindows

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/msteinert/pam"
	"github.com/rivo/tview"
)

// ShowConfirmDeleteModal отображает модальное окно подтверждения удаления.
func ShowConfirmDeleteModal(path string, node *tview.TreeNode, app *tview.Application, pages *tview.Pages, treeView *tview.TreeView) {
	log.Println("Вызов showConfirmDeleteModal")
	confirmModal := tview.NewModal().
		SetText("Вы уверены, что хотите удалить:\n" + path + "?").
		AddButtons([]string{"Да", "Нет"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			log.Println("Закрытие модального окна подтверждения удаления")
			pages.RemovePage("confirmDeleteModal")
			if buttonLabel == "Да" {
				err := os.RemoveAll(path)
				if err != nil {
					if strings.Contains(err.Error(), "permission denied") {
						promptSudoPassword(path, node, app, treeView, pages)
						return
					}
					log.Printf("Failed to delete %s: %v", path, err)
					ShowErrorModal("Ошибка поиска: "+err.Error(), app, pages, treeView)
					return
				}
				removeNodeFromTree(node, treeView)
			}
			app.SetRoot(treeView, true)
			app.SetFocus(treeView)
		})

	pages.AddPage("confirmDeleteModal", confirmModal, true, true)
	app.SetRoot(pages, true)
	app.SetFocus(confirmModal)
	log.Println("Модальное окно подтверждения удаления отображено и фокус установлен на него")
}

// getParentNode возвращает родительский узел для указанного узла.
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

// authenticateWithPam выполняет аутентификацию с использованием PAM.
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

// promptSudoPassword отображает форму ввода пароля sudo.
func promptSudoPassword(path string, node *tview.TreeNode, app *tview.Application, treeView *tview.TreeView, pages *tview.Pages) {
	log.Println("Вызов promptSudoPassword")
	passwordField := tview.NewInputField().
		SetLabel("Введите пароль sudo (PAM): ").
		SetMaskCharacter('*')
	passwordField.
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				pass := passwordField.GetText()
				if err := authenticateWithPam("net0pyr", pass); err != nil {
					log.Printf("Ошибка аутентификации PAM: %v", err)
					showWrongPasswordModal(path, node, app, pages, treeView)
					return
				}
				cmd := exec.Command("sudo", "-S", "rm", "-rf", path)
				cmd.Stdin = strings.NewReader(pass + "\n")
				if err := cmd.Run(); err != nil {
					log.Printf("Ошибка удаления через sudo %s: %v", path, err.Error())
					showWrongPasswordModal(path, node, app, pages, treeView)
					return
				}
				removeNodeFromTree(node, treeView)
				app.SetRoot(treeView, true)
				app.SetFocus(treeView)
			}
		})
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(treeView, 0, 3, false).
		AddItem(passwordField, 1, 1, true)

	app.SetRoot(layout, true)
	app.SetFocus(passwordField)
	log.Println("Форма ввода пароля отображена и фокус установлен на нее")
}

// showWrongPasswordModal отображает модальное окно ошибки ввода пароля.
func showWrongPasswordModal(path string, node *tview.TreeNode, app *tview.Application, pages *tview.Pages, treeView *tview.TreeView) {
	log.Println("Вызов showWrongPasswordModal")
	modal := tview.NewModal().
		SetText("Пароль неверный").
		AddButtons([]string{"Отмена", "Ввести заново"}).
		SetDoneFunc(func(i int, label string) {
			log.Println("Закрытие модального окна ошибки ввода пароля")
			pages.RemovePage("wrongPassModal")
			if label == "Ввести заново" {
				promptSudoPassword(path, node, app, treeView, pages)
				return
			}
			app.SetRoot(treeView, true)
			app.SetFocus(treeView)
		})

	pages.AddPage("wrongPassModal", modal, true, true)
	app.SetRoot(pages, true)
	app.SetFocus(modal)
	log.Println("Модальное окно ошибки ввода пароля отображено и фокус установлен на него")
}

// removeNodeFromTree удаляет узел из дерева.
func removeNodeFromTree(node *tview.TreeNode, treeView *tview.TreeView) {
	parent := getParentNode(treeView.GetRoot(), node)
	if parent != nil {
		parent.RemoveChild(node)
	}
}
