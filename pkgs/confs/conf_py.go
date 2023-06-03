package confs

import (
	"os"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/gvc/pkgs/utils"
)

type PyConf struct {
	WinAmd64       string   `koanf:"win_amd64"`
	WinArm64       string   `koanf:"win_arm64"`
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
			tui.PrintError(err)
		}
	}
}

func (that *PyConf) Reset() {
	that.WinAmd64 = "https://gitee.com/moqsien/gvc/releases/download/v1/portable-amd64.zip"
	that.WinArm64 = "https://gitee.com/moqsien/gvc/releases/download/v1/portable-arm64.zip"
	that.PyenvWin = "https://gitee.com/moqsien/gvc/releases/download/v1/pyenv-win.zip"
	that.PyenvWinNeeded = "https://gitee.com/moqsien/gvc/releases/download/v1/pyenv_win_needed.zip"
	that.PyenvUnix = "https://gitee.com/moqsien/gvc/releases/download/v1/pyenv-unix.zip"
	that.PyenvReadline = []string{
		"https://gitee.com/moqsien/gvc/releases/download/v1/readline-8.2.tar.gz",
		"https://gitee.com/moqsien/gvc/releases/download/v1/readline-8.1.tar.gz",
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
