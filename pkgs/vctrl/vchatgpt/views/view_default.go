package views

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/vtui"
)

type DefaultMsg struct{}

var DefaultCmd tea.Cmd = func() tea.Msg {
	return vtui.Message{Name: vtui.Default}
}

var DefaultInit = func() tea.Cmd {
	return DefaultCmd
}

type DefaultView struct {
	*ViewBase
	help help.Model
}

func NewDefaultView() (dv *DefaultView) {
	dv = &DefaultView{
		ViewBase: NewBase(vtui.Default),
		help:     help.New(),
	}
	dv.Enabled = true
	return
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
				return nil, nil
			}
			return nil, nil
		},
	})
	return ml
}

func (that *DefaultView) View() string {
	if that.help.ShowAll {
		return that.help.View(that.Model.GetKeys())
	}
	r := ""
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#f67ba0")).
		Width(50).Height(3).PaddingTop(1).Align(lipgloss.Center)

	welcomeStr := style.Render("Welcome to GVC ChatGPT")
	helpStr := HelpStyle.Render(" (Press ctrl+h to show/hide help info for [Shortcut Keys].)")
	quitStr := HelpStyle.Render(" (Press esc/ctrl+c to quit GVC ChatGPT Tui.)")
	r = lipgloss.JoinVertical(lipgloss.Center, welcomeStr, "\n", helpStr, quitStr)
	return r
}
