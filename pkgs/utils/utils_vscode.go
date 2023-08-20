package utils

import (
	"fmt"
	"os"
	"strings"

	tui "github.com/moqsien/goutils/pkgs/gtui"
)

func AddNewlineToVscodeSettings(key, value, settingsPath string) {
	// vscodeSettingsPath := filepath.Join(GetWinAppdataEnv(), `Code\User\settings.json`)
	if ok, _ := PathIsExist(settingsPath); ok {
		if vsContent, err := os.ReadFile(settingsPath); err == nil {
			strContent := strings.TrimSuffix(string(vsContent), "\n")
			if strings.Contains(strContent, fmt.Sprintf(`"%s"`, key)) {
				tui.PrintInfo(fmt.Sprintf(`"%s" already exists in: %s`, key, settingsPath))
				return
			}
			cList := strings.Split(strContent, "\n")
			length := len(cList)
			if strings.Contains(cList[length-1], "}") {
				if length-2 >= 0 {
					line := strings.TrimSuffix(cList[length-2], "\n")
					if !strings.HasSuffix(line, ",") {
						line = line + ","
					}
					line += "\n"
				}

				cList = append(cList[:length-2], fmt.Sprintf(`    "%s": "%s"`, key, value), "}")
			}
			strContent = strings.Join(cList, "\n")
			os.WriteFile(settingsPath, []byte(strContent), os.ModePerm)
		}
	}
}
