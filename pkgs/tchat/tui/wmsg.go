package tui

import (
	"github.com/epiclabs-io/winman"
	"github.com/rivo/tview"
)

type MsgCallback func(clicked string)

type WMsg struct {
	name     string
	title    string
	text     string
	buttons  []string
	callback MsgCallback
	tui      ITui
}

func NewWMsg(name, title, text string, buttons []string, cb MsgCallback) *WMsg {
	return &WMsg{
		name:     name,
		title:    title,
		text:     text,
		buttons:  buttons,
		callback: cb,
	}
}

func (that *WMsg) SetTui(t ITui) {
	that.tui = t
}

func (that *WMsg) Name() string {
	return that.name
}

func (that *WMsg) GetWindow() *winman.WindowBase {
	msgBox := winman.NewWindow()
	message := tview.NewTextView().SetText(that.text).SetTextAlign(tview.AlignCenter)
	buttonBar := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	content := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(message, 0, 1, false).
		AddItem(buttonBar, 1, 0, true)

	msgBox.SetRoot(content)
	msgBox.SetTitle(that.title).
		SetRect(4, 2, 30, 6)
	msgBox.Draggable = true
	msgBox.Modal = true

	for _, buttonText := range that.buttons {
		button := func(buttonText string) *tview.Button {
			return tview.NewButton(buttonText).SetSelectedFunc(func() {
				msgBox.Hide()
				that.callback(buttonText)
			})
		}(buttonText)
		buttonBar.AddItem(button, 0, 1, true)
	}

	return msgBox
}
