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
		"windows":      "https://gitlab.com/moqsien/gvc_resources/-/raw/main/typst_x64_windows.zip",
		"linux_amd64":  "https://gitlab.com/moqsien/gvc_resources/-/raw/main/typst_x64_linux.tar.xz",
		"linux_arm64":  "https://gitlab.com/moqsien/gvc_resources/-/raw/main/typst_arm_linux.tar.xz",
		"darwin_amd64": "https://gitlab.com/moqsien/gvc_resources/-/raw/main/typst_x64_macos.tar.xz",
		"darwin_arm64": "https://gitlab.com/moqsien/gvc_resources/-/raw/main/typst_arm_macos.tar.xz",
	}
	that.GithubUrls = map[string]string{
		"windows":      "https://github.com/typst/typst/releases/latest/download/typst-x86_64-pc-windows-msvc.zip",
		"linux_amd64":  "https://github.com/typst/typst/releases/latest/download/typst-x86_64-unknown-linux-musl.tar.xz",
		"linux_arm64":  "https://github.com/typst/typst/releases/latest/download/typst-aarch64-unknown-linux-musl.tar.xz",
		"darwin_arm64": "https://github.com/typst/typst/releases/latest/download/typst-aarch64-apple-darwin.tar.xz",
		"darwin_amd64": "https://github.com/typst/typst/releases/latest/download/typst-x86_64-apple-darwin.tar.xz",
	}
}
