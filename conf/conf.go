package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
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
	return filepath.Join(GetGVCWorkDir(), DefaultGVConfigFileName)
}

type GVConfig struct {
	GitUserName   string `json:"git_username"`
	GitToken      string `json:"git_token"`
	GiteeUserName string `json:"gitee_username"`
	GiteeToken    string `json:"gitee_token"`
	Password      string `json:"password"`
	PicRepo       string `json:"pic_repo"`
	BackupRepo    string `json:"backup_repo"`
	LocalProxy    string `json:"local_proxy"`
	ReverseProxy  string `json:"reverse_proxy"`
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

/*
Get/Set config values.
*/
func (c *GVConfig) GetGitUserName() string {
	c.Load()
	if c.GitUserName == "" {
		fmt.Println(gprint.CyanStr(`Please set your github username:`))
		var username string
		fmt.Scanln(&username)
		c.GitUserName = strings.TrimSpace(username)
		c.Save()
	}
	return c.GitUserName
}

func (c *GVConfig) GetGitToken() string {
	c.Load()
	if c.GitToken == "" {
		fmt.Println(gprint.CyanStr(`Please set your github token:`))
		var token string
		fmt.Scanln(&token)
		c.GitToken = strings.TrimSpace(token)
		c.Save()
	}
	return c.GitToken
}

func (c *GVConfig) GetGiteeUserName() string {
	c.Load()
	if c.GiteeUserName == "" {
		fmt.Println(gprint.CyanStr(`Please set your gitee username:`))
		var username string
		fmt.Scanln(&username)
		c.GiteeUserName = strings.TrimSpace(username)
		c.Save()
	}
	return c.GiteeUserName
}

func (c *GVConfig) GetGiteeToken() string {
	c.Load()
	if c.GiteeToken == "" {
		fmt.Println(gprint.CyanStr(`Please set your gitee token:`))
		var token string
		fmt.Scanln(&token)
		c.GiteeToken = strings.TrimSpace(token)
		c.Save()
	}
	return c.GiteeToken
}

func (c *GVConfig) GetPassword() string {
	c.Load()
	if c.Password == "" {
		fmt.Println(gprint.CyanStr(`Please set your password for encryting files:`))
		var password string
		fmt.Scanln(&password)
		c.Password = strings.TrimSpace(password)
		c.Save()
	}
	return c.Password
}

func (c *GVConfig) GetPicRepo() string {
	c.Load()
	if c.PicRepo == "" {
		fmt.Println(gprint.CyanStr(`Please set your picture repo name:`))
		var repo string
		fmt.Scanln(&repo)
		c.PicRepo = strings.TrimSpace(repo)
		c.Save()
	}
	return c.PicRepo
}

func (c *GVConfig) GetConfPath() string {
	c.Load()
	if c.BackupRepo == "" {
		fmt.Println(gprint.CyanStr(`Please set your backup repo name:`))
		var repo string
		fmt.Scanln(&repo)
		c.BackupRepo = strings.TrimSpace(repo)
		c.Save()
	}
	return c.BackupRepo
}

func (c *GVConfig) GetLocalProxy() string {
	c.Load()
	if c.LocalProxy == "" {
		fmt.Println(gprint.CyanStr(`Please set your local proxy:`))
		var proxy string
		fmt.Scanln(&proxy)
		c.LocalProxy = strings.TrimSpace(proxy)
		c.Save()
	}
	return c.LocalProxy
}

func (c *GVConfig) GetReverseProxy() string {
	c.Load()
	if c.ReverseProxy == "" {
		fmt.Println(gprint.CyanStr(`Please set your reverse proxy:`))
		var proxy string
		fmt.Scanln(&proxy)
		c.ReverseProxy = strings.TrimSpace(proxy)
		c.Save()
	}
	return c.ReverseProxy
}
