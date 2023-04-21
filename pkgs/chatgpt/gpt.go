package chatgpt

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Asutorufa/yuhaiin/pkg/net/interfaces/proxy"
	"github.com/Asutorufa/yuhaiin/pkg/node/register"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/point"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/protocol"
	"github.com/sashabaranov/go-openai"
)

type ChatBot struct {
	Client *openai.Client
}

func NewBot(conf *ChatConf) (cb *ChatBot) {
	cb = &ChatBot{}
	oconf := openai.DefaultConfig(conf.APIKey)
	oconf.BaseURL = conf.BaseUrl
	oconf.OrgID = conf.OrgId
	if conf.APIType != openai.APITypeOpenAI {
		oconf.APIType = conf.APIType
		oconf.APIVersion = conf.APIVersion
		oconf.Engine = conf.Engine
	}
	oconf.HTTPClient = cb.GetHttpClient(conf.ProxyUri, conf.Timeout)
	cb.Client = openai.NewClientWithConfig(oconf)
	return
}

func (that *ChatBot) parseProxy(proxyUri string) (scheme, host string, port int) {
	if u, err := url.Parse(proxyUri); err == nil {
		scheme = u.Scheme
		host = u.Hostname()
		port, _ = strconv.Atoi(u.Port())
		if port == 0 {
			port = 80
		}
	}
	return
}

func (that *ChatBot) getSocks5Client(host string, port int32, timeout int) (h *http.Client) {
	node := &point.Point{
		Protocols: []*protocol.Protocol{
			{
				Protocol: &protocol.Protocol_Simple{
					Simple: &protocol.Simple{
						Host:             host,
						Port:             port,
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
	if proxy_result, err := register.Dialer(node); err == nil {
		h = &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					add, err := proxy.ParseAddress(proxy.PaseNetwork(network), addr)
					if err != nil {
						return nil, fmt.Errorf("parse address failed: %w", err)
					}
					add.WithContext(ctx)
					return proxy_result.Conn(add)
				}},
			Timeout: time.Second * time.Duration(timeout),
		}
	} else {
		fmt.Println("[Socks5 Dialer error] ", err)
		h = &http.Client{}
	}
	return
}

func (that *ChatBot) GetHttpClient(proxyUri string, timeout int) (h *http.Client) {
	scheme, host, port := that.parseProxy(proxyUri)
	if host != "" && port != 0 {
		switch scheme {
		case "http", "https":
			if pUri, err := url.Parse(proxyUri); err == nil {
				h = &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(pUri),
					},
					Timeout: time.Second * time.Duration(timeout),
				}
			}
		case "socks5":
			h = that.getSocks5Client(host, int32(port), timeout)
		default:
			panic("unsupported proxy type")
		}
	}
	if h == nil {
		h = &http.Client{}
	}
	return
}

func (that *ChatBot) Send(conf *ConvsationConf, msgs []openai.ChatCompletionMessage, collector chan string) error {
	request := openai.ChatCompletionRequest{
		Model:       conf.Model,
		Messages:    msgs,
		MaxTokens:   conf.MaxTokens,
		Temperature: conf.Temperature,
		N:           1,
	}
	if conf.EnableStream {
		stream, err := that.Client.CreateChatCompletionStream(context.Background(), request)
		if err != nil {
			return err
		}
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					close(collector)
					return nil
				}
				return err
			}
			content := resp.Choices[0].Delta.Content
			collector <- content
		}
	} else {
		resp, err := that.Client.CreateChatCompletion(context.Background(), request)
		if err != nil {
			return err
		}
		content := resp.Choices[0].Message.Content
		collector <- content
	}
	return nil
}

func Run() {
	cnf := &ChatConf{
		APIKey:   "",
		BaseUrl:  "https://api.openai.com/v1",
		ProxyUri: "socks5://localhost:2019",
		Timeout:  180,
	}
	cb := NewBot(cnf)

	cf := &ConvsationConf{
		Model:        "gpt-3.5-turbo",
		MaxTokens:    1024,
		EnableStream: true,
	}

	result := make(chan string, 10)
	go cb.Send(cf, []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "用golang写一个快排",
		},
	}, result)
	mList := []string{}
	for a := range result {
		mList = append(mList, a)
	}
	fmt.Println(strings.Join(mList, ""))
}
