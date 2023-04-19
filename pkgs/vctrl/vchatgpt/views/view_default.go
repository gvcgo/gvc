package views

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/vtui"
)

var DefaultMsg vtui.MsgType = "default"

var DefaultCmd tea.Cmd = func() tea.Msg {
	return vtui.Message{
		Name: vtui.Default,
		Type: DefaultMsg,
	}
}

var DefaultInit = func() tea.Cmd {
	return DefaultCmd
}

type DefaultView struct {
	*ViewBase
	help              help.Model
	addKeysRegistered bool
}

func NewDefaultView() (dv *DefaultView) {
	dv = &DefaultView{
		ViewBase: NewBase(vtui.Default),
		help:     help.New(),
	}
	dv.Enabled = true
	dv.ViewName = string(DefaultMsg)
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
			if that.Model != nil && !that.addKeysRegistered {
				vList := that.Model.GetViews()
				for _, v := range vList {
					v.AdditionalKeys()
				}
				that.addKeysRegistered = true
			}
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
		Func: func(m1 *vtui.Message) (tea.Cmd, error) {
			switch m1.Type {
			case DefaultMsg:
				that.Enabled = true
				that.help.ShowAll = false
			default:
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
