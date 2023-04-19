package views

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/chatgpt"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/vtui"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

var textareaKeys = textarea.KeyMap{
	CharacterForward:           key.NewBinding(key.WithKeys("right", "ctrl+f")),
	CharacterBackward:          key.NewBinding(key.WithKeys("left", "ctrl+b")),
	WordForward:                key.NewBinding(key.WithKeys("alt+right", "alt+f")),
	WordBackward:               key.NewBinding(key.WithKeys("alt+left", "alt+b")),
	LineNext:                   key.NewBinding(key.WithKeys("down")),
	LinePrevious:               key.NewBinding(key.WithKeys("up")),
	DeleteWordBackward:         key.NewBinding(key.WithKeys("alt+backspace", "ctrl+w")),
	DeleteWordForward:          key.NewBinding(key.WithKeys("alt+delete", "alt+d")),
	DeleteAfterCursor:          key.NewBinding(key.WithKeys("ctrl+k")),
	DeleteBeforeCursor:         key.NewBinding(key.WithKeys("ctrl+u")),
	InsertNewline:              key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "insert new line")),
	DeleteCharacterBackward:    key.NewBinding(key.WithKeys("backspace")),
	DeleteCharacterForward:     key.NewBinding(key.WithKeys("delete")),
	LineStart:                  key.NewBinding(key.WithKeys("home", "ctrl+a")),
	LineEnd:                    key.NewBinding(key.WithKeys("end", "ctrl+e")),
	Paste:                      key.NewBinding(key.WithKeys("ctrl+v", "alt+v"), key.WithHelp("ctrl+v", "paste")),
	InputBegin:                 key.NewBinding(key.WithKeys("alt+<", "ctrl+home")),
	InputEnd:                   key.NewBinding(key.WithKeys("alt+>", "ctrl+end")),
	CapitalizeWordForward:      key.NewBinding(key.WithKeys("alt+c")),
	LowercaseWordForward:       key.NewBinding(key.WithKeys("alt+l")),
	UppercaseWordForward:       key.NewBinding(key.WithKeys("alt+u")),
	TransposeCharacterBackward: key.NewBinding(key.WithDisabled()),
}

var viewportKeys = viewport.KeyMap{
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	HalfPageUp:   key.NewBinding(key.WithDisabled()),
	HalfPageDown: key.NewBinding(key.WithDisabled()),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	),
}

var ChatgptInit vtui.InitHandler = func() tea.Cmd {
	cmds := []tea.Cmd{tea.EnterAltScreen}
	return tea.Batch(cmds...)
}

type InputMode int

const (
	InputModelSingleLine InputMode = iota
	InputModelMultiLine
)

func savePeriodically() tea.Cmd {
	return tea.Tick(15*time.Second, func(time.Time) tea.Msg {
		return vtui.NewMessage(
			chatgpt.VChatName,
			chatgpt.SaveMsg,
			struct{}{},
		)
	})
}

type ChatgptView struct {
	*ViewBase
	Conf        *chatgpt.ChatGPTConf
	Chatgpt     *chatgpt.VChat
	ConvManager *chatgpt.ConvManager
	viewport    viewport.Model
	spinner     spinner.Model
	textarea    textarea.Model
	rd          *glamour.TermRenderer
	spinning    bool
	// inputMode   InputMode
	err        error
	historyIdx int
	width      int
	height     int
}

func NewChatgptView() (cv *ChatgptView) {
	ta := textarea.New()
	ta.Placeholder = "send your message..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = -1
	ta.SetWidth(50)
	ta.SetHeight(1)
	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	vp := viewport.New(50, 5)
	spinner := spinner.New(spinner.WithSpinner(spinner.Points))
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithEnvironmentConfig(),
		glamour.WithWordWrap(0), // we do hard-wrapping ourselves
	)

	conf := chatgpt.NewChatGptConf()
	conf.Reload()
	cv = &ChatgptView{
		textarea:    ta,
		viewport:    vp,
		spinner:     spinner,
		rd:          renderer,
		Conf:        conf,
		ConvManager: chatgpt.NewConvManager(conf),
		Chatgpt:     chatgpt.NewVChat(conf),
	}
	cv.historyIdx = cv.ConvManager.Curr().Len()
	cv.viewport.KeyMap = viewportKeys
	cv.textarea.KeyMap = textareaKeys

	cv.ViewBase = NewBase("chatgpt_view")
	return
}

/*
textarea and viewport
*/
func (that *ChatgptView) handlePanels(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	that.textarea, cmd = that.textarea.Update(msg)
	cmds = append(cmds, cmd)
	that.viewport, cmd = that.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (that *ChatgptView) handleWindowSize(msg tea.Msg) (cmd tea.Cmd) {
	if message, ok := msg.(tea.WindowSizeMsg); ok {
		that.width = message.Width
		that.height = message.Height
		that.viewport.Width = message.Width
		// that.viewport.Height = message.Height - that.textarea.Height() - lipgloss.Height(that.RenderFooter())
		that.textarea.SetWidth(message.Width)
		// that.viewport.SetContent(that.RenderConversation(that.viewport.Width))
		that.viewport.GotoBottom()
	}
	return
}

func (that *ChatgptView) handleSpinnerTick(msg tea.Msg) (cmd tea.Cmd) {
	if message, ok := msg.(spinner.TickMsg); ok && that.spinning {
		that.spinner, cmd = that.spinner.Update(message)
	}
	return
}

func (that *ChatgptView) ExtraCmdHandlers() []vtui.CmdHandler {
	return []vtui.CmdHandler{
		that.handlePanels,
		that.handleWindowSize,
		that.handleSpinnerTick,
	}
}

func (that *ChatgptView) AdditionalKeys() {
	if that.Model != nil {
		that.Model.RegisterKeys(that.viewport.KeyMap, "chatgpt_viewport")
		that.Model.RegisterKeys(that.textarea.KeyMap, "chatgpt_textarea")
	}
}

func (that *ChatgptView) Keys() vtui.KeyList {
	kl := vtui.KeyList{}
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key: key.NewBinding(
			key.WithKeys("ctrl+w"),
			key.WithHelp("ctrl+w", "show chatgpt tui."),
		),
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			that.Enable()
			that.Model.DisableOthers(that.ViewName)
			return nil, nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit a message."),
		),
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			if that.Chatgpt.Answering {
				return nil, nil
			}
			input := strings.TrimSpace(that.textarea.Value())
			if input == "" {
				return nil, nil
			}
			that.ConvManager.Curr().AddQuestion(input)
			var cmds []tea.Cmd
			cmds = append(
				cmds,
				that.Chatgpt.Query(that.ConvManager.Curr().Config, that.ConvManager.Curr().GetContextMessages()),
			)
			// Start answer spinner
			that.spinning = true
			cmds = append(
				cmds, func() tea.Msg {
					return that.spinner.Tick()
				},
			)
			that.renderViewport()
			that.textarea.Reset()
			that.textarea.Blur()
			that.textarea.Placeholder = ""
			that.historyIdx = that.ConvManager.Curr().Len()
			return tea.Batch(cmds...), nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key: key.NewBinding(
			key.WithKeys("ctrl+m"),
			key.WithHelp("ctrl+m", "Create new conversation."),
		),
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			if that.Chatgpt.Answering {
				return nil, nil
			}
			that.err = nil
			// TODO change config when creating new conversation
			that.ConvManager.New(that.ConvManager.Conf.Conversation)
			that.renderViewport()
			that.historyIdx = 0
			return nil, nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key: key.NewBinding(
			key.WithKeys("ctrl+x"),
			key.WithHelp("ctrl+x", "forget context"),
		),
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			if that.Chatgpt.Answering {
				return nil, nil
			}
			that.err = nil
			that.ConvManager.Curr().ForgetContext()
			that.renderViewport()
			return nil, nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key: key.NewBinding(
			key.WithKeys("alt+r"),
			key.WithHelp("alt+r", "remove current conversation"),
		),
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			if that.Chatgpt.Answering {
				return nil, nil
			}
			that.err = nil
			that.ConvManager.RemoveCurr()
			that.renderViewport()
			that.historyIdx = that.ConvManager.Curr().Len()
			return nil, nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key: key.NewBinding(
			key.WithKeys("ctrl+left"),
			key.WithHelp("ctrl+left", "previous conversation"),
		),
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			if that.Chatgpt.Answering {
				return nil, nil
			}
			that.err = nil
			that.ConvManager.Prev()
			that.renderViewport()
			that.historyIdx = that.ConvManager.Curr().Len()
			return nil, nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key: key.NewBinding(
			key.WithKeys("ctrl+right"),
			key.WithHelp("ctrl+right", "next conversation"),
		),
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			if that.Chatgpt.Answering {
				return nil, nil
			}
			that.err = nil
			that.ConvManager.Next()
			that.renderViewport()
			that.historyIdx = that.ConvManager.Curr().Len()
			return nil, nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("ctrl+j", "multiline mode"),
		),
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			// if that.inputMode == InputModelSingleLine {
			// 	UseMultiLineInputMode(&m)
			// 	m.textarea.ShowLineNumbers = true
			// 	m.textarea.SetHeight(2)
			// 	m.viewport.Height = m.height - m.textarea.Height() - lipgloss.Height(m.RenderFooter())
			// } else {
			// 	UseSingleLineInputMode(&m)
			// 	m.textarea.ShowLineNumbers = false
			// 	m.textarea.SetHeight(1)
			// 	m.viewport.Height = m.height - m.textarea.Height() - lipgloss.Height(m.RenderFooter())
			// }
			// that.viewport.SetContent(that.RenderConversation(that.viewport.Width))
			return nil, nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key: key.NewBinding(
			key.WithKeys("ctrl+y"),
			key.WithHelp("ctrl+y", "copy last answer"),
		),
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			if that.Chatgpt.Answering || that.ConvManager.Curr().LastAnswer() == "" {
				return nil, nil
			}
			_ = clipboard.WriteAll(that.ConvManager.Curr().LastAnswer())
			return nil, nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "next question"),
		),
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			if that.Chatgpt.Answering {
				return nil, nil
			}
			idx := that.historyIdx + 1
			if idx >= that.ConvManager.Curr().Len() {
				that.historyIdx = that.ConvManager.Curr().Len()
				that.textarea.SetValue("")
			} else {
				that.textarea.SetValue(that.ConvManager.Curr().GetQuestion(idx))
				that.historyIdx = idx
			}
			return nil, nil
		},
	})
	kl = append(kl, &vtui.ShortcutKey{
		Name: that.ViewName,
		Key: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "previous question"),
		),
		Func: func(km tea.KeyMsg) (tea.Cmd, error) {
			if that.Chatgpt.Answering {
				return nil, nil
			}
			idx := that.historyIdx - 1
			if idx < 0 {
				idx = 0
			}
			q := that.ConvManager.Curr().GetQuestion(idx)
			that.textarea.SetValue(q)
			that.historyIdx = idx
			return nil, nil
		},
	})
	return kl
}

func (that *ChatgptView) Msgs() vtui.MessageList {
	ml := vtui.MessageList{}
	ml = append(ml, &vtui.Message{
		Name: chatgpt.VChatName,
		Func: func(msg *vtui.Message) (tea.Cmd, error) {
			if msg.Name == chatgpt.VChatName {
				switch msg.Type {
				case chatgpt.DeltaAnswerMsg:
					if c, ok := msg.Content.(string); ok {
						that.ConvManager.Curr().UpdatePending(c, false)
						that.err = nil
						that.renderViewport()
						return that.Chatgpt.Recv(), nil
					}
				case chatgpt.AnswerMsg:
					if c, ok := msg.Content.(string); ok {
						that.ConvManager.Curr().UpdatePending(c, true)
						that.spinning = false
						that.Chatgpt.Done()
						that.renderViewport()
						that.renderTextarea()
					}
				case chatgpt.SaveMsg:
					_ = that.ConvManager.Dump()
					return savePeriodically(), nil
				case chatgpt.ErrorMsg:
					if err, ok := msg.Content.(error); ok {
						if err == io.EOF {
							if that.ConvManager.Curr().PendingAnswer() == "" {
								that.err = errors.New("unexpected EOF, please try again")
							}
						} else {
							that.err = err
						}
						that.spinning = false
						that.ConvManager.Curr().UpdatePending("", true)
						that.Chatgpt.Done()
						that.renderViewport()
						that.renderTextarea()
					}
				default:
				}
			}
			return nil, nil
		},
	})
	return ml
}

func (that *ChatgptView) View() string {
	if that.Enabled {
		if that.width == 0 || that.height == 0 {
			return "Initializing..."
		}
		return lipgloss.JoinVertical(
			lipgloss.Left,
			that.viewport.View(),
			that.textarea.View(),
			that.RenderFooter(),
		)
	}
	return ""
}

func (that *ChatgptView) renderViewport() {
	that.viewport.SetContent(that.RenderConversation(that.viewport.Width))
	that.viewport.GotoBottom()
}

func (that *ChatgptView) renderTextarea() {
	that.textarea.Placeholder = "Send a message..."
	that.textarea.Focus()
}

var (
	SenderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	BotStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	ErrorStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1"))
	FooterStyle = lipgloss.NewStyle().
			Height(1).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("8")).
			Faint(true)
)

func (that *ChatgptView) RenderFooter() string {
	if that.err != nil {
		return FooterStyle.Render(ErrorStyle.Render(fmt.Sprintf("error: %v", that.err)))
	}

	var columns []string
	if that.spinning {
		columns = append(columns, that.spinner.View())
	} else {
		columns = append(columns, that.spinner.Spinner.Frames[0])
	}

	if that.ConvManager.Len() > 0 {
		columns = append(columns,
			fmt.Sprintf("%s %d/%d", vtui.ConversationIcon,
				that.ConvManager.Idx+1,
				that.ConvManager.Len()))
	}

	question := that.textarea.Value()
	if that.ConvManager.Curr().Len() > 0 || len(question) > 0 {
		tokens := that.ConvManager.Curr().GetContextTokens()
		if len(question) > 0 {
			tokens += chatgpt.CountTokens(that.ConvManager.Curr().Config.Model, question) + 5
		}
		columns = append(columns, fmt.Sprintf("%s %d", vtui.TokenIcon, tokens))
	}

	columns = append(columns, fmt.Sprintf("%s ctrl+h", vtui.HelpIcon))

	prompt := that.ConvManager.Curr().Config.Prompt
	prompt = fmt.Sprintf("%s %s", vtui.PromptIcon, prompt)
	columns = append(columns, prompt)

	totalWidth := lipgloss.Width(strings.Join(columns, ""))
	padding := (that.width - totalWidth) / (len(columns) - 1)
	if padding < 0 {
		padding = 2
	}

	if totalWidth+(len(columns)-1)*padding > that.width {
		remainingSpace := that.width - (lipgloss.Width(
			strings.Join(columns[:len(columns)-1], ""),
		) + (len(columns)-2)*padding + 3)
		columns[len(columns)-1] = columns[len(columns)-1][:remainingSpace] + "..."
	}

	footer := strings.Join(columns, strings.Repeat(" ", padding))
	footer = FooterStyle.Render(footer)
	return footer
}

func (that *ChatgptView) RenderConversation(maxWidth int) string {
	builder := &strings.Builder{}
	current := that.ConvManager.Curr()
	if current == nil {
		return ""
	}
	for _, m := range current.Forgotten {
		that.renderYou(m.Question, maxWidth, builder)
		that.renderChatgpt(m.Answer, maxWidth, builder)
	}
	if len(current.Forgotten) > 0 {
		builder.WriteString(lipgloss.NewStyle().PaddingLeft(5).Faint(true).Render("----- New Session -----"))
		builder.WriteString("\n")
	}
	for _, q := range current.Context {
		that.renderYou(q.Question, maxWidth, builder)
		that.renderChatgpt(q.Answer, maxWidth, builder)
	}
	if current.Pending != nil {
		that.renderYou(current.Pending.Question, maxWidth, builder)
		that.renderChatgpt(current.Pending.Answer, maxWidth, builder)
	}
	return builder.String()
}

func (that *ChatgptView) renderYou(content string, maxWidth int, builder *strings.Builder) {
	builder.WriteString(SenderStyle.Render("You: "))
	if utils.ContainsCJK(content) {
		content = wrap.String(content, maxWidth-5)
	} else {
		content = wordwrap.String(content, maxWidth-5)
	}
	content, _ = that.rd.Render(content)
	builder.WriteString(utils.EnsureTrailingNewline(content))
}

func (that *ChatgptView) renderChatgpt(content string, maxWidth int, builder *strings.Builder) {
	if content == "" {
		return
	}
	builder.WriteString(BotStyle.Render("ChatGPT: "))
	if utils.ContainsCJK(content) {
		content = wrap.String(content, maxWidth-5)
	} else {
		content = wordwrap.String(content, maxWidth-5)
	}
	content, _ = that.rd.Render(content)
	builder.WriteString(utils.EnsureTrailingNewline(content))
}
