package confs

import (
	"fmt"
	"os"

	"github.com/moqsien/gvc/pkgs/utils"
)

type GradleConf struct {
	OfficialUrl string `koanf:"official_url"`
	path        string
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
			fmt.Println("[mkdir Failed] ", that.path)
		}
	}
}

func (that *GradleConf) Reset() {
	that.OfficialUrl = "https://gradle.org/releases/"
}
