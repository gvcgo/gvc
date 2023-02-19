package vctrl

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

var (
	gPattern string = `# GVC Start
export PATH="$PATH:%s"
# GVC End`
	winBatPattern string = `@echo off
setx Path "%s;%s"
@echo on`
	winBatfile string = filepath.Join(config.GVCWorkDir, "gvcenv.bat")
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
		}
	}
	config.New() // init dirs and files
}

func setEnvForGVC(genvs string) {
	shellrc := utils.GetShellRcFile()
	if shellrc != utils.Win {
		if file, err := os.Open(shellrc); err == nil {
			defer file.Close()
			content, err := io.ReadAll(file)
			if err == nil {
				c := string(content)
				os.WriteFile(fmt.Sprintf("%s.backup", shellrc), content, 0644)
				if !strings.Contains(c, "# GVC Start") {
					s := fmt.Sprintf("%v\n%s", c, genvs)
					os.WriteFile(shellrc, []byte(strings.ReplaceAll(s, utils.GetHomeDir(), "$HOME")), 0644)
				}
			}
		}
	} else {
		content := fmt.Sprintf(winBatPattern, `%Path%`, config.GVCWorkDir)
		if err := os.WriteFile(winBatfile, []byte(content), os.ModePerm); err != nil {
			fmt.Println("[create batfile failed] ", err)
			return
		}
		if err := exec.Command("cmd", "/c", "start", winBatfile).Run(); err != nil {
			fmt.Println("[execute batfile failed] ", err)
			return
		}
		fmt.Println("Set env successed!")
	}
}
