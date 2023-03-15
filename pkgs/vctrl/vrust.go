package vctrl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
)

type RustInstaller struct {
	*downloader.Downloader
	Conf *config.GVConfig
}

func NewRustInstaller() (ri *RustInstaller) {
	ri = &RustInstaller{
		Downloader: &downloader.Downloader{},
		Conf:       config.New(),
	}
	return
}

func (that *RustInstaller) getInstaller() (fPath string) {
	that.Timeout = 10 * time.Minute
	if runtime.GOOS == utils.Windows {
		that.Url = that.Conf.Rust.UrlWin
		fPath = filepath.Join(config.RustFilesDir, that.Conf.Rust.FileNameWin)

	} else {
		that.Url = that.Conf.Rust.UrlUnix
		fPath = filepath.Join(config.RustFilesDir, that.Conf.Rust.FileNameUnix)
	}
	that.GetFile(fPath, os.O_CREATE|os.O_WRONLY, 0777)
	return
}

func (that *RustInstaller) SetAccelerationEnv() {
	if runtime.GOOS == utils.Windows {
		if os.Getenv(config.DistServerEnvName) == "" {
			utils.SetWinEnv(config.DistServerEnvName, that.Conf.Rust.DistServer)
			utils.SetWinEnv(config.UpdateRootEnvName, that.Conf.Rust.UpdateRoot)
		}
	} else {
		if os.Getenv(config.DistServerEnvName) == "" {
			eRust := fmt.Sprintf(config.RustEnvPattern,
				config.DistServerEnvName,
				that.Conf.Rust.DistServer,
				config.UpdateRootEnvName,
				that.Conf.Rust.UpdateRoot)
			utils.SetUnixEnv(eRust)
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
		if err := exec.Command("cmd", "/c", iPath).Run(); err != nil {
			fmt.Println("[Execute installer errored] ", err)
		}
	} else {
		cmd := exec.Command("sh", iPath)
		cmd.Env = that.getEnv()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Println("[Execute installer errored] ", err)
		}
	}
}
