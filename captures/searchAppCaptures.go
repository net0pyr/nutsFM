package captures

import (
	"log"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
)

// SearchAppCaptures обрабатывает события клавиатуры для приложения поиска.
func SearchAppCaptures(event *tcell.EventKey) *tcell.EventKey {
	switch {
	case event.Key() == tcell.KeyEsc:
		SearchApp.Stop()
		CreateMainApp(*RootDir)
		return nil
	case event.Key() == tcell.KeyEnter:
		node := SearchTreeView.GetCurrentNode()
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
			log.Printf("Не удалось получить доступ к %s: %v", path, err)
			return event
		}

		if info.IsDir() {
			SearchApp.Stop()
			CreateMainApp(path)
		} else {
			cmd := exec.Command("xdg-open", path)
			err := cmd.Start()
			if err != nil {
				log.Printf("Не удалось открыть файл %s: %v", path, err)
			}
		}
	default:
		return event
	}
	return event
}
