package vtui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	CmdHandler  func(tea.Msg) tea.Cmd
	InitHandler func() tea.Cmd
)

type IView interface {
	Name() string
	View() string
	Keys() KeyList
	Msgs() MessageList
	ExtraCmdHandlers() []CmdHandler
	IsEnabled() bool
	Disable()
	Enable()
	SetModel(IModel)
}

type IModel interface {
	GetKeys() help.KeyMap
	DisableOthers(string)
	EnableDefault()
}

type Model struct {
	Views            map[string]IView
	Keys             KeyList
	Msgs             MessageList
	ExtraCmdHandlers []CmdHandler
	InitList         []InitHandler
}

func NewModel() (m *Model) {
	m = &Model{
		Views: make(map[string]IView, 0),
		Keys:  make(KeyList, 0),
		Msgs:  make(MessageList, 0),
	}
	return
}

func (that *Model) GetKeys() help.KeyMap {
	return &that.Keys
}

func (that *Model) DisableOthers(name string) {
	for _, v := range that.Views {
		if v.Name() == name {
			continue
		}
		v.Disable()
	}
}

func (that *Model) EnableDefault() {
	for _, v := range that.Views {
		if v.Name() == Default {
			v.Enable()
		}
	}
}

func (that *Model) AddInit(f func() tea.Cmd) {
	that.InitList = append(that.InitList, f)
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
	if v != nil {
		that.Keys = append(that.Keys, v.Keys()...)
		that.Msgs = append(that.Msgs, v.Msgs()...)
		that.Views[v.Name()] = v
		v.SetModel(that)
		that.ExtraCmdHandlers = v.ExtraCmdHandlers()
	}
}

func (that *Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, f := range that.InitList {
		cmds = append(cmds, f())
	}
	return tea.Batch(cmds...)
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
	for _, f := range that.ExtraCmdHandlers {
		cmd = f(msg)
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
