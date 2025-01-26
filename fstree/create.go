// Пакет fstree предоставляет функции для работы с деревом файловой системой в приложении nutsFM.
package fstree

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Create создает дерево файлов и директорий для указанного пути.
func Create(node *tview.TreeNode, path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Не удалось прочитать директорию %s: %v", path, err)
		return
	}

	// Добавляем узел для родительской директории
	parentNode := tview.NewTreeNode("..").
		SetReference(filepath.Join(path, ".."))
	parentNode.SetColor(tcell.ColorPink)
	node.AddChild(parentNode)

	// Добавляем узлы для всех файлов и директорий в текущей директории
	for _, file := range files {
		childPath := filepath.Join(path, file.Name())
		childNode := tview.NewTreeNode(file.Name()).
			SetReference(childPath)

		info, err := os.Stat(childPath)
		if err != nil {
			log.Printf("Не удалось получить информацию о файле %s: %v", childPath, err)
			continue
		}

		if info.IsDir() {
			childNode.SetColor(tcell.ColorPink)
		} else {
			childNode.SetColor(tcell.ColorGreen)
		}

		childNode.Collapse()
		node.AddChild(childNode)
	}
}
