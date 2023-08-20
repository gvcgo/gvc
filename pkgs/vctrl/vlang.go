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
	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Vlang struct {
	Conf    *config.GVConfig
	env     *utils.EnvsHandler
	fetcher *request.Fetcher
	checker *SumChecker
}

func NewVlang() (vl *Vlang) {
	vl = &Vlang{
		Conf:    config.New(),
		fetcher: request.NewFetcher(),
		env:     utils.NewEnvsHandler(),
	}
	vl.checker = NewSumChecker(vl.Conf)
	vl.env.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *Vlang) download(force bool) string {
	optionList := []string{
		"From Gitlab[Default]",
		"From Github",
	}
	sel := tui.NewSelect(optionList, func(s string) string {
		var vUrls map[string]string
		switch s {
		case optionList[1]:
			vUrls = that.Conf.Vlang.VlangUrls
		default:
			vUrls = that.Conf.Vlang.VlangGitlabUrls
		}
		return vUrls[runtime.GOOS]
	})
	sel.SetDefaultText("Choose your download URL")
	that.fetcher.Url = sel.Start()
	if that.fetcher.Url != "" {
		fpath := filepath.Join(config.VlangFilesDir, "vlang.zip")
		if strings.Contains(that.fetcher.Url, "gitlab.com") && !that.checker.IsUpdated(fpath, that.fetcher.Url) {
			tui.PrintInfo("Current version is already the latest.")
			return fpath
		}
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
		tui.PrintInfo("Vlang is already installed.")
		return
	} else {
		os.RemoveAll(config.VlangRootDir)
	}
	if err := archiver.Unarchive(zipFilePath, config.VlangFilesDir); err != nil {
		os.RemoveAll(config.VlangRootDir)
		os.RemoveAll(zipFilePath)
		tui.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
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
	optionList := []string{
		"From Gitlab[Default]",
		"From Github",
	}
	sel := tui.NewSelect(optionList, func(s string) string {
		var vUrls map[string]string
		switch s {
		case optionList[1]:
			vUrls = that.Conf.Vlang.AnalyzerUrls
		default:
			vUrls = that.Conf.Vlang.AnalyzerGitlabUrls
		}
		key := runtime.GOOS
		if runtime.GOOS == utils.MacOS {
			key = fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
		}
		return vUrls[key]
	})
	sel.SetDefaultText("Choose your download URL")
	that.fetcher.Url = sel.Start()

	if that.fetcher.Url != "" {
		fpath := filepath.Join(config.VlangFilesDir, "analyzer.zip")
		if strings.Contains(that.fetcher.Url, "gitlab.com") && !that.checker.IsUpdated(fpath, that.fetcher.Url) {
			tui.PrintInfo("Current version is already the latest.")
			return
		}
		that.fetcher.Timeout = 20 * time.Minute
		that.fetcher.SetThreadNum(3)
		if ok, _ := utils.PathIsExist(fpath); !ok {
			if err := that.fetcher.DownloadAndDecompress(fpath, config.VlangFilesDir, true); err == nil {
				tui.PrintSuccess(fpath)
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

// TODO: protoc installation
// TODO: sudo for windows, using grpc
