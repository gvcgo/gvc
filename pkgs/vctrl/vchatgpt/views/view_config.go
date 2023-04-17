package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/chatgpt"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/vtui"
)

var (
	FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	CursorStyle         = FocusedStyle.Copy()
	NoStyle             = lipgloss.NewStyle()
	HelpStyle           = BlurredStyle.Copy()
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	FocusedButton = FocusedStyle.Copy().Render("[ Submit ]")
	BlurredButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Submit"))
)

type ChatgptConfView struct {
	*ViewBase
	Conf       *chatgpt.ChatGPTConf
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func NewConfView() (cv *ChatgptConfView) {
	cv = &ChatgptConfView{
		ViewBase: NewBase("chatgpt_conf"),
		Conf:     chatgpt.NewChatGptConf(),
	}
	cv.Conf.GetOptions()
	idx := 0
	for _, opt := range cv.Conf.OptList {
		t := textinput.New()
		t.CursorStyle = CursorStyle
		t.CharLimit = 100
		t.Placeholder = opt.KName
		t.SetValue(opt.String())
		switch idx {
		case 0:
			t.Focus()
			t.PromptStyle = FocusedStyle
			t.TextStyle = FocusedStyle
		default:
		}
		cv.inputs[idx] = t
		idx++
	}
	return
}

func (that *ChatgptConfView) Keys() vtui.KeyList {
	kl := vtui.KeyList{}
	kl = append(kl, &vtui.ShortcutKey{
		Name: "ChangeCursorMode",
		Key:  key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "Change cursor mode.")),
		Func: func(m tea.Msg, sk *vtui.ShortcutKey) error {
			that.cursorMode++
			if that.cursorMode > cursor.CursorHide {
				that.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(that.inputs))
			for i := range that.inputs {
				cmds[i] = that.inputs[i].Cursor.SetMode(that.cursorMode)
			}
			sk.Cmd = tea.Batch(cmds...)
			return nil
		},
	})
	return kl
}

func (that *ChatgptConfView) Msgs() vtui.MessageList {
	ml := vtui.MessageList{}
	return ml
}

func (that *ChatgptConfView) View() string {
	return ""
}
