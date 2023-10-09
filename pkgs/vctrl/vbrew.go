package vctrl

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/gtea/selector"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Homebrew struct {
	Conf    *config.GVConfig
	envs    *utils.EnvsHandler
	fetcher *request.Fetcher
}

func NewHomebrew() (hb *Homebrew) {
	hb = &Homebrew{
		Conf:    config.New(),
		fetcher: request.NewFetcher(),
		envs:    utils.NewEnvsHandler(),
	}
	hb.envs.SetWinWorkDir(config.GVCDir)
	return
}

func (that *Homebrew) getShellScript() string {
	fPath := filepath.Join(config.HomebrewFileDir, "homebrew.sh")
	if ok, _ := utils.PathIsExist(fPath); !ok {
		that.fetcher.Url = that.Conf.Homebrew.ShellScriptUrl
		if size := that.fetcher.GetAndSaveFile(fPath); size > 0 {
			return fPath
		}
		return ""
	}
	return fPath
}

func (that *Homebrew) SetEnv() {
	itemList := selector.NewItemList()
	itemList.Add("from mirrors.tuna.tsinghua.edu.cn", that.Conf.Homebrew.TsingHua)
	itemList.Add("form mirrors.ustc.edu.cn", that.Conf.Homebrew.USTC)
	sel := selector.NewSelector(
		itemList,
		selector.WidthEnableMulti(false),
		selector.WithEnbleInfinite(true),
		selector.WithWidth(40),
		selector.WithHeight(10),
		selector.WithTitle("Choose a homebrew mirror"),
	)
	sel.Run()
	value := sel.Value()[0]
	envMap := value.(map[string]string)
	if len(envMap) > 0 {
		envars := fmt.Sprintf(utils.HOMEbrewEnv,
			envMap["HOMEBREW_API_DOMAIN"],
			envMap["HOMEBREW_BOTTLE_DOMAIN"],
			envMap["HOMEBREW_BREW_GIT_REMOTE"],
			envMap["HOMEBREW_CORE_GIT_REMOTE"],
			envMap["HOMEBREW_PIP_INDEX_URL"])
		that.envs.UpdateSub(utils.SUB_BREW, envars)
	}
}

func (that *Homebrew) Install() {
	if runtime.GOOS != utils.Windows {
		script := that.getShellScript()
		if _, err := utils.ExecuteSysCommand(false, "sh", script); err != nil {
			gprint.PrintError("%+v", err)
			return
		}
		that.SetEnv()
	} else {
		gprint.PrintError("Homebrew does not support Windows.")
	}
}
