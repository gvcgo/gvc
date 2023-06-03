package confs

import (
	"os"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/gvc/pkgs/utils"
)

type TypstConf struct {
	GithubUrls map[string]string `koanf:"github_urls"`
	GiteeUrls  map[string]string `koanf:"gitee_urls"`
	path       string
}

func NewTypstConf() (r *TypstConf) {
	r = &TypstConf{
		path: TypstFilesDir,
	}
	r.setup()
	return
}

func (that *TypstConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			tui.PrintError(err)
		}
	}
}

func (that *TypstConf) Reset() {
	that.GiteeUrls = map[string]string{
		utils.Windows: "https://gitee.com/moqsien/gvc/releases/download/v1/typst-windows.zip",
		utils.MacOS:   "https://gitee.com/moqsien/gvc/releases/download/v1/typst-darwin.zip",
		utils.Linux:   "https://gitee.com/moqsien/gvc/releases/download/v1/typst-linux.zip",
	}
	that.GithubUrls = map[string]string{
		utils.Windows: "https://github.com/typst/typst/releases/download/v0.1.0/typst-x86_64-pc-windows-msvc.zip",
		utils.MacOS:   "https://github.com/typst/typst/releases/download/v0.1.0/typst-x86_64-apple-darwin.tar.gz",
		utils.Linux:   "https://github.com/typst/typst/releases/download/v0.1.0/typst-x86_64-unknown-linux-gnu.tar.gz",
	}
}
