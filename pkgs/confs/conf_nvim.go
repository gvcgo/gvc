package confs

import (
	"os"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/gvc/pkgs/utils"
)

// type NUrl struct {
// 	Url  string `koanf:"url"`
// 	Name string `koanf:"name"`
// 	Ext  string `koanf:"ext"`
// }

type NVimConf struct {
	NvimUrl       string `koanf,json:"nvim_url"`
	NeovideUrl    string `koanf,json:"neovide_url"`
	TreeSitterUrl string `koanf,json:"treesitter_url"`
	FzFUrl        string `koanf,json:"fzf_url"`
	FdUrl         string `koanf,josn:"fd_url"`
	RipgrepUrl    string `koanf,json:"ripgrep_url"`
	GNvimUrl      string `koanf,json:"gnvim_url"`
	GlowUrl       string `koanf,json:"glow_url"`
	AstroNvimUrl  string `koanf,json:"astro_nvim_url"`
	PluginsUrl    string `koanf:"plugins_url"`
	GithubProxy   string `koanf:"github_proxy"`
	path          string
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
	that.NvimUrl = "https://github.com/neovim/neovim/releases/latest/"
	that.NeovideUrl = "https://github.com/neovide/neovide/releases/latest/"
	that.TreeSitterUrl = "https://github.com/tree-sitter/tree-sitter/releases/latest/"
	that.FzFUrl = "https://github.com/junegunn/fzf/releases/latest/"
	that.FdUrl = "https://github.com/sharkdp/fd/releases/latest/"
	that.RipgrepUrl = "https://github.com/BurntSushi/ripgrep/releases/latest/"
	that.GlowUrl = "https://github.com/charmbracelet/glow/releases/latest/"
	that.GNvimUrl = "git@github.com:gvcgo/gnvim.git"
	that.AstroNvimUrl = "https://codeload.github.com/AstroNvim/AstroNvim/zip/refs/heads/main"
	that.PluginsUrl = "https://gitlab.com/moqsien/gvc_resources/uploads/753afef9d38f8f6224d221770d25c9a3/nvim-plugins.zip"
	that.GithubProxy = "https://ghproxy.com/"
}
