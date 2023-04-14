package vchat

import (
	openai "github.com/sashabaranov/go-openai"
)

type VChat struct {
	Conf      *ChatGPTConf
	Client    *openai.Client
	Stream    *openai.ChatCompletionStream
	Answering bool
}
