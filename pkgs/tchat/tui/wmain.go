package tui

import (
	"github.com/epiclabs-io/winman"
	"github.com/rivo/tview"
)

type WMain struct {
	name    string
	tui     ITui
	Input   *tview.InputField
	Message *tview.TextArea
}

func NewMainWindow() *WMain {
	return &WMain{
		name:    "main",
		Input:   tview.NewInputField(),
		Message: tview.NewTextArea(),
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

	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	rightFlex.AddItem(that.Message, 0, 8, false)
	rightFlex.AddItem(that.Input, 0, 1, false)

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
		OnClick:   that.tui.Quit,
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
