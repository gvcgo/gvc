package views

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/vtui"
)

type DefaultMsg struct{}

type DefaultView struct {
	*ViewBase
	help     help.Model
	viewport viewport.Model
}

func NewDefaultView() *DefaultView {
	return &DefaultView{
		ViewBase: NewBase(vtui.Default),
		help:     help.New(),
		viewport: viewport.New(50, 5),
	}
}

func (that *DefaultView) Keys() vtui.KeyList {
	kl := vtui.KeyList{}
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key:  key.NewBinding(key.WithKeys("esc", "ctrl+c"), key.WithHelp("esc", "Quit tui."), key.WithHelp("ctrl+c", "Quit tui.")),
		Cmd:  tea.Quit,
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key:  key.NewBinding(key.WithKeys("ctrl+h"), key.WithHelp("ctrl+h", "Show/hide help info.")),
		Func: func(m tea.Msg, sk *vtui.ShortcutKey) error {
			that.help.ShowAll = !that.help.ShowAll
			that.Enabled = true
			return nil
		},
	})
	return kl
}

func (that *DefaultView) Msgs() vtui.MessageList {
	return vtui.MessageList{}
}

func (that *DefaultView) View() string {
	r := ""
	if that.help.ShowAll {
		r = that.help.View(that.Model.GetKeys())
	}
	return r
}
