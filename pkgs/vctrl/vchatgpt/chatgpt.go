package vchatgpt

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/Asutorufa/yuhaiin/pkg/net/interfaces/proxy"
	"github.com/Asutorufa/yuhaiin/pkg/node/register"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/point"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/protocol"
	config "github.com/moqsien/gvc/pkgs/confs"
	openai "github.com/sashabaranov/go-openai"
)

type ChatGpt struct {
	Client *openai.Client
	Conf   *config.GVConfig
}

func NewChatGpt() (cg *ChatGpt) {
	cg = &ChatGpt{
		Conf: config.New(),
	}
	return
}

func (that *ChatGpt) SetAppKey(key string) {
	that.Conf.Reload()
	that.Conf.Chatgpt.AppKey = key
	that.Conf.Restore()
	that.Conf.Reload()
}

func (that *ChatGpt) SetProxyPort(p int) {
	that.Conf.Reload()
	that.Conf.Chatgpt.LocalProxyPort = p
	that.Conf.Restore()
	that.Conf.Reload()
}

func (that *ChatGpt) getClientWithProxy() {
	if that.Conf.Chatgpt.LocalProxyPort == 0 {
		fmt.Println("Please set your local socks5 port!")
		return
	}
	node := &point.Point{
		Protocols: []*protocol.Protocol{
			{
				Protocol: &protocol.Protocol_Simple{
					Simple: &protocol.Simple{
						Host:             "127.0.0.1",
						Port:             int32(that.Conf.Chatgpt.LocalProxyPort),
						PacketConnDirect: true,
					},
				},
			},
			{
				Protocol: &protocol.Protocol_Socks5{
					Socks5: &protocol.Socks5{},
				},
			},
		},
	}
	pro, err := register.Dialer(node)

	if err != nil {
		fmt.Println("[Dialer error] ", err)
		return
	}
	t := that.Conf.Chatgpt.ProxyTimeout
	if t == 0 {
		t = 120
	}
	t = 600
	c := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				add, err := proxy.ParseAddress(proxy.PaseNetwork(network), addr)
				if err != nil {
					return nil, fmt.Errorf("parse address failed: %w", err)
				}
				add.WithContext(ctx)
				return pro.Conn(add)
			}}, Timeout: time.Duration(t) * time.Second,
	}

	appkey := that.Conf.Chatgpt.AppKey
	oc := openai.DefaultConfig(appkey)
	oc.HTTPClient = c
	that.Client = openai.NewClientWithConfig(oc)
}

func (that *ChatGpt) GetClient(withProxy, force bool) {
	if that.Conf.Chatgpt.AppKey == "" {
		fmt.Println("Please set your openai AppKey!")
		return
	}
	if that.Client == nil || force {
		if withProxy {
			that.getClientWithProxy()
		} else {
			that.Client = openai.NewClient(that.Conf.Chatgpt.AppKey)
		}
	}
}

func (that *ChatGpt) Chat(content string) {
	that.GetClient(true, false)
	count := 0

	for count < 5 {
		count++
		resp, err := that.Client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: content,
					},
				},
			},
		)

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			continue
		}
		fmt.Println(resp.Choices[0].Message.Content)
		break
	}
}
