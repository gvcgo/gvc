package chatgpt

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/avast/retry-go"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/vtui"
	openai "github.com/sashabaranov/go-openai"
)

const VChatName string = "vchatgpt"

type VChat struct {
	Conf      *ChatGPTConf
	Client    *openai.Client
	Stream    *openai.ChatCompletionStream
	Answering bool
}

func NewVChat(conf *ChatGPTConf) (vc *VChat) {
	vc = &VChat{
		Conf: conf,
	}
	vc.initVChat()
	return
}

func (that *VChat) initVChat() {
	if that.Conf == nil {
		return
	}
	config := openai.DefaultConfig(that.Conf.APIKey)
	config.OrgID = that.Conf.OrgID
	if that.Conf.Endpoint != "" {
		config.BaseURL = that.Conf.Endpoint
	}
	if that.Conf.APIType != string(openai.APITypeOpenAI) {
		config.APIType = openai.APIType(that.Conf.APIType)
		config.APIVersion = that.Conf.APIVersion
		config.Engine = that.Conf.Engine
	}
	hclient := that.Conf.GetHttpClient()
	if hclient != nil {
		config.HTTPClient = hclient
	}
	that.Client = openai.NewClientWithConfig(config)
}

func (that *VChat) Ask(conf *ConversationConf, question string, out io.Writer) error {
	req := openai.ChatCompletionRequest{
		Model: conf.Model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: that.Conf.SearchPrompt(conf.Prompt)},
			{Role: openai.ChatMessageRoleUser, Content: question},
		},
		MaxTokens:   conf.MaxTokens,
		Temperature: conf.Temperature,
		N:           1,
	}

	if conf.Stream {
		stream, err := that.Client.CreateChatCompletionStream(context.Background(), req)
		if err != nil {
			return err
		}
		defer stream.Close()
		for {
			resp, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					_, _ = fmt.Fprintln(out)
					break
				}
				return err
			}
			content := resp.Choices[0].Delta.Content
			_, _ = fmt.Fprint(out, content)
		}
	} else {
		resp, err := that.Client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			return err
		}
		content := resp.Choices[0].Message.Content
		_, _ = fmt.Fprintln(out, content)
	}

	return nil
}

var (
	DeltaAnswerMsg vtui.MsgType = "chatgpt_deltaanswer"
	AnswerMsg      vtui.MsgType = "chatgpt_answer"
	SaveMsg        vtui.MsgType = "chatgpt_save"
	ErrorMsg       vtui.MsgType = "chatgpt_error"
)

func (that *VChat) Query(conf *ConversationConf, mList []openai.ChatCompletionMessage) (r tea.Cmd) {
	that.Answering = true
	r = func() (msg tea.Msg) {
		executor := func() error {
			req := openai.ChatCompletionRequest{
				Model:       conf.Model,
				Messages:    mList,
				MaxTokens:   conf.MaxTokens,
				Temperature: conf.Temperature,
				N:           1,
			}
			if conf.Stream {
				stream, err := that.Client.CreateChatCompletionStream(context.Background(), req)
				that.Stream = stream
				if err != nil {
					return vtui.ErrMsg(err)
				}
				resp, err := stream.Recv()
				if err != nil {
					return err
				}
				content := resp.Choices[0].Delta.Content
				msg = vtui.NewMessage(VChatName, DeltaAnswerMsg, content)
			} else {
				resp, err := that.Client.CreateChatCompletion(context.Background(), req)
				if err != nil {
					return vtui.ErrMsg(err)
				}
				content := resp.Choices[0].Message.Content
				msg = vtui.NewMessage(VChatName, AnswerMsg, content)
			}
			return nil
		}
		if err := retry.Do(executor, retry.Attempts(3), retry.LastErrorOnly(true)); err != nil {
			msg = vtui.NewMessage(VChatName, ErrorMsg, err)
		}
		return
	}
	return
}

// func (that *VChat) Send(conf *ConversationConf, messages []openai.ChatCompletionMessage) (r tea.Cmd) {
// 	that.Answering = true
// 	r = func() (msg tea.Msg) {
// 		executor := func() error {
// 			req := openai.ChatCompletionRequest{
// 				Model:       conf.Model,
// 				Messages:    messages,
// 				MaxTokens:   conf.MaxTokens,
// 				Temperature: conf.Temperature,
// 				N:           1,
// 			}

// 			if conf.Stream {
// 				stream, err := that.Client.CreateChatCompletionStream(context.Background(), req)
// 				that.Stream = stream
// 				if err != nil {
// 					return vtui.ErrMsg(err)
// 				}
// 				resp, err := stream.Recv()
// 				if err != nil {
// 					return err
// 				}
// 				content := resp.Choices[0].Delta.Content
// 				msg = vtui.DeltaAnswerMsg(content)
// 			} else {
// 				resp, err := that.Client.CreateChatCompletion(context.Background(), req)
// 				if err != nil {
// 					return vtui.ErrMsg(err)
// 				}
// 				content := resp.Choices[0].Message.Content
// 				msg = vtui.AnswerMsg(content)
// 			}
// 			return nil
// 		}

// 		if err := retry.Do(executor, retry.Attempts(3), retry.LastErrorOnly(true)); err != nil {
// 			return vtui.ErrMsg(err)
// 		}
// 		return
// 	}
// 	return
// }

func (that *VChat) Recv() (r tea.Cmd) {
	r = func() (msg tea.Msg) {
		if resp, err := that.Stream.Recv(); err != nil {
			msg = vtui.NewMessage(VChatName, ErrorMsg, err)
		} else {
			msg = vtui.NewMessage(VChatName, DeltaAnswerMsg, resp.Choices[0].Delta.Content)
		}
		return
	}
	return
}

// func (that *VChat) Receive() tea.Cmd {
// 	return func() tea.Msg {
// 		resp, err := that.Stream.Recv()
// 		if err != nil {
// 			return vtui.ErrMsg(err)
// 		}
// 		content := resp.Choices[0].Delta.Content
// 		return vtui.DeltaAnswerMsg(content)
// 	}
// }

func (that *VChat) Done() {
	if that.Stream != nil {
		that.Stream.Close()
	}
	that.Stream = nil
	that.Answering = false
}
