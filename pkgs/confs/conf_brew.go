package confs

import (
	"fmt"
	"os"

	"github.com/moqsien/gvc/pkgs/utils"
)

type HomebrewConf struct {
	ShellScriptUrl string `koanf:"shell_script_url"`
	path           string
}

func NewHomebrewConf() (r *HomebrewConf) {
	r = &HomebrewConf{
		path: HomebrewFileDir,
	}
	r.setup()
	return
}

func (that *HomebrewConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", that.path)
		}
	}
}

func (that *HomebrewConf) Reset() {
	that.ShellScriptUrl = "https://gitee.com/moqsien/gvc/raw/master/homebrew.sh"
}
