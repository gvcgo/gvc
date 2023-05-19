package confs

import (
	"os"

	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/tui"
)

type GradleConf struct {
	OfficialUrl      string `koanf:"official_url"`
	OfficialCheckUrl string `koanf:"official_check_url"`
	path             string
}

func NewGradleConf() (r *GradleConf) {
	r = &GradleConf{
		path: GradleRoot,
	}
	r.setup()
	return
}

func (that *GradleConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			tui.PrintError(err)
		}
	}
}

func (that *GradleConf) Reset() {
	that.OfficialUrl = "https://gradle.org/releases/"
	that.OfficialCheckUrl = "https://gradle.org/release-checksums/"
}
