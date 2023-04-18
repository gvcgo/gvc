package views

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/vtui"
)

type DefaultMsg struct{}

var DefaultCmd tea.Cmd = func() tea.Msg {
	return vtui.Message{Name: vtui.Default}
}

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
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			return tea.Quit, nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key:  key.NewBinding(key.WithKeys("ctrl+h"), key.WithHelp("ctrl+h", "Show/hide help info.")),
		Func: func(m tea.KeyMsg) (tea.Cmd, error) {
			that.help.ShowAll = !that.help.ShowAll
			that.Enabled = true
			that.Model.DisableOthers(that.Name())
			return nil, nil
		},
	})
	return kl
}

func (that *DefaultView) Msgs() vtui.MessageList {
	ml := vtui.MessageList{}
	ml = append(ml, &vtui.Message{
		Name: that.ViewName,
		Func: func(m1 tea.Msg) (tea.Cmd, error) {
			if m, ok := m1.(*vtui.Message); ok && m.Name == that.ViewName {
				that.Enabled = true
				that.help.ShowAll = false
				return tea.ClearScreen, nil
			}
			return nil, nil
		},
	})
	return ml
}

func (that *DefaultView) View() string {
	r := ""
	if that.help.ShowAll {
		r = that.help.View(that.Model.GetKeys())
	}
	return r
}
