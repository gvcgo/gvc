package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type GsudoConf struct {
	GitlabUrl string `koanf:"gitlab_url"`
	path      string
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
	that.GitlabUrl = "https://gitlab.com/moqsien/gvc_resources/-/raw/main/gsudo_portable.zip"
}
