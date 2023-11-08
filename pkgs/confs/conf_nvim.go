package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type NUrl struct {
	Url  string `koanf:"url"`
	Name string `koanf:"name"`
	Ext  string `koanf:"ext"`
}

type NVimConf struct {
	Urls        map[string]*NUrl `koanf:"urls"`
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
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *NVimConf) Reset() {
	that.Urls = map[string]*NUrl{
		"darwin": {
			"https://github.com/neovim/neovim/releases/latest/download/nvim-macos.tar.gz",
			"nvim_macos",
			".tar.gz",
		},
		"linux": {
			"https://github.com/neovim/neovim/releases/latest/download/nvim-linux64.tar.gz",
			"nvim_linux64",
			".tar.gz",
		},
		"windows": {
			"https://github.com/neovim/neovim/releases/latest/download/nvim-win64.zip",
			"nvim_win64",
			".zip",
		},
	}
	that.PluginsUrl = "https://gitlab.com/moqsien/gvc_resources/uploads/753afef9d38f8f6224d221770d25c9a3/nvim-plugins.zip"
	that.GithubProxy = "https://ghproxy.com/"
}
