package captures

import (
	"log"
	"os"

	"github.com/net0pyr/nutsFM/fstree"

	"github.com/gdamore/tcell/v2"
)

// ViewCaptures обрабатывает события клавиатуры для представления дерева файлов.
func ViewCaptures(event *tcell.EventKey) *tcell.EventKey {
	switch {
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
			log.Printf("Не удалось получить доступ к %s: %v", path, err)
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
	default:
		return event
	}
}
