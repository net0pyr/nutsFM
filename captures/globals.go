// Пакет captures содержит необходимые горячие клавиши для приложения nutsFM.
package captures

import (
	"github.com/rivo/tview"
)

// TreeView представляет текущее дерево файлов.
var TreeView *tview.TreeView

// App представляет основное TUI-приложение.
var App *tview.Application

// RootPath хранит путь к корневой директории.
var RootPath *string

// CreateSearchApp функция для создания и запуска приложения поиска.
var CreateSearchApp func(searchTreeView *tview.TreeView, rootDir string)

// CreateMainApp функция для создания и запуска основного приложения.
var CreateMainApp func(rootDir string)

// SearchApp представляет приложение поиска.
var SearchApp *tview.Application

// SearchTreeView представляет дерево файлов в приложении поиска.
var SearchTreeView *tview.TreeView

// RootDir хранит путь к корневой директории для поиска.
var RootDir *string

// Pages представляет страницы приложения.
var Pages *tview.Pages

// SortField представляет поле ввода для сортировки.
var SortField = tview.NewInputField().
	SetLabel("/").
	SetFieldWidth(30)
