package tui

import (
	"github.com/epiclabs-io/winman"
	"github.com/rivo/tview"
)

type ITui interface {
	SetFocus(tview.Primitive)
	Quit(string)
	SearchWindow(string) *winman.WindowBase
}

type IWindow interface {
	SetTui(ITui)
	GetWindow() *winman.WindowBase
	Name() string
}
