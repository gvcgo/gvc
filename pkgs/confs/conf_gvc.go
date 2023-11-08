package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type GvcConf struct {
	Urls map[string]string `koanf:"github_urls"`
	path string
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
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *GvcConf) Reset() {
	that.Urls = map[string]string{
		"windows_amd64": "https://github.com/moqsien/gvc/releases/latest/download/gvc_windows-amd64.zip",
		"windows_arm64": "https://github.com/moqsien/gvc/releases/latest/download/gvc_windows-arm64.zip",
		"linux_amd64":   "https://github.com/moqsien/gvc/releases/latest/download/gvc_linux-amd64.zip",
		"linux_arm64":   "https://github.com/moqsien/gvc/releases/latest/download/gvc_linux-arm64.zip",
		"darwin_amd64":  "https://github.com/moqsien/gvc/releases/latest/download/gvc_darwin-amd64.zip",
		"darwin_arm64":  "https://github.com/moqsien/gvc/releases/latest/download/gvc_darwin-arm64.zip",
	}
}
