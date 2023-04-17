package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/vtui"
)

type DefaultMsg struct{}

type DefaultView struct {
	enabled bool
	name    string
	keys    vtui.KeyList
}

func NewDefaultView() *DefaultView {
	return &DefaultView{
		enabled: false,
		name:    "default",
		keys:    make(vtui.KeyList, 0),
	}
}

func (that *DefaultView) Name() string {
	return that.name
}

func (that *DefaultView) IsEnabled() bool {
	return that.enabled
}

func (that *DefaultView) Keys() vtui.KeyList {
	kl := vtui.KeyList{}
	kl = append(kl, &vtui.ShortcutKey{
		Key:  "esc",
		Help: "quit tui.",
		Cmd:  tea.Quit,
	})
	return kl
}

func (that *DefaultView) Msgs() vtui.MessageList {
	return vtui.MessageList{}
}

func (that *DefaultView) View() string {
	return ""
}
