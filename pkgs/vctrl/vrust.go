package vctrl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/moqsien/goutils/pkgs/gtea/confirm"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type RustInstaller struct {
	Conf    *config.GVConfig
	env     *utils.EnvsHandler
	fetcher *request.Fetcher
}

func NewRustInstaller() (ri *RustInstaller) {
	ri = &RustInstaller{
		fetcher: request.NewFetcher(),
		Conf:    config.New(),
		env:     utils.NewEnvsHandler(),
	}
	ri.env.SetWinWorkDir(config.GVCDir)
	return
}

func (that *RustInstaller) getInstaller() (fPath string) {
	that.fetcher.Timeout = 10 * time.Minute
	if runtime.GOOS == utils.Windows {
		that.fetcher.Url = that.Conf.Rust.UrlWin
		fPath = filepath.Join(config.RustFilesDir, that.Conf.Rust.FileNameWin)

	} else {
		that.fetcher.Url = that.Conf.Rust.UrlUnix
		fPath = filepath.Join(config.RustFilesDir, that.Conf.Rust.FileNameUnix)
	}
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(that.fetcher.Url)
	that.fetcher.GetAndSaveFile(fPath)
	return
}

func (that *RustInstaller) SetAccelerationEnv() {
	cfm := confirm.NewConfirm(confirm.WithTitle("Set RUSTUP_DIST_SERVER/RUSTUP_UPDATE_ROOT to 'mirrors.ustc.edu.cn' or not?"))
	cfm.Run()
	result := cfm.Result()
	if !result {
		return
	}
	if runtime.GOOS == utils.Windows {
		if os.Getenv(config.DistServerEnvName) == "" {
			envList := map[string]string{
				config.DistServerEnvName: that.Conf.Rust.DistServer,
				config.UpdateRootEnvName: that.Conf.Rust.UpdateRoot,
			}
			that.env.SetEnvForWin(envList)
		}
	} else {
		if os.Getenv(config.DistServerEnvName) == "" {
			that.env.UpdateSub(utils.SUB_RUST, fmt.Sprintf(utils.RustEnv,
				config.DistServerEnvName,
				that.Conf.Rust.DistServer,
				config.UpdateRootEnvName,
				that.Conf.Rust.UpdateRoot))
		}
	}
}

func (that *RustInstaller) getEnv() (r []string) {
	r = os.Environ()
	if !strings.Contains(strings.Join(r, " "), config.DistServerEnvName) {
		cfm := confirm.NewConfirm(confirm.WithTitle("Set RUSTUP_DIST_SERVER/RUSTUP_UPDATE_ROOT to 'mirrors.ustc.edu.cn' or not?"))
		cfm.Run()
		result := cfm.Result()
		if result {
			r = append(r, fmt.Sprintf("%s=%s", config.DistServerEnvName, that.Conf.Rust.DistServer))
			r = append(r, fmt.Sprintf("%s=%s", config.UpdateRootEnvName, that.Conf.Rust.UpdateRoot))
		}
	}
	return
}

func (that *RustInstaller) Install() {
	that.SetAccelerationEnv()
	iPath := that.getInstaller()
	if runtime.GOOS == utils.Windows {
		os.Setenv("PATH", fmt.Sprintf("%s;%s", config.RustFilesDir, os.Getenv("PATH")))
		c := exec.Command(that.Conf.Rust.FileNameWin)
		c.Env = os.Environ()
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		if err := c.Run(); err != nil {
			gprint.PrintInfo(fmt.Sprintf("You can install rust by running rustup-init.exe in %s.", iPath))
		}
	} else {
		cmd := exec.Command("sh", iPath)
		cmd.Env = that.getEnv()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			gprint.PrintError(fmt.Sprintf("Execute installer failed: %+v", err))
		}
	}
}
