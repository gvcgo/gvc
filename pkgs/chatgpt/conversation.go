package chatgpt

import "github.com/sashabaranov/go-openai"

type QandA struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Conversation struct {
	Conf        *ConvsationConf `json:"conf"`
	Prompt      string          `json:"prompt"`
	QandAList   []*QandA        `json:"qa_list"`
	Current     *QandA          `json:"current"`
	bot         *ChatBot
	isAnswering bool
}

func NewConversation(cnf *ConvsationConf, prompt string, bot *ChatBot) *Conversation {
	return &Conversation{
		Conf:   cnf,
		Prompt: prompt,
		bot:    bot,
	}
}

func (that *Conversation) IsAnswering() bool {
	return that.isAnswering
}

func (that *Conversation) GetConversationContext() (r []openai.ChatCompletionMessage) {
	if that.Prompt != "" {
		r = append(r, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: that.Prompt,
		})
	}
	for _, qa := range that.QandAList {
		r = append(r, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: qa.Question,
		})
		r = append(r, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: qa.Answer,
		})
	}
	if that.Current != nil {
		r = append(r, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: that.Current.Question,
		})
	}
	return
}

func (that *Conversation) Ask(question string) {
	if !that.isAnswering {
		if that.Current != nil {
			that.QandAList = append(that.QandAList, that.Current)
			that.Current = nil
		}
		that.Current = &QandA{
			Question: question,
		}
		collector := make(chan string)
		go that.bot.Send(that.Conf, that.GetConversationContext(), collector)
		that.isAnswering = true
		for ans := range collector {
			that.Current.Answer += ans
		}
		that.isAnswering = false
	}
}

func (that *Conversation) RenderConversation() (r string) {
	return
}
