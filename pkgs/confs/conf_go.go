package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type GoConf struct {
	CompilerUrls []string `koanf:"compiler_urls"`
	AliRepoUrl   string   `koanf:"ali_repo_url"`
	Proxies      []string `koanf:"proxies"`
	SearchUrl    string   `koanf:"search_url"`
	path         string
}

func NewGoConf() (r *GoConf) {
	r = &GoConf{
		path: GoFilesDir,
	}
	r.setup()
	return
}

func (that *GoConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *GoConf) Reset() {
	that.CompilerUrls = []string{
		"https://golang.google.cn/dl/",
		"https://go.dev/dl/",
		"https://studygolang.com/dl",
	}
	that.AliRepoUrl = "https://mirrors.aliyun.com/golang/"
	that.Proxies = []string{
		"https://goproxy.cn,direct",
		"https://repo.huaweicloud.com/repository/goproxy/,direct",
	}
	that.SearchUrl = `https://pkg.go.dev/search?limit=100&m=package&q=%s#more-results`
}
