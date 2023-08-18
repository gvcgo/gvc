package confs

import (
	"os"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/gvc/pkgs/utils"
)

type VlangConf struct {
	VlangGiteeUrls map[string]string `koanf:"vlang_gitee_url"`
	VlangUrls      map[string]string `koanf:"vlang_url"`
	path           string
}

func NewVlangConf() (r *VlangConf) {
	r = &VlangConf{
		path: VlangFilesDir,
	}
	r.setup()
	return
}

func (that *VlangConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			tui.PrintError(err)
		}
	}
}

func (that *VlangConf) Reset() {
	that.VlangGiteeUrls = map[string]string{
		utils.MacOS:   "https://gitlab.com/moqsien/gvc_resources/-/raw/main/vlang_macos.zip",
		utils.Linux:   "https://gitlab.com/moqsien/gvc_resources/-/raw/main/vlang_linux.zip",
		utils.Windows: "https://gitlab.com/moqsien/gvc_resources/-/raw/main/vlang_windows.zip",
	}
	that.VlangUrls = map[string]string{
		utils.MacOS:   "https://github.com/vlang/v/releases/latest/download/v_macos.zip",
		utils.Linux:   "https://github.com/vlang/v/releases/latest/download/v_linux.zip",
		utils.Windows: "https://github.com/vlang/v/releases/latest/download/v_windows.zip",
	}
}
