package vctrl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/query"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/tui"
)

type Homebrew struct {
	Conf    *config.GVConfig
	envs    *utils.EnvsHandler
	fetcher *query.Fetcher
}

func NewHomebrew() (hb *Homebrew) {
	hb = &Homebrew{
		Conf:    config.New(),
		fetcher: query.NewFetcher(),
		envs:    utils.NewEnvsHandler(),
	}
	hb.envs.SetWinWorkDir(config.GVCWorkDir)
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

// TODO: pterm options
func (that *Homebrew) SetEnv() {
	mirror := ""
	fmt.Println("Choose a Mirror Site in China:")
	fmt.Println("1) TsingHua")
	fmt.Println("2) USTC")
	fmt.Scan(&mirror)
	switch mirror {
	case "1":
		envMap := that.Conf.Homebrew.TsingHua
		envars := fmt.Sprintf(utils.HOMEbrewEnv,
			envMap["HOMEBREW_API_DOMAIN"],
			envMap["HOMEBREW_BOTTLE_DOMAIN"],
			envMap["HOMEBREW_BREW_GIT_REMOTE"],
			envMap["HOMEBREW_CORE_GIT_REMOTE"],
			envMap["HOMEBREW_PIP_INDEX_URL"])
		that.envs.UpdateSub(utils.SUB_BREW, envars)
	case "2":
		envMap := that.Conf.Homebrew.USTC
		envars := fmt.Sprintf(utils.HOMEbrewEnv,
			envMap["HOMEBREW_API_DOMAIN"],
			envMap["HOMEBREW_BOTTLE_DOMAIN"],
			envMap["HOMEBREW_BREW_GIT_REMOTE"],
			envMap["HOMEBREW_CORE_GIT_REMOTE"],
			envMap["HOMEBREW_PIP_INDEX_URL"])
		that.envs.UpdateSub(utils.SUB_BREW, envars)
	default:
		l := tui.NewLog("Unknown Mirror Choice!")
		l.Info()
	}
}

func (that *Homebrew) Install() {
	if runtime.GOOS != utils.Windows {
		script := that.getShellScript()
		cmd := exec.Command("sh", script)
		cmd.Env = os.Environ()
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			tui.PrintError(err)
			return
		}
		that.SetEnv()
	} else {
		l := tui.NewLog("[Homebrew does not support Windows]")
		l.Info()
	}
}
