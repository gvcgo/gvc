package vtui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	ErrMsg         error
	DeltaAnswerMsg string
	AnswerMsg      string
	SaveMsg        struct{}
)

const (
	Default string = "default"
)

type KeyFunc func(tea.KeyMsg) (tea.Cmd, error)

type ShortcutKey struct {
	Name string
	Key  key.Binding
	Func KeyFunc
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
					cmd, _ := k.Func(m)
					cmds = append(cmds, cmd)
				}
			}
		}
	}

	if len(cmds) > 0 {
		return tea.Batch(cmds...)
	}
	return nil
}

type MsgFunc func(msg tea.Msg) (tea.Cmd, error)

type Message struct {
	Name string
	Func MsgFunc
}

type MessageList []*Message

func (that *MessageList) UpdateByMessage(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	if m, ok := msg.(*Message); ok && m != nil {
		if m.Func != nil {
			cmd, _ := m.Func(msg)
			cmds = append(cmds, cmd)
		}
	}
	if len(cmds) > 0 {
		return tea.Batch(cmds...)
	}
	return nil
}
