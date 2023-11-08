package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type PyConf struct {
	PyenvUnix      string   `koanf:"pyenv_unix"`
	PyenvReadline  []string `koanf:"pyenv_readline"`
	PyenvWin       string   `koanf:"pyenv_win"`
	PyenvWinNeeded string   `koanf:"pyenv_win_needed"`
	PypiProxies    []string `koanf:"pypi_proxies"`
	PyBuildUrls    []string `koanf:"python_build_urls"`
	PyBuildUrl     string   `koanf:"python_build_url"`
	path           string
}

func NewPyConf() (r *PyConf) {
	r = &PyConf{
		path: PythonFilesDir,
	}
	r.setup()
	return
}

func (that *PyConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *PyConf) Reset() {
	that.PyenvWin = "https://github.com/pyenv-win/pyenv-win/archive/refs/heads/master.zip"
	that.PyenvWinNeeded = "https://gitlab.com/moqsien/gvc_resources/uploads/45d0ad242f9abb45b5a09b9634d3be73/pyenv_win_needed.zip"
	that.PyenvUnix = "https://github.com/pyenv/pyenv/archive/refs/heads/master.zip"
	that.PyenvReadline = []string{
		"https://gitlab.com/moqsien/gvc_resources/uploads/06845bbd8f73ce5c24e1d4f5761829a1/readline-8.1.tar.gz",
		"https://gitlab.com/moqsien/gvc_resources/uploads/5ae9516bd13038839b0aa102dada0a14/readline-8.2.tar.gz",
	}
	that.PypiProxies = []string{
		"https://pypi.tuna.tsinghua.edu.cn/simple/",
		"https://pypi.mirrors.ustc.edu.cn/simple/",
		"http://mirrors.aliyun.com/pypi/simple/",
		"http://pypi.douban.com/simple/",
	}
	that.PyBuildUrls = []string{
		"https://registry.npmmirror.com/-/binary/python/",
		"https://npm.taobao.org/mirrors/python/",
	}
	that.PyBuildUrl = "https://npm.taobao.org/mirrors/python/"
}
