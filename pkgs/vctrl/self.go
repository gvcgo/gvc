package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/TwiN/go-color"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
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
	if ok, _ := utils.PathIsExist(config.GVCWorkDir); !ok {
		os.MkdirAll(config.GVCWorkDir, os.ModePerm)
	}
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
	var r string
	fmt.Println(color.InYellow("Are you sure to delete gvc and the softwares it installed? [Y/N]"))
	fmt.Scan(&r)
	r = strings.TrimSpace(r)
	if strings.ToLower(r) == "y" || strings.ToLower(r) == "yes" {
		that.env.RemoveSubs()
		fmt.Println(color.InYellow("Restore your config files to webdav? [Y/N]"))
		fmt.Scan(&r)
		r = strings.TrimSpace(r)
		if strings.ToLower(r) == "y" || strings.ToLower(r) == "yes" || r == "" {
			dav := NewGVCWebdav()
			dav.GatherAndPushSettings()
		}
		if ok, _ := utils.PathIsExist(config.GVCWorkDir); ok {
			os.RemoveAll(config.GVCWorkDir)
		}
	} else {
		fmt.Println(color.InGreen("Uninstall gvc has been aborted."))
	}
}

func (that *Self) ShowInstallPath() {
	fmt.Println(color.InCyan("======================================================"))
	fmt.Println("[gvc] is installed in dir: ", color.InGreen(config.GVCWorkDir))
	fmt.Println(color.InCyan("======================================================"))
}

const (
	VERSION = "1.2.0"
)

func (that *Self) ShowVersion() {
	fmt.Println(color.InGreen("***========================================================***"))
	fmt.Println(color.InPurple("   GVC Version: ") + color.InYellow("v"+VERSION))
	fmt.Println(color.InPurple("   Github:      ") + color.InYellow("https://github.com/moqsien/gvc"))
	fmt.Println(color.InPurple("   Gitee:       ") + color.InYellow("https://gitee.com/moqsien/gvc_tools"))
	fmt.Println(color.InPurple("   Email:       ") + color.InYellow("moqsien@foxmail.com"))
	fmt.Println(color.InGreen("***========================================================***"))
}
