package modalWindows

import (
	"github.com/rivo/tview"
)

// ShowErrorModal отображает модальное окно с сообщением об ошибке.
func ShowErrorModal(message string, app *tview.Application, pages *tview.Pages, treeView *tview.TreeView) {
	errorModal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Закрыть"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage("errorModal")
			app.SetRoot(treeView, true)
			app.SetFocus(treeView)
		})

	pages.AddPage("errorModal", errorModal, true, true)
	app.SetRoot(pages, true)
	app.SetFocus(errorModal)
}
