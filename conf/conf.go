package conf

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gvcgo/goutils/pkgs/gutils"
)

const (
	DefaultGVConfigFileName string = "gvc.conf"
)

func GetGVCWorkDir() string {
	homeDir, _ := os.UserHomeDir()
	r := filepath.Join(homeDir, ".gvc")
	if ok, _ := gutils.PathIsExist(r); !ok {
		os.MkdirAll(r, os.ModePerm)
	}
	return r
}

func GetConfPath() string {
	return filepath.Join(GetGVCWorkDir(), ".gvc", DefaultGVConfigFileName)
}

type GVConfig struct {
	GitToken   string `json:"git_token"`
	GiteeToken string `json:"gitee_token"`
	Password   string `json:"password"`
	PicRepo    string `json:"pic_repo"`
	LocalProxy string `json:"local_proxy"`
}

func NewGVConfig() *GVConfig {
	cfg := &GVConfig{}
	cfg.Load()
	return cfg
}

func (c *GVConfig) Load() {
	if content, err := os.ReadFile(GetConfPath()); err == nil {
		json.Unmarshal(content, c)
	}
}

func (c *GVConfig) Save() {
	content, _ := json.MarshalIndent(c, "", "    ")
	os.WriteFile(GetConfPath(), content, os.ModePerm)
}
