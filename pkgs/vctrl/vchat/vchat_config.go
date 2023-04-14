package vchat

import (
	"fmt"
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
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
	Conversation ConversationConf  `koanf:"conversation"`
	ProxyType    string            `koanf:"proxy_type"`
	ProxyPort    int               `koanf:"proxy_port"`
	k            *koanf.Koanf
	parser       *yaml.YAML
	path         string
}

func NewChatGptConf() (cc *ChatGPTConf) {
	cc = &ChatGPTConf{
		k:      koanf.New("."),
		parser: yaml.Parser(),
		path:   config.ChatgptFilesDir,
	}
	cc.setup()
	return
}

func (that *ChatGPTConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", that.path)
		}
	}
}
