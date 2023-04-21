package chatgpt

import "github.com/sashabaranov/go-openai"

type ChatConf struct {
	APIKey     string         `json,koanf:"api_key"`
	BaseUrl    string         `json,koanf:"base_url"`
	APIType    openai.APIType `json,koanf:"api_type"`
	APIVersion string         `json,koanf:"api_version"`
	Engine     string         `json,koanf:"engine"`
	OrgId      string         `json,koanf:"org_id"`
	ProxyUri   string         `json,koanf:"proxy_uri"`
	Timeout    int            `json,koanf:"timeout"`
}

type ConvsationConf struct {
	Model        string  `json,koanf:"model"`
	MaxTokens    int     `json,koanf:"max_token"`
	Temperature  float32 `json,koanf:"temperature"`
	EnableStream bool    `json,koanf:"enable_stream"`
}
