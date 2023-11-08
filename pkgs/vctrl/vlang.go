package vctrl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gogf/gf/os/genv"
	"github.com/mholt/archiver/v3"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Vlang struct {
	Conf    *config.GVConfig
	env     *utils.EnvsHandler
	fetcher *request.Fetcher
}

func NewVlang() (vl *Vlang) {
	vl = &Vlang{
		Conf:    config.New(),
		fetcher: request.NewFetcher(),
		env:     utils.NewEnvsHandler(),
	}
	vl.env.SetWinWorkDir(config.GVCDir)
	return
}

func (that *Vlang) download(force bool) string {
	that.fetcher.Url = that.Conf.Vlang.VlangUrls[runtime.GOOS]
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(that.fetcher.Url)

	if that.fetcher.Url != "" {
		fpath := filepath.Join(config.VlangFilesDir, "vlang.zip")
		if force {
			os.RemoveAll(fpath)
		}
		that.fetcher.Timeout = 20 * time.Minute
		that.fetcher.SetThreadNum(3)
		if ok, _ := utils.PathIsExist(fpath); !ok || force {
			if size := that.fetcher.GetAndSaveFile(fpath); size > 0 {
				return fpath
			} else {
				os.RemoveAll(fpath)
			}
		} else if ok && !force {
			return fpath
		}
	}
	return ""
}

func (that *Vlang) Install(force bool) {
	zipFilePath := that.download(force)
	if ok, _ := utils.PathIsExist(config.VlangRootDir); ok && !force {
		gprint.PrintInfo("Vlang is already installed.")
		return
	} else {
		os.RemoveAll(config.VlangRootDir)
	}
	if err := archiver.Unarchive(zipFilePath, config.VlangFilesDir); err != nil {
		os.RemoveAll(config.VlangRootDir)
		os.RemoveAll(zipFilePath)
		gprint.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
		return
	}
	if ok, _ := utils.PathIsExist(config.VlangRootDir); ok {
		that.CheckAndInitEnv()
	}
}

func (that *Vlang) CheckAndInitEnv() {
	if runtime.GOOS != utils.Windows {
		vlangEnv := fmt.Sprintf(utils.VlangEnv, config.VlangRootDir)
		that.env.UpdateSub(utils.SUB_VLANG, vlangEnv)
	} else {
		envList := map[string]string{
			"PATH": config.VlangRootDir,
		}
		that.env.SetEnvForWin(envList)
	}
}

func (that *Vlang) InstallVAnalyzerForVscode() {
	key := runtime.GOOS
	if key == utils.MacOS {
		key = fmt.Sprintf("%s_%s", key, runtime.GOARCH)
	}
	that.fetcher.Url = that.Conf.Vlang.AnalyzerUrls[key]
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(that.fetcher.Url)
	if that.fetcher.Url != "" {
		fpath := filepath.Join(config.VlangFilesDir, "analyzer.zip")
		that.fetcher.Timeout = 20 * time.Minute
		that.fetcher.SetThreadNum(3)
		if ok, _ := utils.PathIsExist(fpath); !ok {
			if err := that.fetcher.DownloadAndDecompress(fpath, config.VlangFilesDir, true); err == nil {
				gprint.PrintSuccess(fpath)
			} else {
				fmt.Println(err)
				os.RemoveAll(fpath)
				return
			}
		}
		binName := "v-analyzer"
		if runtime.GOOS == utils.Windows {
			binName = "v-analyzer.exe"
		}
		binPath := filepath.Join(config.VlangFilesDir, binName)
		if ok, _ := utils.PathIsExist(binPath); ok {
			cnf := NewGVCWebdav()
			filesToSync := cnf.GetFilesToSync()
			vscodeSettingsPath := filesToSync[config.CodeUserSettingsBackupFileName]
			if runtime.GOOS == utils.Windows {
				binPath = strings.ReplaceAll(binPath, `\`, `\\`)
			}
			utils.AddNewlineToVscodeSettings("v-analyzer.serverPath", binPath, vscodeSettingsPath)
		} else {
			return
		}
		// install extension for vscode
		cmd := exec.Command("code", "--install-extension", "vosca.vscode-v-analyzer")
		cmd.Env = genv.All()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Run()
	}
}
