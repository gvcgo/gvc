package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type VlangConf struct {
	VlangGitlabUrls    map[string]string `koanf:"vlang_gitlab_url"`
	VlangUrls          map[string]string `koanf:"vlang_url"`
	AnalyzerUrls       map[string]string `koanf:"analyzer_url"`
	AnalyzerGitlabUrls map[string]string `koanf:"analyzer_gitlab_url"`
	path               string
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
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *VlangConf) Reset() {
	that.VlangGitlabUrls = map[string]string{
		utils.MacOS:   "https://gitlab.com/moqsien/gvc_resources/-/raw/main/vlang_macos.zip",
		utils.Linux:   "https://gitlab.com/moqsien/gvc_resources/-/raw/main/vlang_linux.zip",
		utils.Windows: "https://gitlab.com/moqsien/gvc_resources/-/raw/main/vlang_windows.zip",
	}
	that.VlangUrls = map[string]string{
		utils.MacOS:   "https://github.com/vlang/v/releases/latest/download/v_macos.zip",
		utils.Linux:   "https://github.com/vlang/v/releases/latest/download/v_linux.zip",
		utils.Windows: "https://github.com/vlang/v/releases/latest/download/v_windows.zip",
	}
	that.AnalyzerUrls = map[string]string{
		utils.Windows:  "https://github.com/v-analyzer/v-analyzer/releases/latest/download/v-analyzer-windows-x86_64.zip",
		utils.Linux:    "https://github.com/v-analyzer/v-analyzer/releases/latest/download/v-analyzer-linux-x86_64.zip",
		"darwin_amd64": "https://github.com/v-analyzer/v-analyzer/releases/latest/download/v-analyzer-darwin-x86_64.zip",
		"darwin_arm64": "https://github.com/v-analyzer/v-analyzer/releases/latest/download/v-analyzer-darwin-arm64.zip",
	}
	that.AnalyzerGitlabUrls = map[string]string{
		utils.Windows:  "https://gitlab.com/moqsien/gvc_resources/-/raw/main/v_analyzer_windows_x86_64.zip",
		utils.Linux:    "https://gitlab.com/moqsien/gvc_resources/-/raw/main/v_analyzer_linux_x86_64.zip",
		"darwin_amd64": "https://gitlab.com/moqsien/gvc_resources/-/raw/main/v_analyzer_darwin_x86_64.zip",
		"darwin_arm64": "https://gitlab.com/moqsien/gvc_resources/-/raw/main/v_analyzer_darwin_arm64.zip",
	}
}
