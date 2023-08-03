package vctrl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/moqsien/asciinema/cmd"
	"github.com/moqsien/asciinema/util"
	tui "github.com/moqsien/goutils/pkgs/gtui"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

func getName(base string) string {
	if base == "" {
		return base
	}
	return strings.Split(base, ".")[0]
}

func handleFilePath(fpath string) (title, result string) {
	cwd, _ := os.Getwd()
	if fpath == "" {
		return "default_cast", filepath.Join(cwd, "default.cast")
	}
	base := filepath.Base(fpath)
	if base == fpath {
		return getName(base), filepath.Join(cwd, base)
	}
	return getName(base), fpath
}

// asciinema
type AsciiCast struct {
	runner *cmd.Runner
}

func NewAsciiCast() *AsciiCast {
	os.Setenv(util.DefaultHomeEnv, config.GVCBackupDir)
	return &AsciiCast{
		runner: cmd.New(),
	}
}

func (that *AsciiCast) Rec(fPath string) {
	that.runner.FilePath, that.runner.Title = handleFilePath(fPath)
	that.runner.Rec()
}

func (that *AsciiCast) Play(fPath string) {
	that.runner.FilePath, that.runner.Title = handleFilePath(fPath)
	that.runner.Play()
}

func (that *AsciiCast) Auth() {
	authUrl, info := that.runner.Auth()
	tui.PrintInfo(info)
	var cmd *exec.Cmd
	if runtime.GOOS == utils.MacOS {
		cmd = exec.Command("open", authUrl)
	} else if runtime.GOOS == utils.Linux {
		cmd = exec.Command("x-www-browser", authUrl)
	} else if runtime.GOOS == utils.Windows {
		cmd = exec.Command("cmd", "/c", "start", authUrl)
	} else {
		return
	}
	if err := cmd.Run(); err != nil {
		tui.PrintError(fmt.Sprintf("Execution failed: %+v", err))
	}
}

func (that *AsciiCast) Upload(fPath string) {
	that.runner.FilePath, that.runner.Title = handleFilePath(fPath)
	if respStr, err := that.runner.Upload(); err == nil {
		tui.PrintInfo(respStr)
	}
}
