package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type HomebrewConf struct {
	ShellScriptUrl string            `koanf:"shell_script_url"`
	TsingHua       map[string]string `koanf:"tsing_hua"`
	USTC           map[string]string `koanf:"ustc"`
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
			gprint.PrintError("%+v", err)
		}
	}
}

/*
export HOMEBREW_BREW_GIT_REMOTE="https://mirrors.ustc.edu.cn/brew.git"
export HOMEBREW_CORE_GIT_REMOTE="https://mirrors.ustc.edu.cn/homebrew-core.git"
export HOMEBREW_BOTTLE_DOMAIN="https://mirrors.ustc.edu.cn/homebrew-bottles"
export HOMEBREW_API_DOMAIN="https://mirrors.ustc.edu.cn/homebrew-bottles/api"
*/
func (that *HomebrewConf) Reset() {
	that.ShellScriptUrl = "https://gitee.com/moqsien/gvc/raw/master/homebrew.sh"
	that.TsingHua = map[string]string{
		"HOMEBREW_API_DOMAIN":      "https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles/api",
		"HOMEBREW_BOTTLE_DOMAIN":   "https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles",
		"HOMEBREW_BREW_GIT_REMOTE": "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/brew.git",
		"HOMEBREW_CORE_GIT_REMOTE": "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/homebrew-core.git",
		"HOMEBREW_PIP_INDEX_URL":   "https://pypi.tuna.tsinghua.edu.cn/simple",
	}
	that.USTC = map[string]string{
		"HOMEBREW_API_DOMAIN":      "https://mirrors.ustc.edu.cn/homebrew-bottles/api",
		"HOMEBREW_BOTTLE_DOMAIN":   "https://mirrors.ustc.edu.cn/homebrew-bottles",
		"HOMEBREW_BREW_GIT_REMOTE": "https://mirrors.ustc.edu.cn/brew.git",
		"HOMEBREW_CORE_GIT_REMOTE": "https://mirrors.ustc.edu.cn/homebrew-core.git",
		"HOMEBREW_PIP_INDEX_URL":   "https://mirrors.ustc.edu.cn/pypi/web/simple",
	}
}
