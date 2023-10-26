package confs

import (
	"os"
	"path/filepath"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type GPTConf struct {
	WorkDir string
}

func NewGPTConf() (gconf *GPTConf) {
	gconf = &GPTConf{
		WorkDir: filepath.Join(GVCInstallDir, "gpt_files"),
	}
	gconf.setup()
	return
}

func (that *GPTConf) setup() {
	if ok, _ := utils.PathIsExist(that.WorkDir); !ok {
		if err := os.MkdirAll(that.WorkDir, os.ModePerm); err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *GPTConf) Reset() {
	that.WorkDir = filepath.Join(GVCInstallDir, "gpt_files")
}
