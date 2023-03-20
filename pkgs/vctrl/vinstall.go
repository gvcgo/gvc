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

func SelfInstall() {
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
			setEnvForGVC()
			if runtime.GOOS == utils.Windows {
				setShortcut()
			}
		}
	}
	// init dirs and files
	config.New()
}

var gvcEnv string = `export  PATH="$PATH:%s"`

func setEnvForGVC() {
	if runtime.GOOS != utils.Windows {
		eh := utils.NewEnvsHandler()
		eh.UpdateSub(utils.SUB_GVC, fmt.Sprintf(gvcEnv, config.GVCWorkDir))
	} else {
		utils.SetWinEnv("PATH", config.GVCWorkDir)
	}
}

func setShortcut() {
	if ok, _ := utils.PathIsExist(filepath.Join(config.GVCWorkDir, "gvc.exe")); ok {
		exec.Command("wscript", config.GVCShortcutCommand...)
	}
}
