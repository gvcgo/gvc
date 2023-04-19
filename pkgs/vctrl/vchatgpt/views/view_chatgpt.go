package views

import (
	"fmt"
	"strings"

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
	inputMode   InputMode
	err         error
	historyIdx  int
	width       int
	height      int
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
	return
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
		that.handleWindowSize,
		that.handleSpinnerTick,
	}
}

func (that *ChatgptView) Keys() vtui.KeyList {
	kl := vtui.KeyList{}
	return kl
}

func (that *ChatgptView) Msgs() vtui.MessageList {
	ml := vtui.MessageList{}
	return ml
}

func (that *ChatgptView) View() string {
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
