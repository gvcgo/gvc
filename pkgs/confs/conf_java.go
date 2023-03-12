package confs

import (
	"fmt"
	"os"

	"github.com/moqsien/gvc/pkgs/utils"
)

type JavaConf struct {
	CompilerUrl string `koanf:"compiler_url"`
	path        string
}

func NewJavaConf() (r *JavaConf) {
	r = &JavaConf{
		path: JavaFilesDir,
	}
	r.setup()
	return
}

func (that *JavaConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", that.path)
		}
	}
}

func (that *JavaConf) Reset() {
	that.CompilerUrl = "https://www.oracle.com/java/technologies/downloads/"
}
