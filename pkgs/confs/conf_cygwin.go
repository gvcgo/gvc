package confs

import (
	"fmt"
	"os"

	"github.com/moqsien/gvc/pkgs/utils"
)

type CygwinConf struct {
	InstallerUrl string   `koanf:"installer_url"`
	MirrorUrls   []string `koanf:"mirror_url"`
	PackagesUrl  string   `koanf:"packages_url"`
	path         string
}

func NewCygwinConf() (r *CygwinConf) {
	r = &CygwinConf{
		path: CygwinFilesDir,
	}
	r.setup()
	return
}

func (that *CygwinConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", that.path)
		}
	}
}

func (that *CygwinConf) Reset() {
	that.InstallerUrl = "https://gitee.com/moqsien/gvc/releases/download/v1/cygwin-installer.exe"
	that.PackagesUrl = "https://gitee.com/moqsien/gvc/releases/download/v1/cygwin-packages.yml"
	that.MirrorUrls = []string{
		"https://mirrors.ustc.edu.cn/cygwin/",
		"https://mirrors.zju.edu.cn/cygwin/",
		"https://mirrors.tuna.tsinghua.edu.cn/cygwin/",
		"https://mirrors.aliyun.com/cygwin/",
	}
}
