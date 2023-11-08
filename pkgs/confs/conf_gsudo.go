package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type GsudoConf struct {
	Url  string `koanf:"github_url"`
	path string
}

func NewGsudoConf() (r *GsudoConf) {
	r = &GsudoConf{
		path: GsudoFilePath,
	}
	r.setup()
	return
}

func (that *GsudoConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *GsudoConf) Reset() {
	that.Url = "https://github.com/gerardog/gsudo/releases/latest/download/gsudo.portable.zip"
}
