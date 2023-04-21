package chatgpt

import "github.com/sashabaranov/go-openai"

type ChatConf struct {
	APIKey     string         `json:"api_key"`
	BaseUrl    string         `json:"base_url"`
	APIType    openai.APIType `json:"api_type"`
	APIVersion string         `json:"api_version"`
	Engine     string         `json:"engine"`
	OrgId      string         `json:"org_id"`
	ProxyUri   string         `json:"proxy_uri"`
	Timeout    int            `json:"timeout"`
}

type ConvsationConf struct {
	Model        string  `json:"model"`
	MaxTokens    int     `json:"max_token"`
	Temperature  float32 `json:"temperature"`
	EnableStream bool    `json:"enable_stream"`
}
