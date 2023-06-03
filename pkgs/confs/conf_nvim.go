package confs

import (
	"os"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/gvc/pkgs/utils"
)

type NUrl struct {
	Url  string `koanf:"url"`
	Name string `koanf:"name"`
	Ext  string `koanf:"ext"`
}

type NVimConf struct {
	Urls        map[string]*NUrl `koanf:"urls"`
	ChecksumUrl string           `koanf:"checksum_url"`
	PluginsUrl  string           `koanf:"plugins_url"`
	GithubProxy string           `koanf:"github_proxy"`
	path        string
}

func NewNVimConf() (r *NVimConf) {
	r = &NVimConf{
		path: NVimFileDir,
	}
	r.setup()
	return
}

func (that *NVimConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			tui.PrintError(err)
		}
	}
}

func (that *NVimConf) Reset() {
	that.Urls = map[string]*NUrl{
		"darwin": {
			"https://gitee.com/moqsien/gvc/releases/download/v1/nvim-macos.tar.gz",
			"nvim-macos",
			".tar.gz",
		},
		"linux": {
			"https://gitee.com/moqsien/gvc/releases/download/v1/nvim-linux64.tar.gz",
			"nvim-linux64",
			".tar.gz",
		},
		"windows": {
			"https://gitee.com/moqsien/gvc/releases/download/v1/nvim-win64.zip",
			"nvim-win64",
			".zip",
		},
	}
	that.ChecksumUrl = "https://gitee.com/moqsien/gvc/releases/download/v1/nvim-sha256.txt"
	that.PluginsUrl = "https://gitee.com/moqsien/gvc/releases/download/v1/nvim-plugins.zip"
	that.GithubProxy = "https://ghproxy.com/"
}
