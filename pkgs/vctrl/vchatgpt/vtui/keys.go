package vtui

import tea "github.com/charmbracelet/bubbletea"

type ShortcutKey struct {
	Key  string
	Help string
	Func func() error
	Cmd  tea.Cmd
}

type KeyList []*ShortcutKey

func (that *KeyList) UpdateByKeys(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	if m, ok := msg.(tea.KeyMsg); ok {
		for _, key := range *that {
			if m.String() == key.Key {
				if key.Func != nil {
					key.Func()
				}
				if key.Cmd != nil {
					cmds = append(cmds, key.Cmd)
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
	Func func() error
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
			m.Func()
		}
	}
	if len(cmds) > 0 {
		return tea.Batch(cmds...)
	}
	return nil
}
