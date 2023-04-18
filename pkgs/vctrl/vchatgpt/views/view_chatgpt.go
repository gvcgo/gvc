package views

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/chatgpt"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/vtui"
)

var textareaKeys = textarea.KeyMap{
	CharacterForward:        key.NewBinding(key.WithKeys("right", "ctrl+f")),
	CharacterBackward:       key.NewBinding(key.WithKeys("left", "ctrl+b")),
	WordForward:             key.NewBinding(key.WithKeys("alt+right", "alt+f")),
	WordBackward:            key.NewBinding(key.WithKeys("alt+left", "alt+b")),
	LineNext:                key.NewBinding(key.WithKeys("down")),
	LinePrevious:            key.NewBinding(key.WithKeys("up")),
	DeleteWordBackward:      key.NewBinding(key.WithKeys("alt+backspace", "ctrl+w")),
	DeleteWordForward:       key.NewBinding(key.WithKeys("alt+delete", "alt+d")),
	DeleteAfterCursor:       key.NewBinding(key.WithKeys("ctrl+k")),
	DeleteBeforeCursor:      key.NewBinding(key.WithKeys("ctrl+u")),
	InsertNewline:           key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "insert new line")),
	DeleteCharacterBackward: key.NewBinding(key.WithKeys("backspace")),
	DeleteCharacterForward:  key.NewBinding(key.WithKeys("delete")),
	LineStart:               key.NewBinding(key.WithKeys("home", "ctrl+a")),
	LineEnd:                 key.NewBinding(key.WithKeys("end", "ctrl+e")),
	Paste:                   key.NewBinding(key.WithKeys("ctrl+v", "alt+v"), key.WithHelp("ctrl+v", "paste")),
	InputBegin:              key.NewBinding(key.WithKeys("alt+<", "ctrl+home")),
	InputEnd:                key.NewBinding(key.WithKeys("alt+>", "ctrl+end")),

	CapitalizeWordForward: key.NewBinding(key.WithKeys("alt+c")),
	LowercaseWordForward:  key.NewBinding(key.WithKeys("alt+l")),
	UppercaseWordForward:  key.NewBinding(key.WithKeys("alt+u")),

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

	// shortcut keys
	cv.viewport.KeyMap = viewportKeys
	cv.textarea.KeyMap = textareaKeys
	return
}

func (that *ChatgptView) ExtraCmdHandlers() []vtui.CmdHandler {
	return []vtui.CmdHandler{}
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
	return ""
}
