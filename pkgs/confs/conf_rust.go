package confs

import (
	"fmt"
	"os"

	"github.com/moqsien/gvc/pkgs/utils"
)

type RustConf struct {
	UrlUnix    string `koanf:"url_unix"`
	UrlWin     string `koanf:"url_win"`
	DistServer string `koanf:"RUSTUP_DIST_SERVER"`
	UpdateRoot string `koanf:"RUSTUP_UPDATE_ROOT"`
	path       string
}

func NewRustConf() (r *RustConf) {
	r = &RustConf{
		path: RustFilesDir,
	}
	r.setup()
	return
}

func (that *RustConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", that.path)
		}
	}
}

func (that *RustConf) Reset() {
	that.UrlWin = "https://static.rust-lang.org/rustup/dist/i686-pc-windows-gnu/rustup-init.exe"
	that.UrlUnix = "https://sh.rustup.rs"
	that.DistServer = "https://mirrors.ustc.edu.cn/rust-static"
	that.UpdateRoot = "https://mirrors.ustc.edu.cn/rust-static/rustup"
}
