package tui

import (
	"github.com/epiclabs-io/winman"
	"github.com/moqsien/gvc/pkgs/tchat/gpt"
	"github.com/rivo/tview"
)

type WConf struct {
	name string
	Form *tview.Form
	Conf *gpt.ChatGPTConf
	Tui  ITui
}

func NewConfWindow() (r *WConf) {
	return &WConf{
		name: "chatgpt_conf",
		Form: tview.NewForm(),
		Conf: gpt.NewChatGptConf(),
	}
}

func (that *WConf) SetTui(t ITui) {
	that.Tui = t
}

func (that *WConf) GetWindow() *winman.WindowBase {
	return nil
}

func (that *WConf) Name() string {
	return that.name
}
