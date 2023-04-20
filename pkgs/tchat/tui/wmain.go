package tui

import (
	"github.com/epiclabs-io/winman"
	"github.com/rivo/tview"
)

type WMain struct {
	name     string
	tui      ITui
	Input    *tview.InputField
	Message  *tview.TextArea
	TreeNode *tview.TreeNode
	TreeView *tview.TreeView
}

func NewMainWindow() *WMain {
	return &WMain{
		name:     "main",
		Input:    tview.NewInputField(),
		Message:  tview.NewTextArea(),
		TreeNode: tview.NewTreeNode("+ New Chat"),
		TreeView: tview.NewTreeView(),
	}
}

func (that *WMain) Name() string {
	return that.name
}

func (that *WMain) SetTui(t ITui) {
	that.tui = t
}

func (that *WMain) GetWindow() *winman.WindowBase {
	mainFlex := tview.NewFlex()
	that.Input.SetBorder(true)
	that.Message.SetBorder(true)

	form := tview.NewForm()
	form.AddButton("config", func() {
		if confWin := that.tui.SearchWindow("config"); confWin != nil {
			confWin.Show()
			that.tui.SetFocus(confWin)
		}
	})

	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	rightFlex.AddItem(that.Message, 0, 8, false)
	rightFlex.AddItem(that.Input, 0, 1, false)
	rightFlex.AddItem(form, 0, 1, false)

	mainFlex.AddItem(that.TreeView, 0, 3, false)
	mainFlex.AddItem(rightFlex, 0, 7, false)

	window := winman.NewWindow()
	window.SetRoot(mainFlex).
		SetResizable(true).
		SetDraggable(true).
		SetModal(false)

	window.SetBorder(true).SetTitle("ChatGPT").SetTitleAlign(tview.AlignCenter)
	window.SetRect(0, 0, 50, 30)

	window.AddButton(&winman.Button{
		Symbol:    'X',
		Alignment: winman.ButtonLeft,
		OnClick:   that.quit,
	})

	var maxMinButton *winman.Button
	maxMinButton = &winman.Button{
		Symbol:    '▴',
		Alignment: winman.ButtonRight,
		OnClick: func() {
			if window.IsMaximized() {
				window.Restore()
				maxMinButton.Symbol = '▴'
			} else {
				window.Maximize()
				maxMinButton.Symbol = '▾'
			}
		},
	}

	window.AddButton(maxMinButton)
	window.Show()
	return window
}

func (that *WMain) quit() {
	that.tui.Quit(that.name)
}
