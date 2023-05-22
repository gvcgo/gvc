package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/tui"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

type Self struct {
	Conf *config.GVConfig
	env  *utils.EnvsHandler
}

func NewSelf() (s *Self) {
	s = &Self{
		Conf: config.New(),
		env:  utils.NewEnvsHandler(),
	}
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
	if strings.Contains(ePath, filepath.Join(utils.GetHomeDir(), ".gvc")) {
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
	VERSION = "1.2.x"
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
			"%s\n%s\n%s\n%s",
			pterm.LightCyan("   Version:     ")+pterm.LightYellow("v"+VERSION),
			pterm.LightCyan("   Github:      ")+pterm.LightYellow("https://github.com/moqsien/gvc"),
			pterm.LightCyan("   Gitee:       ")+pterm.LightYellow("https://gitee.com/moqsien/gvc_tools"),
			pterm.LightCyan("   Email:       ")+pterm.LightYellow("moqsien@foxmail.com"),
		)
	pterm.Println(str)
}
