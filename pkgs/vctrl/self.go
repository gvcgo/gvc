package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

type Self struct {
	Conf    *config.GVConfig
	env     *utils.EnvsHandler
	checker *SumChecker
}

func NewSelf() (s *Self) {
	s = &Self{
		Conf: config.New(),
		env:  utils.NewEnvsHandler(),
	}
	s.checker = NewSumChecker(s.Conf)
	s.env.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *Self) setEnv() {
	if runtime.GOOS != utils.Windows {
		that.env.UpdateSub(utils.SUB_GVC, fmt.Sprintf(utils.GvcEnv, config.GVCWorkDir))
	} else {
		that.env.SetEnvForWin(map[string]string{"PATH": config.GVCWorkDir})
	}
}

func (that *Self) setShortcut() {
	switch runtime.GOOS {
	case utils.Windows:
		fPath := filepath.Join(config.GVCWorkDir, "gvc.exe")
		if ok, _ := utils.PathIsExist(fPath); ok {
			newPath := filepath.Join(config.GVCWorkDir, "g.exe")
			os.RemoveAll(newPath)
			utils.CopyFile(fPath, newPath)
		}
	default:
		fPath := filepath.Join(config.GVCWorkDir, "gvc")
		if ok, _ := utils.PathIsExist(fPath); ok {
			utils.MkSymLink(fPath, filepath.Join(config.GVCWorkDir, "g"))
		}
	}
}

func (that *Self) Install() {
	utils.MakeDirs(config.GVCWorkDir)
	ePath, _ := os.Executable()
	if strings.Contains(ePath, filepath.Join(utils.GetHomeDir(), ".gvc")) && !strings.Contains(ePath, "bin_temp") {
		// call the installed exe is not allowed.
		return
	}
	name := filepath.Base(ePath)
	if strings.HasSuffix(ePath, "/gvc") || strings.HasSuffix(ePath, "gvc.exe") {
		if _, err := utils.CopyFile(ePath, filepath.Join(config.GVCWorkDir, name)); err == nil {
			that.setEnv()
			that.setShortcut()
		}
		// reset config file to default.
		that.Conf.SetDefault()
		that.Conf.Restore()
	}
	nPath := filepath.Join(filepath.Dir(ePath), ".neobox_encrypt_key.json")
	if ok, _ := utils.PathIsExist(nPath); ok {
		os.RemoveAll(nPath)
	}
}

func (that *Self) Uninstall() {
	confirmPrinter := pterm.DefaultInteractiveConfirm
	confirmPrinter.DefaultText = "Confirm to remove gvc. "
	confirmPrinter.TextStyle = &pterm.Style{pterm.FgRed}
	if result, _ := confirmPrinter.Show(); result {
		pterm.Println()
		that.env.RemoveSubs()
		cp := pterm.DefaultInteractiveConfirm
		cp.DefaultText = "Confirm to save config files to WebDAV. "
		if r, _ := cp.Show(); r {
			dav := NewGVCWebdav()
			dav.GatherAndPushSettings()
		}
		pterm.Println()
		if ok, _ := utils.PathIsExist(config.GVCWorkDir); ok {
			os.RemoveAll(config.GVCWorkDir)
		}
	} else {
		tui.PrintInfo("Remove has been aborted.")
	}
}

func (that *Self) ShowPath() {
	str := pterm.DefaultBox.
		WithRightPadding(1).
		WithLeftPadding(1).
		WithTopPadding(1).
		WithBottomPadding(1).
		Sprintf("%s: %s\n%s: %s\n%s: %s",
			pterm.Cyan("Installation Dir"), pterm.Green(config.GVCWorkDir),
			pterm.Cyan("GVC Config Path"), pterm.Green(config.GVConfigPath),
			pterm.Cyan("GVC Webdav Config Path"), pterm.Green(config.GVCWebdavConfigPath))
	pterm.Println(str)
}

const (
	VERSION = "1.3.x"
)

func (that *Self) ShowVersion() {
	name, _ := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("G", pterm.FgCyan.ToStyle()),
		putils.LettersFromStringWithStyle("VC", pterm.FgLightMagenta.ToStyle()),
	).Srender()

	pterm.Println(name)
	str := pterm.DefaultBox.
		WithRightPadding(2).
		WithLeftPadding(2).
		WithTopPadding(2).
		WithBottomPadding(2).
		Sprintf(
			"%s\n%s\n%s",
			pterm.LightCyan("   Version:     ")+pterm.LightYellow("v"+VERSION),
			pterm.LightCyan("   Github:      ")+pterm.LightYellow("https://github.com/moqsien/gvc"),
			pterm.LightCyan("   Email:       ")+pterm.LightYellow("moqsien@foxmail.com"),
		)
	pterm.Println(str)
}

func (that *Self) CheckLatestVersion(currentVersion string) {
	latest := that.checker.GetLatestGVCVersion()
	if currentVersion == latest {
		tui.PrintInfo(fmt.Sprintf("Current version: %s is already the latest.", currentVersion))
		return
	}
	confirmPrinter := pterm.DefaultInteractiveConfirm
	confirmPrinter.DefaultText = "To download the latest version for GVC or not. "
	confirmPrinter.TextStyle = &pterm.Style{pterm.FgRed}
	if result, _ := confirmPrinter.Show(); result {
		that.download()
	}
}

func (that *Self) download() {
	dUrl := that.Conf.GVC.GitlabUrls[fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)]
	if dUrl != "" {
		fPath := filepath.Join(config.GVCBinTempDir, "gvc.zip")
		fetcher := request.NewFetcher()
		fetcher.SetUrl(dUrl)
		fetcher.Timeout = 20 * time.Minute
		temBinPath := filepath.Join(config.GVCBinTempDir, "gvc")
		if runtime.GOOS == utils.Windows {
			temBinPath += ".exe"
		}
		// remove old files before get a new one
		os.RemoveAll(fPath)
		os.RemoveAll(temBinPath)
		if err := fetcher.DownloadAndDecompress(fPath, config.GVCBinTempDir, true); err != nil {
			tui.PrintError(err)
		} else {
			tui.PrintSuccess(config.GVCBinTempDir)
		}
	}
}
