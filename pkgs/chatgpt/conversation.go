package chatgpt

type QandA struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Conversation struct {
	Prompt    string   `json:"prompt"`
	QandAList []*QandA `json:"qa_list"`
	Current   *QandA   `json:"current"`
	Bot       *ChatBot
}
