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
	enabled  bool
	name     string
	keys     vtui.KeyList
	model    vtui.IModel
	help     help.Model
	viewport viewport.Model
}

func NewDefaultView() *DefaultView {
	return &DefaultView{
		enabled:  false,
		name:     "default",
		keys:     make(vtui.KeyList, 0),
		help:     help.New(),
		viewport: viewport.New(50, 5),
	}
}

func (that *DefaultView) SetModel(m vtui.IModel) {
	that.model = m
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
		Name: vtui.Quit,
		Key:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Quit tui.")),
		Cmd:  tea.Quit,
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: vtui.Help,
		Key:  key.NewBinding(key.WithKeys("ctrl+h"), key.WithHelp("ctrl+h", "Show/hide help info.")),
		Func: func(m tea.Msg) error {
			that.help.ShowAll = !that.help.ShowAll
			that.enabled = true
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
		r = that.help.View(that.model.GetKeys())
	}
	return r
}
