package vctrl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	config "github.com/moqsien/gvc/pkgs/confs"
	downloader "github.com/moqsien/gvc/pkgs/fetcher"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Homebrew struct {
	Conf *config.GVConfig
	*downloader.Downloader
	envs *utils.EnvsHandler
}

func NewHomebrew() (hb *Homebrew) {
	hb = &Homebrew{
		Conf:       config.New(),
		Downloader: &downloader.Downloader{},
		envs:       utils.NewEnvsHandler(),
	}
	return
}

func (that *Homebrew) getShellScript() string {
	fPath := filepath.Join(config.HomebrewFileDir, "homebrew.sh")
	if ok, _ := utils.PathIsExist(fPath); !ok {
		that.Url = that.Conf.Homebrew.ShellScriptUrl
		if size := that.GetFile(fPath, os.O_CREATE|os.O_WRONLY, 0777); size > 0 {
			return fPath
		}
		return ""
	}
	return fPath
}

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
		fmt.Println("Unknown Mirror Choice!")
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
			fmt.Println("[Install Homebrew failed] ", err)
			return
		}
		that.SetEnv()
	} else {
		fmt.Println("[Homebrew does not support Windows]")
	}
}
