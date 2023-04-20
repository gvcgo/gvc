package tui

import (
	"github.com/epiclabs-io/winman"
	"github.com/moqsien/gvc/pkgs/tchat/gpt"
	"github.com/rivo/tview"
)

type WConf struct {
	name   string
	Form   *tview.Form
	Conf   *gpt.ChatGPTConf
	tui    ITui
	window *winman.WindowBase
}

func NewConfWindow() (r *WConf) {
	return &WConf{
		name: "config",
		Form: tview.NewForm(),
		Conf: gpt.NewChatGptConf(),
	}
}

func (that *WConf) SetTui(t ITui) {
	that.tui = t
}

func (that *WConf) GetWindow() *winman.WindowBase {
	form := tview.NewForm()
	window := winman.NewWindow().
		SetRoot(form).
		SetResizable(true).
		SetDraggable(true).
		SetModal(false)
	fields := that.Conf.GetOptOrder()
	optList := that.Conf.GetOptList()
	for _, optName := range fields {
		form.AddInputField(optName, optList[optName].String(), 100, nil, nil)
	}
	form.AddButton("save", func() {}).AddButton("cancel", func() {})

	window.SetBorder(true).SetTitle("config").SetTitleAlign(tview.AlignCenter)
	window.SetRect(8, 4, 100, 50)
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
	that.window = window
	return window
}

func (that *WConf) Name() string {
	return that.name
}

func (that *WConf) quit() {
	that.tui.Quit(that.name)
}
