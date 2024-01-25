package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type TypstConf struct {
	TypstUrl string `json,koanf:"typst_url"`
	path     string
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
	that.TypstUrl = "https://github.com/typst/typst/releases/latest/"
}
