package captures

import (
	"github.com/rivo/tview"
)

var TreeView *tview.TreeView
var App *tview.Application
var RootPath *string
var CreateSearchApp func(searchTreeView *tview.TreeView, rootDir string)
var CreateMainApp func(rootDir string)
var SearchApp *tview.Application
var SearchTreeView *tview.TreeView
var RootDir *string

var SortField = tview.NewInputField().
	SetLabel("/").
	SetFieldWidth(30)
