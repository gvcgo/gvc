package vtui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type IView interface {
	Name() string
	View() string
	Keys() KeyList
	Msgs() MessageList
	IsEnabled() bool
}

type Model struct {
	Views    map[string]IView
	Keys     KeyList
	Msgs     MessageList
	initFunc func() tea.Cmd
}

func NewModel() (m *Model) {
	m = &Model{
		Views: make(map[string]IView, 0),
		Keys:  make(KeyList, 0),
		Msgs:  make(MessageList, 0),
	}
	return
}

func (that *Model) SetInit(f func() tea.Cmd) {
	that.initFunc = f
}

func (that *Model) RegisterView(v IView) {
	if that.Views == nil {
		that.Views = make(map[string]IView, 0)
	}
	if that.Keys == nil {
		that.Keys = make(KeyList, 0)
	}
	if that.Msgs == nil {
		that.Msgs = make(MessageList, 0)
	}
	that.Keys = append(that.Keys, v.Keys()...)
	that.Msgs = append(that.Msgs, v.Msgs()...)
	that.Views[v.Name()] = v
}

func (that *Model) Init() tea.Cmd {
	if that.initFunc != nil {
		return that.initFunc()
	}
	return nil
}

func (that *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	cmd := that.Keys.UpdateByKeys(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	cmd = that.Msgs.UpdateByMessage(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return that, tea.Batch(cmds...)
}

func (that *Model) View() string {
	result := []string{}
	for _, view := range that.Views {
		if view.IsEnabled() {
			result = append(result, view.View())
		}
	}
	return strings.Join(result, "\n")
}
