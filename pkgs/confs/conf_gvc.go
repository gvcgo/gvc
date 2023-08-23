package confs

import (
	"os"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/gvc/pkgs/utils"
)

type GvcConf struct {
	GitlabUrls map[string]string `koanf:"gitlab_urls"`
	path       string
}

func NewGvcConf() (r *GvcConf) {
	r = &GvcConf{
		path: GVCBinTempDir,
	}
	r.setup()
	return
}

func (that *GvcConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			tui.PrintError(err)
		}
	}
}

func (that *GvcConf) Reset() {
	that.GitlabUrls = map[string]string{
		"windows_amd64": "https://gitlab.com/moqsien/gvc_resources/-/raw/main/gvc_windows-amd64.zip",
		"windows_arm64": "https://gitlab.com/moqsien/gvc_resources/-/raw/main/gvc_windows-arm64.zip",
		"linux_amd64":   "https://gitlab.com/moqsien/gvc_resources/-/raw/main/gvc_linux-amd64.zip",
		"linux_arm64":   "https://gitlab.com/moqsien/gvc_resources/-/raw/main/gvc_linux-arm64.zip",
		"darwin_amd64":  "https://gitlab.com/moqsien/gvc_resources/-/raw/main/gvc_darwin-amd64.zip",
		"darwin_arm64":  "https://gitlab.com/moqsien/gvc_resources/-/raw/main/gvc_darwin-arm64.zip",
	}
}
