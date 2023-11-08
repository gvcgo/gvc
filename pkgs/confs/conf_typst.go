package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type TypstConf struct {
	GithubUrls map[string]string `koanf:"github_urls"`
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
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *TypstConf) Reset() {
	that.GithubUrls = map[string]string{
		"windows":      "https://github.com/typst/typst/releases/latest/download/typst-x86_64-pc-windows-msvc.zip",
		"linux_amd64":  "https://github.com/typst/typst/releases/latest/download/typst-x86_64-unknown-linux-musl.tar.xz",
		"linux_arm64":  "https://github.com/typst/typst/releases/latest/download/typst-aarch64-unknown-linux-musl.tar.xz",
		"darwin_arm64": "https://github.com/typst/typst/releases/latest/download/typst-aarch64-apple-darwin.tar.xz",
		"darwin_amd64": "https://github.com/typst/typst/releases/latest/download/typst-x86_64-apple-darwin.tar.xz",
	}
}
