package vctrl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Self struct {
	Conf *config.GVConfig
	env  *utils.EnvsHandler
}

func NewSelf() (s *Self) {
	return &Self{
		Conf: config.New(),
		env:  utils.NewEnvsHandler(),
	}
}

func (that *Self) setEnv() {
	if runtime.GOOS != utils.Windows {
		that.env.UpdateSub(utils.SUB_GVC, fmt.Sprintf(utils.GvcEnv, config.GVCWorkDir))
	} else {
		// utils.SetWinEnv("PATH", config.GVCWorkDir)
		gEnv := map[string]string{
			"PATH": config.GVCWorkDir,
		}
		that.env.SetEnvForWin(gEnv)
	}
}

func (that *Self) setShortcut() {
	if ok, _ := utils.PathIsExist(filepath.Join(config.GVCWorkDir, "gvc.exe")); ok {
		exec.Command("wscript", config.GVCShortcutCommand...)
	}
}

func (that *Self) Install() {
	if runtime.GOOS == utils.Windows {
		that.env.HintsForWin(1)
	}
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
			if runtime.GOOS == utils.Windows {
				that.setShortcut()
			}
		}
	}
	// init dirs and files
	config.New()
}

func (that *Self) Uninstall() {
	var r string
	fmt.Println("Are you sure to delete gvc and the softwares it installed?[Y/N]")
	fmt.Scan(&r)
	r = strings.TrimSpace(r)
	if strings.ToLower(r) == "y" || strings.ToLower(r) == "yes" {
		that.env.RemoveSubs()
		fmt.Println("Restore your config files to webdav?[Y/N]")
		fmt.Scan(&r)
		r = strings.TrimSpace(r)
		if strings.ToLower(r) == "y" || strings.ToLower(r) == "yes" || r == "" {
			that.Conf.Push()
		}
		if ok, _ := utils.PathIsExist(config.GVCWorkDir); ok {
			os.RemoveAll(config.GVCWorkDir)
		}
	} else {
		fmt.Println("Uninstall gvc aborted.")
	}
}

func (that *Self) ShowInstallPath() {
	fmt.Println("===================================")
	fmt.Println("[gvc] is installed @", config.GVCWorkDir)
	fmt.Println("[gvc] 安装目录: ", config.GVCWorkDir)
	fmt.Println("===================================")
}
