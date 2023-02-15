package vctrl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/moqsien/gvc/pkgs/config"
	"github.com/moqsien/gvc/pkgs/utils"
)

var gPattern string = `# GVC Start
export PATH="$PATH:%s"
# GVC End`

func SelfInstall() {
	if ok, _ := utils.PathIsExist(config.GVCWorkDir); !ok {
		os.MkdirAll(config.GVCWorkDir, os.ModePerm)
	}
	ePath, _ := os.Executable()
	name := filepath.Base(ePath)
	if strings.HasSuffix(ePath, "/gvc") || strings.HasSuffix(ePath, "/gvc.exe") {
		if _, err := utils.CopyFile(ePath, filepath.Join(config.GVCWorkDir, name)); err == nil {
			genvs := fmt.Sprintf(gPattern, config.GVCWorkDir)
			setEnvForGVC(genvs)
		}
	}
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
		fmt.Println(utils.Win)
	}
}
