package vctrl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gvcgo/asciinema/cmd"
	"github.com/gvcgo/asciinema/util"
	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	config "github.com/gvcgo/gvc/pkgs/confs"
	"github.com/gvcgo/gvc/pkgs/utils"
)

const (
	AsciinemaConfName string = "asciinema_config"
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
	os.Setenv(util.DefaultHomeEnv, config.GVCDir)
	return &AsciiCast{
		runner: cmd.New(),
	}
}

func (that *AsciiCast) Rec(fPath string) {
	that.runner.Title, that.runner.FilePath = handleFilePath(fPath)
	that.runner.Rec()
}

func (that *AsciiCast) Play(fPath string) {
	that.runner.Title, that.runner.FilePath = handleFilePath(fPath)
	that.runner.Play()
}

func (that *AsciiCast) Auth() {
	authUrl, info := that.runner.Auth()
	gprint.PrintInfo(info)
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
		gprint.PrintError(fmt.Sprintf("Execution failed: %+v", err))
	}
}

func (that *AsciiCast) Upload(fPath string) {
	that.runner.Title, that.runner.FilePath = handleFilePath(fPath)
	if respStr, err := that.runner.Upload(); err == nil {
		gprint.PrintInfo(respStr)
	}
}

func (that *AsciiCast) HandleAsciinemaConf(toDownload bool) {
	fPath := filepath.Join(config.GVCDir, AsciinemaConfName)
	remoteFileName := AsciinemaConfName
	repoSyncer := NewSynchronizer()
	if toDownload {
		// download and deploy.
		repoSyncer.DownloadFile(
			fPath,
			remoteFileName,
			EncryptByAES,
		)
	} else {
		repoSyncer.UploadFile(
			fPath,
			remoteFileName,
			EncryptByAES,
		)
	}
}
