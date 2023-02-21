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

var (
	gPattern string = `# GVC Start
export PATH="$PATH:%s"
# GVC End`
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
			genvs := fmt.Sprintf(gPattern, config.GVCWorkDir)
			setEnvForGVC(genvs)
			if runtime.GOOS == "windows" {
				setShortcut()
			}
		}
	}
	// init dirs and files
	config.New()
}

func setEnvForGVC(genvs string) {
	shellrc := utils.GetShellRcFile()
	if shellrc != utils.Win {
		utils.SetUnixEnv(genvs)
	} else {
		utils.SetWinEnv("PATH", config.GVCWorkDir)
	}
	fmt.Println("GVC set env successed!")
}

func setShortcut() {
	if ok, _ := utils.PathIsExist(filepath.Join(config.GVCWorkDir, "gvc.exe")); ok {
		exec.Command("wscript", config.GVCShortcutCommand...)
	}
}
