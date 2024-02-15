package confs

import (
	"os"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/gvc/pkgs/utils"
)

type VlangConf struct {
	VlangUrl    string `koanf,json:"vlang_url"`
	AnalyzerUrl string `koanf,json:"analyzer_url"`
	path        string
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
	that.VlangUrl = "https://github.com/vlang/v/releases/latest/"
	that.AnalyzerUrl = "https://github.com/v-analyzer/v-analyzer/releases/latest/"
}
