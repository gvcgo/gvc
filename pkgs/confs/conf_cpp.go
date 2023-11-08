package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type CppConf struct {
	MsysInstallerUrl   string            `koanf:"msys_installer_url"`
	MsysMirrorUrls     map[string]string `koanf:"msys_mirror_urls"`
	CygwinInstallerUrl string            `koanf:"installer_url"`
	CygwinMirrorUrls   []string          `koanf:"mirror_url"`
	VCpkgUrl           string            `koanf:"vcpkg_url"`
	VCpkgToolUrl       string            `koanf:"vcpkg_tool_url"`
	WinVCpkgToolUrls   map[string]string `koanf:"win_vcpkg_tool_urls"`
	path               string
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
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *CppConf) Reset() {
	that.MsysInstallerUrl = "https://github.com/msys2/msys2-installer/releases/latest/msys2-x86_64-latest.exe"
	that.MsysMirrorUrls = map[string]string{
		"mirrorlist.msys":    "https://mirrors.tuna.tsinghua.edu.cn/msys2/msys/$arch/",
		"mirrorlist.mingw64": "https://mirrors.tuna.tsinghua.edu.cn/msys2/mingw/x86_64/",
		"mirrorlist.clang64": "https://mirrors.tuna.tsinghua.edu.cn/msys2/mingw/clang64/",
		"mirrorlist.mingw":   "https://mirrors.tuna.tsinghua.edu.cn/msys2/mingw/$repo/",
		"mirrorlist.mingw32": "https://mirrors.tuna.tsinghua.edu.cn/msys2/mingw/i686/",
		"mirrorlist.clang32": "https://mirrors.tuna.tsinghua.edu.cn/msys2/mingw/clang32/",
		"mirrorlist.ucrt64":  "https://mirrors.tuna.tsinghua.edu.cn/msys2/mingw/ucrt64/",
	}
	that.CygwinInstallerUrl = "https://www.cygwin.com/setup-x86_64.exe"
	that.CygwinMirrorUrls = []string{
		"https://mirrors.ustc.edu.cn/cygwin/",
		"https://mirrors.zju.edu.cn/cygwin/",
		"https://mirrors.tuna.tsinghua.edu.cn/cygwin/",
		"https://mirrors.aliyun.com/cygwin/",
	}
	that.VCpkgUrl = "https://github.com/microsoft/vcpkg-tool/archive/refs/heads/main.zip"
	that.VCpkgToolUrl = "https://github.com/microsoft/vcpkg-tool/archive/refs/heads/main.zip"
	that.WinVCpkgToolUrls = map[string]string{
		"arm64": "https://github.com/microsoft/vcpkg-tool/releases/latest/download/vcpkg-arm64.exe",
		"amd64": "https://github.com/microsoft/vcpkg-tool/releases/latest/download/vcpkg.exe",
	}
}
