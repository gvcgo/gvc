package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type CodeConf struct {
	StableUrl      string   `koanf:"stable_url"`
	CdnUrl         string   `koanf:"cdn_url"`
	DownloadUrl    string   `koanf:"download_url"`
	ExtIdentifiers []string `koanf:"ext_identifiers"`
	path           string
}

func NewCodeConf() (r *CodeConf) {
	r = &CodeConf{
		path: CodeFileDir,
	}
	r.setup()
	return
}

func (that *CodeConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *CodeConf) Reset() {
	that.StableUrl = "az764295.vo.msecnd.net"
	that.CdnUrl = "vscode.cdn.azure.cn"
	that.DownloadUrl = "https://code.visualstudio.com/sha?build=stable"
	that.ExtIdentifiers = []string{
		"moqsien.easynotes",
		"doggy8088.go-extension-pack",
		"galkowskit.go-interface-annotations",
		"liuchao.go-struct-tag",
		"tabnine.tabnine-vscode",
		"gruntfuggly.todo-tree",
		"zxh404.vscode-proto3",
		"premparihar.gotestexplorer",
		"ms-python.python",
		"ms-python.vscode-pylance",
		"donjayamanne.python-environment-manager",
		"alefragnani.project-manager",
		"yzhang.markdown-all-in-one",
		"mhutchie.git-graph",
		"asvetliakov.vscode-neovim",
		"ms-ceintl.vscode-language-pack-zh-hans",
		"bracketpaircolordlw.bracket-pair-color-dlw",
		"rust-lang.rust-analyzer",
		"vue.volar",
		"joe-re.sql-language-server",
		"akamud.vscode-theme-onedark",
		"pkief.material-icon-theme",
	}
}
