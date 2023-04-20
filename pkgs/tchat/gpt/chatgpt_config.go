package gpt

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/Asutorufa/yuhaiin/pkg/net/interfaces/proxy"
	"github.com/Asutorufa/yuhaiin/pkg/node/register"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/point"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/protocol"
	"github.com/gogf/gf/util/gconv"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	openai "github.com/sashabaranov/go-openai"
)

type Option struct {
	Value reflect.Value
	Type  string
	Field string
	KName string
}

func (that *Option) String() string {
	value := that.Value.Elem().FieldByName(that.Field)
	var r any
	switch that.Type {
	case "string":
		return value.String()
	case "bool":
		r = value.Bool()
	case "int":
		r = value.Int()
	case "float32", "float64", "float":
		r = value.Float()
	default:
		return ""
	}
	return gconv.String(r)
}

type ConversationConf struct {
	Prompt        string  `koanf:"prompt"`
	ContextLength int     `koanf:"context_length"`
	Model         string  `koanf:"model"`
	Stream        bool    `koanf:"stream"`
	Temperature   float32 `koanf:"temperature"`
	MaxTokens     int     `koanf:"max_tokens"`
}

type ChatGPTConf struct {
	Endpoint     string            `koanf:"endpoint"`
	APIKey       string            `koanf:"api_key"`
	APIType      string            `koanf:"api_type"`
	APIVersion   string            `koanf:"api_version"`
	Engine       string            `koanf:"engine"`
	OrgID        string            `koanf:"org_id"`
	ProxyType    string            `koanf:"proxy_type"`
	ProxyUri     string            `koanf:"proxy_uri"`
	ProxyTimeout int               `koanf:"proxy_timeout"`
	Conversation *ConversationConf `koanf:"conversation"`
	Prompts      map[string]string `koanf:"prompts"`
	optList      map[string]*Option
	optOrder     []string
	k            *koanf.Koanf
	parser       *yaml.YAML
	path         string
}

func NewChatGptConf() (cc *ChatGPTConf) {
	cc = &ChatGPTConf{
		Conversation: &ConversationConf{},
		k:            koanf.New("."),
		parser:       yaml.Parser(),
		path:         filepath.Join(config.ChatgptFilesDir, config.ChatgptConfigFileName),
		optList:      map[string]*Option{},
	}
	cc.setup()
	return
}

func (that *ChatGPTConf) GetOptList() (optList map[string]*Option) {
	return that.optList
}

func (that *ChatGPTConf) GetOptOrder() []string {
	return that.optOrder
}

func (that *ChatGPTConf) setup() {
	if ok, _ := utils.PathIsExist(config.ChatgptFilesDir); !ok {
		if err := os.MkdirAll(config.ChatgptFilesDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", that.path)
		}
	}
}

func (that *ChatGPTConf) Reset() {
	that.Endpoint = "https://api.openai.com/v1"
	that.APIKey = ""
	that.APIType = string(openai.APITypeOpenAI)
	that.Prompts = map[string]string{
		"default":    "You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible.",
		"translator": "I want you to act as an English translator, spelling corrector and improver. I will speak to you in any language and you will detect the language, translate it and answer in the corrected and improved version of my text, in English. I want you to replace my simplified A0-level words and sentences with more beautiful and elegant, upper level English words and sentences. The translation should be natural, easy to understand, and concise. Keep the meaning same, but make them more literary. I want you to only reply the correction, the improvements and nothing else, do not write explanations.",
		"shell":      "Return a one-line bash command with the functionality I will describe. Return ONLY the command ready to run in the terminal. The command should do the following:",
	}
	that.ProxyType = "socks5"
	that.ProxyUri = "localhost:2019"
	if that.Conversation != nil {
		that.Conversation.Model = openai.GPT3Dot5Turbo
		that.Conversation.Prompt = "default"
		that.Conversation.ContextLength = 6
		that.Conversation.Stream = true
		that.Conversation.Temperature = 0
		that.Conversation.MaxTokens = 1024
	}
}

func (that *ChatGPTConf) Restore() {
	if ok, _ := utils.PathIsExist(config.ChatgptFilesDir); !ok {
		os.MkdirAll(config.ChatgptFilesDir, os.ModePerm)
	}
	that.k.Load(structs.Provider(*that, "koanf"), nil)
	if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
		os.WriteFile(that.path, b, 0666)
	}
}

func (that *ChatGPTConf) Reload() {
	err := that.k.Load(file.Provider(that.path), that.parser)
	if err != nil {
		fmt.Println("[Config Load Failed] ", err)
		return
	}
	that.k.UnmarshalWithConf("", that, koanf.UnmarshalConf{Tag: "koanf"})
}

func (that *ChatGPTConf) GetOptions() {
	that.Reload()
	that.optOrder = []string{}
	cVal := reflect.ValueOf(that)
	valType := cVal.Type().Elem()
	for i := 0; i < valType.NumField(); i++ {
		tag := valType.Field(i).Tag
		koanfTag := tag.Get("koanf")
		if koanfTag != "" {
			switch koanfTag {
			case "conversation":
				if that.Conversation != nil {
					conVal := reflect.ValueOf(that.Conversation)
					conType := conVal.Type().Elem()
					for j := 0; j < conType.NumField(); j++ {
						ctag := conType.Field(j).Tag
						ckoanfTag := ctag.Get("koanf")
						that.optOrder = append(that.optOrder, ckoanfTag)
						if ckoanfTag != "" {
							that.optList[ckoanfTag] = &Option{
								Value: conVal,
								Type:  conType.Field(j).Type.Kind().String(),
								Field: conType.Field(j).Name,
								KName: ckoanfTag,
							}
						}
					}
				}
			default:
				if koanfTag != "" {
					typ := valType.Field(i).Type.Kind().String()
					if !strings.Contains(typ, "map") {
						that.optOrder = append(that.optOrder, koanfTag)
						that.optList[koanfTag] = &Option{
							Value: cVal,
							Type:  typ,
							Field: valType.Field(i).Name,
							KName: koanfTag,
						}
					}
				}
			}
		}
	}
}

func (that *ChatGPTConf) ShowOpts() []string {
	if len(that.optList) == 0 {
		that.GetOptions()
	}
	r := []string{}
	for k := range that.optList {
		r = append(r, k)
	}
	return r
}

func (that *ChatGPTConf) SetConfField(kName, value string) {
	that.Reload()
	if len(that.optList) == 0 {
		that.GetOptions()
	}
	kName = strings.ReplaceAll(kName, " ", "")
	value = strings.ReplaceAll(value, " ", "")
	f, ok := that.optList[kName]
	if !ok {
		return
	}
	val, fName, typeStr := f.Value, f.Field, f.Type

	if fName == "" {
		fmt.Println("Cannot find option: ", kName)
		return
	}
	field := val.Elem().FieldByName(fName)
	switch typeStr {
	case "string":
		if field.CanSet() {
			field.SetString(value)
		}
	case "int":
		if field.CanInt() {
			field.SetInt(gconv.Int64(value))
		}
	case "float32":
		if field.CanFloat() {
			field.SetFloat(gconv.Float64(value))
		}
	case "bool":
		if field.CanSet() {
			field.SetBool(gconv.Bool(value))
		}
	default:
		fmt.Println("Unsupported type!")
	}
	that.Restore()
}

func (that *ChatGPTConf) GetHttpClient() *http.Client {
	that.Reload()
	var (
		host string
		port int32
		u    *url.URL
		err  error
	)
	if u, err = url.Parse(that.ProxyUri); err == nil {
		host = u.Host
		port = gconv.Int32(u.Port())
	} else {
		return nil
	}
	if port != 0 && host != "" && (that.ProxyType == "socks5" || that.ProxyType == "") {
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
		pro, err := register.Dialer(node)
		if err != nil {
			fmt.Println("[Dialer error] ", err)
			return nil
		}
		if that.ProxyTimeout == 0 {
			that.ProxyTimeout = 300
		}
		r := &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					add, err := proxy.ParseAddress(proxy.PaseNetwork(network), addr)
					if err != nil {
						return nil, fmt.Errorf("parse address failed: %w", err)
					}
					add.WithContext(ctx)
					return pro.Conn(add)
				}}, Timeout: time.Duration(that.ProxyTimeout) * time.Second,
		}
		return r
	} else if port != 0 && host != "" {
		r := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(u),
			},
		}
		return r
	}
	return nil
}

func (that *ChatGPTConf) SearchPrompt(key string) string {
	prompt := that.Prompts[key]
	if prompt == "" {
		return key
	}
	return prompt
}
