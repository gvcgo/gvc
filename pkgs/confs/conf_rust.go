package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type RustConf struct {
	UrlUnix      string `koanf:"url_unix"`
	FileNameUnix string `koanf:"filename_unix"`
	UrlWin       string `koanf:"url_win"`
	FileNameWin  string `koanf:"filename_win"`
	DistServer   string `koanf:"RUSTUP_DIST_SERVER"`
	UpdateRoot   string `koanf:"RUSTUP_UPDATE_ROOT"`
	path         string
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
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *RustConf) Reset() {
	that.UrlWin = "https://static.rust-lang.org/rustup/dist/i686-pc-windows-gnu/rustup-init.exe"
	that.FileNameWin = "rustup-init.exe"
	that.UrlUnix = "https://sh.rustup.rs"
	that.FileNameUnix = "rustup-init.sh"
	that.DistServer = "https://mirrors.ustc.edu.cn/rust-static"
	that.UpdateRoot = "https://mirrors.ustc.edu.cn/rust-static/rustup"
}
