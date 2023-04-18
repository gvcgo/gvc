package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/moqsien/gvc/pkgs/utils"
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
		ViewBase: NewBase("chatgpt_config"),
		Conf:     chatgpt.NewChatGptConf(),
	}
	cv.Conf.GetOptions()
	cv.inputs = make([]textinput.Model, len(cv.Conf.GetOptOrder()))
	idx := 0
	maxLength := utils.FindMaxLengthOfStringList(cv.Conf.GetOptOrder())
	for _, kname := range cv.Conf.GetOptOrder() {
		opt := cv.Conf.GetOptList()[kname]
		t := textinput.New()
		t.CursorStyle = CursorStyle
		t.CharLimit = 100
		t.Placeholder = opt.KName
		t.Prompt = fmt.Sprintf("%s: %s", kname, strings.Repeat(" ", maxLength-len(kname)))
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

func (that *ChatgptConfView) inputHandler(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(that.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range that.inputs {
		that.inputs[i], cmds[i] = that.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (that *ChatgptConfView) ExtraCmdHandlers() []vtui.CmdHandler {
	return []vtui.CmdHandler{
		that.inputHandler,
	}
}

func (that *ChatgptConfView) Keys() vtui.KeyList {
	kl := vtui.KeyList{}
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key:  key.NewBinding(key.WithKeys("ctrl+y"), key.WithHelp("ctrl+y", "Show chatgpt_config.")),
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			that.Enabled = true
			that.Model.DisableOthers(that.Name())
			return nil, nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key:  key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "Change cursor mode of [chatgpt_config].")),
		Func: func(m tea.KeyMsg) (tea.Cmd, error) {
			that.cursorMode++
			if that.cursorMode > cursor.CursorHide {
				that.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(that.inputs))
			for i := range that.inputs {
				cmds[i] = that.inputs[i].Cursor.SetMode(that.cursorMode)
			}
			return tea.Batch(cmds...), nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key:  key.NewBinding(key.WithKeys("tab", "shift+tab"), key.WithHelp("tab", "Change focus to next input of [chatgpt_config].")),
		Func: func(msg tea.KeyMsg) (tea.Cmd, error) {
			s := msg.String()
			if s == "shift+tab" {
				that.focusIndex--
			} else {
				that.focusIndex++
			}

			if that.focusIndex > len(that.inputs) {
				that.focusIndex = 0
			} else if that.focusIndex < 0 {
				that.focusIndex = len(that.inputs)
			}

			cmds := make([]tea.Cmd, len(that.inputs))
			for i := 0; i <= len(that.inputs)-1; i++ {
				if i == that.focusIndex {
					// Set focused state
					cmds[i] = that.inputs[i].Focus()
					that.inputs[i].PromptStyle = FocusedStyle
					that.inputs[i].TextStyle = FocusedStyle
					continue
				}
				// Remove focused state
				that.inputs[i].Blur()
				that.inputs[i].PromptStyle = NoStyle
				that.inputs[i].TextStyle = NoStyle
			}

			return tea.Batch(cmds...), nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key:  key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "Submit chatgpt_config.")),
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			if that.focusIndex != len(that.inputs) {
				return nil, nil
			}
			for _, ipt := range that.inputs {
				kname := ipt.Placeholder
				value := ipt.Value()
				that.Conf.SetConfField(kname, value)
			}
			that.Conf.Restore()
			that.Enabled = false
			that.Model.EnableDefault()
			return DefaultCmd, nil
		},
	})
	return kl
}

func (that *ChatgptConfView) Msgs() vtui.MessageList {
	ml := vtui.MessageList{}
	return ml
}

func (that *ChatgptConfView) View() string {
	var b strings.Builder

	for i := range that.inputs {
		b.WriteString(that.inputs[i].View())
		if i < len(that.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	button := &BlurredButton
	if that.focusIndex == len(that.inputs) {
		button = &FocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	b.WriteString(HelpStyle.Render("cursor mode is "))
	b.WriteString(CursorModeHelpStyle.Render(that.cursorMode.String()))
	b.WriteString(HelpStyle.Render(" (ctrl+r to change style)"))
	return b.String()
}
