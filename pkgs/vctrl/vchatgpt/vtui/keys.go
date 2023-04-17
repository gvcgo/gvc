package vtui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	Help string = "help"
	Quit string = "quit"
)

type ShortcutKey struct {
	Name string
	Key  key.Binding
	Func func(tea.Msg) error
	Cmd  tea.Cmd
}

type KeyList []*ShortcutKey

func (that *KeyList) ShortHelp() (r []key.Binding) {
	for _, k := range *that {
		if k.Name == Help || k.Name == Quit {
			r = append(r, k.Key)
		}
	}
	return
}

func (that *KeyList) FullHelp() (r [][]key.Binding) {
	kList := []key.Binding{}
	for _, k := range *that {
		kList = append(kList, k.Key)
	}
	r = append(r, kList)
	return
}

func (that *KeyList) UpdateByKeys(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	if m, ok := msg.(tea.KeyMsg); ok {
		for _, k := range *that {
			if key.Matches(m, k.Key) {
				if k.Func != nil {
					k.Func(msg)
				}
				if k.Cmd != nil {
					cmds = append(cmds, k.Cmd)
				}
			}
		}
	}

	if len(cmds) > 0 {
		return tea.Batch(cmds...)
	}
	return nil
}

type Message struct {
	Func func(tea.Msg) error
	Cmd  tea.Cmd
}

type MessageList []*Message

func (that *MessageList) UpdateByMessage(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	if m, ok := msg.(*Message); ok && m != nil {
		if m.Cmd != nil {
			cmds = append(cmds, m.Cmd)
		}
		if m.Func != nil {
			m.Func(msg)
		}
	}
	if len(cmds) > 0 {
		return tea.Batch(cmds...)
	}
	return nil
}
