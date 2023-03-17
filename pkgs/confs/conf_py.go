package confs

import (
	"fmt"
	"os"

	"github.com/moqsien/gvc/pkgs/utils"
)

type PyConf struct {
	WinAmd64    string   `koanf:"win_amd64"`
	WinArm64    string   `koanf:"win_arm64"`
	PyenvUnix   string   `koanf:"pyenv_unix"`
	PyenvWin    string   `koanf:"pyenv_win"`
	PypiProxies []string `koanf:"pypi_proxies"`
	path        string
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
			fmt.Println("[mkdir Failed] ", that.path)
		}
	}
}

func (that *PyConf) Reset() {
	that.WinAmd64 = "https://gitee.com/moqsien/gvc/releases/download/v1/portable-amd64.zip"
	that.WinArm64 = "https://gitee.com/moqsien/gvc/releases/download/v1/portable-arm64.zip"
	that.PyenvWin = "https://github.com/pyenv-win/pyenv-win/archive/master.zip"
	that.PyenvUnix = "https://github.com/pyenv/pyenv/archive/refs/heads/master.zip"
	that.PypiProxies = []string{
		"https://pypi.tuna.tsinghua.edu.cn/simple/",
		"https://pypi.mirrors.ustc.edu.cn/simple/",
		"http://mirrors.aliyun.com/pypi/simple/",
		"http://pypi.douban.com/simple/",
	}
}
