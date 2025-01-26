// Пакет nodes предоставляет инструменты для работы с узлами дерева файлов в приложении nutsFM.
package nodes

import (
	"github.com/rivo/tview"
)

// TreeView представляет текущее дерево файлов.
var TreeView *tview.TreeView

// SortResults хранит результаты сортировки узлов дерева файлов.
var SortResults = []*tview.TreeNode{}

// CurrentIndex хранит текущий индекс в результатах сортировки.
var CurrentIndex = 0
