package vtui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	Default string = "default"
)

type ShortcutKey struct {
	Name string
	Key  key.Binding
	Func func(tea.Msg, *ShortcutKey) error
	Cmd  tea.Cmd
}

func (that *ShortcutKey) Execute(msg tea.Msg) (err error) {
	if that.Func != nil {
		err = that.Func(msg, that)
	}
	return
}

type KeyList []*ShortcutKey

func (that *KeyList) ShortHelp() (r []key.Binding) {
	for _, k := range *that {
		if k.Name == Default {
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
					k.Execute(msg)
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
	Func func(tea.Msg, *Message) error
	Cmd  tea.Cmd
}

func (that *Message) Execute(msg tea.Msg) (err error) {
	if that.Func != nil {
		return that.Func(msg, that)
	}
	return
}

type MessageList []*Message

func (that *MessageList) UpdateByMessage(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	if m, ok := msg.(*Message); ok && m != nil {
		if m.Func != nil {
			m.Execute(msg)
		}
		if m.Cmd != nil {
			cmds = append(cmds, m.Cmd)
		}
	}
	if len(cmds) > 0 {
		return tea.Batch(cmds...)
	}
	return nil
}
