package tui

import (
	"github.com/epiclabs-io/winman"
	"github.com/rivo/tview"
)

type ITui interface {
	SetFocus(p tview.Primitive)
	Quit()
}

type IWindow interface {
	SetTui(ITui)
	GetWindow() *winman.WindowBase
	Name() string
}
