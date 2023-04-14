package vchat

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gogf/gf/util/gconv"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	openai "github.com/sashabaranov/go-openai"
)

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
	APIType      openai.APIType    `koanf:"api_type"`
	APIVersion   string            `koanf:"api_version"`
	Engine       string            `koanf:"engine"`
	OrgID        string            `koanf:"org_id"`
	Prompts      map[string]string `koanf:"prompts"`
	Conversation *ConversationConf `koanf:"conversation"`
	ProxyType    string            `koanf:"proxy_type"`
	ProxyPort    int               `koanf:"proxy_port"`
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
	}
	cc.setup()
	return
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
	that.Prompts = map[string]string{
		"default":    "You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible.",
		"translator": "I want you to act as an English translator, spelling corrector and improver. I will speak to you in any language and you will detect the language, translate it and answer in the corrected and improved version of my text, in English. I want you to replace my simplified A0-level words and sentences with more beautiful and elegant, upper level English words and sentences. The translation should be natural, easy to understand, and concise. Keep the meaning same, but make them more literary. I want you to only reply the correction, the improvements and nothing else, do not write explanations.",
		"shell":      "Return a one-line bash command with the functionality I will describe. Return ONLY the command ready to run in the terminal. The command should do the following:",
	}
	that.ProxyType = "socks5"
	that.ProxyPort = 2019
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

func (that *ChatGPTConf) ShowConfigOpts(kNames ...string) (val reflect.Value, fName string, fType string) {
	var kName string
	if len(kNames) > 0 {
		kName = kNames[0]
	}
	cVal := reflect.ValueOf(that)
	valType := cVal.Type().Elem()
	idx := 0
	for i := 0; i < valType.NumField(); i++ {
		name := valType.Field(i).Name
		tag := valType.Field(i).Tag
		koanfTag := tag.Get("koanf")

		if koanfTag != "" {
			switch name {
			case "Conversation":
				if that.Conversation != nil {
					conVal := reflect.ValueOf(that.Conversation)
					conType := conVal.Type().Elem()
					for j := 0; j < conType.NumField(); j++ {
						ctag := conType.Field(j).Tag
						ckoanfTag := ctag.Get("koanf")
						idx++
						t := conType.Field(j).Type.Kind().String()
						if kName == "" {
							fmt.Printf("%d. %s [%s]\n", idx, ckoanfTag, t)
						} else if kName == ckoanfTag {
							return conVal, conType.Field(j).Name, t
						}
					}
				}
			default:
				idx++
				t := valType.Field(i).Type.Kind().String()
				if kName == "" {
					fmt.Printf("%d. %s [%s]\n", idx, koanfTag, t)
				} else if kName == koanfTag {
					return cVal, valType.Field(i).Name, t
				}
			}
		}
	}
	return
}

func (that *ChatGPTConf) SetConfField(kName, value string) {
	that.Reload()
	kName = strings.ReplaceAll(kName, " ", "")
	value = strings.ReplaceAll(value, " ", "")
	val, fName, typeStr := that.ShowConfigOpts(kName)
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
