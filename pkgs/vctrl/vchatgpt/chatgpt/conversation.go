package chatgpt

import (
	"log"
	"os"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	openai "github.com/sashabaranov/go-openai"
)

type CManager interface {
	GetChatConf() *ChatGPTConf
}

type QuesAnsw struct {
	Question string `koanf:"question"`
	Answer   string `koanf:"answer"`
}

type Conversation struct {
	Manager       CManager
	contextTokens int
	Config        *ConversationConf `koanf:"config"`
	Forgotten     []*QuesAnsw       `koanf:"forgotten,omitempty"`
	Context       []*QuesAnsw       `koanf:"context,omitempty"`
	Pending       *QuesAnsw         `koanf:"pending,omitempty"`
}

func (that *Conversation) AddQuestion(q string) {
	that.Pending = &QuesAnsw{Question: q}
	that.contextTokens = 0
}

func (that *Conversation) UpdatePending(ans string, done bool) {
	if that.Pending == nil {
		return
	}
	that.Pending.Answer += ans
	if done {
		that.Context = append(that.Context, that.Pending)
		that.contextTokens = 0
		if len(that.Context) > that.Config.ContextLength {
			that.Forgotten = append(that.Forgotten, that.Context[0])
			that.Context = that.Context[1:]
		}
		that.Pending = nil
	}
}

func (that *Conversation) ForgetContext() {
	that.Forgotten = append(that.Forgotten, that.Context...)
	that.Context = nil
	that.contextTokens = 0
}

func (that *Conversation) PendingAnswer() string {
	if that.Pending == nil {
		return ""
	}
	return that.Pending.Answer
}

func (that *Conversation) LastAnswer() string {
	if len(that.Context) == 0 {
		return ""
	}
	return that.Context[len(that.Context)-1].Answer
}

func (that *Conversation) Len() int {
	l := len(that.Forgotten) + len(that.Context)
	if that.Pending != nil {
		l++
	}
	return l
}

func (that *Conversation) GetQuestion(idx int) string {
	if idx < 0 || idx >= that.Len() {
		return ""
	}
	if idx < len(that.Forgotten) {
		return that.Forgotten[idx].Question
	}
	return that.Context[idx-len(that.Forgotten)].Question
}

func (that *Conversation) GetContextMessages() []openai.ChatCompletionMessage {
	messages := make([]openai.ChatCompletionMessage, 0, 2*len(that.Context)+2)

	if conf := that.Manager.GetChatConf(); conf != nil {
		messages = append(
			messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: conf.SearchPrompt(that.Config.Prompt),
			},
		)
	}

	for _, c := range that.Context {
		messages = append(
			messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: c.Question,
			},
		)
		messages = append(
			messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: c.Answer,
			},
		)
	}

	if that.Pending != nil {
		messages = append(
			messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: that.Pending.Question,
			},
		)
	}
	return messages
}

func (that *Conversation) GetContextTokens() int {
	if that.contextTokens == 0 {
		that.contextTokens = CountMessagesTokens(that.Config.Model, that.GetContextMessages())
	}
	return that.contextTokens
}

// Conversation Manager
type ConvManager struct {
	filePath      string
	Conf          *ChatGPTConf
	Conversations []*Conversation `koanf:"conversations"`
	Idx           int             `koanf:"last_idx"`
	k             *koanf.Koanf
	parser        *utils.JSON
}

func NewConvManager(conf *ChatGPTConf) *ConvManager {
	cm := &ConvManager{
		Conf:   conf,
		Idx:    -1,
		k:      koanf.New("."),
		parser: utils.NewJsonParser(),
	}
	cm.filePath = filepath.Join(config.ChatgptFilesDir, config.ChatgptConversationFileName)
	err := cm.Load()
	if err != nil {
		log.Println("Failed to load history:", err)
	}
	return cm
}

func (that *ConvManager) GetChatConf() *ChatGPTConf {
	return that.Conf
}

func (that *ConvManager) Dump() error {
	if ok, _ := utils.PathIsExist(config.ChatgptFilesDir); !ok {
		os.MkdirAll(config.ChatgptFilesDir, os.ModePerm)
	}
	that.k.Load(structs.Provider(*that, "koanf"), nil)
	if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
		os.WriteFile(that.filePath, b, 0666)
	} else {
		return err
	}
	return nil
}

func (that *ConvManager) Load() error {
	if err := that.k.Load(file.Provider(that.filePath), that.parser); err != nil {
		return err
	}
	return that.k.UnmarshalWithConf("", that, koanf.UnmarshalConf{Tag: "koanf"})
}

// Create a new Conversation
func (that *ConvManager) New(conf *ConversationConf) *Conversation {
	c := &Conversation{
		Manager: that,
		Config:  conf,
	}
	that.Conversations = append(that.Conversations, c)
	that.Idx = len(that.Conversations) - 1
	return c
}

func (that *ConvManager) RemoveCurr() {
	if len(that.Conversations) == 0 {
		return
	}
	that.Conversations = append(that.Conversations[:that.Idx], that.Conversations[that.Idx+1:]...)
	if that.Idx >= len(that.Conversations) {
		that.Idx = len(that.Conversations) - 1
	}
}

func (that *ConvManager) Len() int {
	return len(that.Conversations)
}

func (that *ConvManager) Curr() *Conversation {
	if len(that.Conversations) == 0 {
		// create initial conversation using default config
		return that.New(that.Conf.Conversation)
	}
	return that.Conversations[that.Idx]
}

func (that *ConvManager) Prev() *Conversation {
	if len(that.Conversations) == 0 {
		return nil
	}
	that.Idx--
	if that.Idx < 0 {
		that.Idx = 0 // dont wrap around
	}
	return that.Conversations[that.Idx]
}

func (that *ConvManager) Next() *Conversation {
	if len(that.Conversations) == 0 {
		return nil
	}
	that.Idx++
	if that.Idx >= len(that.Conversations) {
		that.Idx = len(that.Conversations) - 1 // dont wrap around
	}
	return that.Conversations[that.Idx]
}
