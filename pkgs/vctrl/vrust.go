package vctrl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	color "github.com/TwiN/go-color"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/query"
	"github.com/moqsien/gvc/pkgs/utils"
)

type RustInstaller struct {
	Conf    *config.GVConfig
	env     *utils.EnvsHandler
	fetcher *query.Fetcher
}

func NewRustInstaller() (ri *RustInstaller) {
	ri = &RustInstaller{
		fetcher: query.NewFetcher(),
		Conf:    config.New(),
		env:     utils.NewEnvsHandler(),
	}
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
	that.fetcher.GetAndSaveFile(fPath)
	return
}

func (that *RustInstaller) SetAccelerationEnv() {
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
		r = append(r, fmt.Sprintf("%s=%s", config.DistServerEnvName, that.Conf.Rust.DistServer))
		r = append(r, fmt.Sprintf("%s=%s", config.UpdateRootEnvName, that.Conf.Rust.UpdateRoot))
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
			fmt.Println("[Rust installer path] You can install rust by running rustup-init.exe @ ", color.InYellow(iPath))
			fmt.Printf("请切换到目录@ %s, 然后执行rustup-init.exe即可开始安装。", color.InYellow(iPath))
		}
	} else {
		cmd := exec.Command("sh", iPath)
		cmd.Env = that.getEnv()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Println(color.InRed("[Execute installer errored] "), err)
		}
	}
}
