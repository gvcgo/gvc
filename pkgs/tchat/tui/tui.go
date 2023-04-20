package tui

import (
	"github.com/epiclabs-io/winman"
	"github.com/rivo/tview"
)

type Window struct {
	Index  int
	Window *winman.WindowBase
}

type TUI struct {
	App        *tview.Application
	Manager    *winman.Manager
	windowList map[string]*Window
	idx        int
}

func NewTui() (r *TUI) {
	return &TUI{
		App:        tview.NewApplication(),
		Manager:    winman.NewWindowManager(),
		windowList: map[string]*Window{},
	}
}

func (that *TUI) Register(w IWindow) {
	if _, ok := that.windowList[w.Name()]; !ok {
		that.idx++
		w.SetTui(that)
		window := w.GetWindow()
		that.windowList[w.Name()] = &Window{
			Window: window,
			Index:  that.idx,
		}
		that.Manager.AddWindow(window)
	} else {
		panic("window already registered")
	}
}

func (that *TUI) Setup() {
	msgBox := NewWMsg("quit",
		"Confirmation",
		"Really quit?",
		[]string{"Yes", "No"},
		func(clicked string) {
			if clicked == "Yes" {
				that.App.Stop()
			}
		})
	that.Register(msgBox)

	confWindow := NewConfWindow()
	that.Register(confWindow)

	mainWindow := NewMainWindow()
	that.Register(mainWindow)
}

func (that *TUI) SearchWindow(name string) *winman.WindowBase {
	if wind, ok := that.windowList[name]; ok {
		return wind.Window
	}
	return nil
}

func (that *TUI) Quit(name string) {
	if msgbox, ok := that.windowList["quit"]; ok && name == "main" {
		w := msgbox.Window
		w.Show()
		that.Manager.Center(w)
		that.SetFocus(w)

	} else {
		if w := that.SearchWindow(name); w != nil {
			w.Hide()
			that.SetFocus(that.Manager)
		}
	}
}

func (that *TUI) SetFocus(p tview.Primitive) {
	go that.App.QueueUpdateDraw(func() {
		that.App.SetFocus(p)
	})
}

func (that *TUI) Run() {
	that.Setup()
	if err := that.App.SetRoot(that.Manager, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
