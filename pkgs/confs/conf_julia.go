package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type JuliaConf struct {
	VersionUrl         string   `koanf:"version_url"`
	VersionUrlOfficial string   `koanf:"version_url_official"`
	MirrorUrls         []string `koanf:"mirror_urls"`
	BaseUrl            string   `koanf:"base_url"`
	PkgServer          string   `koanf:"pkg_server"`
	path               string
}

func NewJuliaConf() (r *JuliaConf) {
	r = &JuliaConf{
		path: JuliaFilesDir,
	}
	r.setup()
	return
}

func (that *JuliaConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *JuliaConf) Reset() {
	that.VersionUrl = "https://mirrors.tuna.tsinghua.edu.cn/julia-releases/bin/versions.json"
	that.VersionUrlOfficial = "https://julialang-s3.julialang.org/bin/versions.json"
	that.MirrorUrls = []string{
		"https://mirrors.ustc.edu.cn/julia-releases/bin/versions.json",
		"https://mirrors.nju.edu.cn/julia-releases/bin/versions.json",
	}
	that.BaseUrl = "https://mirrors.tuna.tsinghua.edu.cn/julia-releases/bin"
	that.PkgServer = "https://mirrors.tuna.tsinghua.edu.cn/julia"
}
