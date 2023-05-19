package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/query"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/tui"
	"github.com/pterm/pterm"
)

type Vlang struct {
	Conf    *config.GVConfig
	env     *utils.EnvsHandler
	fetcher *query.Fetcher
}

func NewVlang() (vl *Vlang) {
	vl = &Vlang{
		Conf:    config.New(),
		fetcher: query.NewFetcher(),
		env:     utils.NewEnvsHandler(),
	}
	vl.env.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *Vlang) download(force bool) string {
	vUrls := that.Conf.Vlang.VlangGiteeUrls

	tui.PrintInfo("Choose your URL to download: ")
	pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
		{Level: 0, Text: "From Gitee (by default & fast in China).", TextStyle: pterm.NewStyle(pterm.FgCyan), Bullet: "1)", BulletStyle: pterm.NewStyle(pterm.FgYellow)},
		{Level: 0, Text: "From Github.", TextStyle: pterm.NewStyle(pterm.FgCyan), Bullet: "2)", BulletStyle: pterm.NewStyle(pterm.FgYellow)},
	}).Render()
	fmt.Print(pterm.Cyan("Input>>"))
	var choice string
	fmt.Scan(&choice)
	choice = strings.TrimSpace(choice)
	if choice == "2" {
		vUrls = that.Conf.Vlang.VlangUrls
	}
	that.fetcher.Url = vUrls[runtime.GOOS]
	if that.fetcher.Url != "" {
		fpath := filepath.Join(config.VlangFilesDir, "vlang.zip")
		if force {
			os.RemoveAll(fpath)
		}
		if ok, _ := utils.PathIsExist(fpath); !ok || force {
			if size := that.fetcher.GetAndSaveFile(fpath); size > 0 {
				return fpath
			} else {
				os.RemoveAll(fpath)
			}
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
