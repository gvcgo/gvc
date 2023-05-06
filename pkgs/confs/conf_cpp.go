package confs

import (
	"fmt"
	"os"

	"github.com/moqsien/gvc/pkgs/utils"
)

type CppConf struct {
	MsysInstallerUrl string            `koanf:"msys_installer_url"`
	MsysMirrorUrls   map[string]string `koanf:"msys_mirror_urls"`
	VCpkgUrl         string            `koanf:"vcpkg_url"`
	path             string
}

func NewCppConf() (r *CppConf) {
	r = &CppConf{
		path: CppFilesDir,
	}
	r.setup()
	return
}

func (that *CppConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", that.path)
		}
	}
}

func (that *CppConf) Reset() {
	that.MsysInstallerUrl = "https://mirrors.tuna.tsinghua.edu.cn/msys2/distrib/x86_64/"
	that.MsysMirrorUrls = map[string]string{
		"mirrorlist.msys":    "https://mirrors.tuna.tsinghua.edu.cn/msys2/msys/$arch/",
		"mirrorlist.mingw64": "https://mirrors.tuna.tsinghua.edu.cn/msys2/mingw/x86_64/",
		"mirrorlist.clang64": "https://mirrors.tuna.tsinghua.edu.cn/msys2/mingw/clang64/",
		"mirrorlist.mingw":   "https://mirrors.tuna.tsinghua.edu.cn/msys2/mingw/$repo/",
		"mirrorlist.mingw32": "https://mirrors.tuna.tsinghua.edu.cn/msys2/mingw/i686/",
		"mirrorlist.clang32": "https://mirrors.tuna.tsinghua.edu.cn/msys2/mingw/clang32/",
		"mirrorlist.ucrt64":  "https://mirrors.tuna.tsinghua.edu.cn/msys2/mingw/ucrt64/",
	}
	that.VCpkgUrl = "https://github.com/microsoft/vcpkg/archive/refs/heads/master.zip"
}
