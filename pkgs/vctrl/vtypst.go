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

type Typst struct {
	Conf    *config.GVConfig
	fetcher *query.Fetcher
	env     *utils.EnvsHandler
}

func NewTypstVersion() (tv *Typst) {
	tv = &Typst{
		Conf:    config.New(),
		fetcher: query.NewFetcher(),
		env:     utils.NewEnvsHandler(),
	}
	tv.env.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *Typst) download(force bool) string {
	vUrls := that.Conf.Typst.GiteeUrls

	tui.PrintInfo("Choose your URL to download: ")
	pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
		{Level: 0, Text: "From Gitee (by default & fast in China).", TextStyle: pterm.NewStyle(pterm.FgCyan), Bullet: "1)", BulletStyle: pterm.NewStyle(pterm.FgYellow)},
		{Level: 0, Text: "From Github.", TextStyle: pterm.NewStyle(pterm.FgCyan), Bullet: "2)", BulletStyle: pterm.NewStyle(pterm.FgYellow)},
	}).Render()
	fmt.Print(pterm.Cyan("Input>> "))
	var choice string
	fmt.Scan(&choice)
	choice = strings.TrimSpace(choice)
	if choice == "2" {
		vUrls = that.Conf.Typst.GithubUrls
	}
	that.fetcher.Url = vUrls[runtime.GOOS]
	suffix := utils.GetExt(that.fetcher.Url)
	if that.fetcher.Url != "" {
		fpath := filepath.Join(config.TypstFilesDir, fmt.Sprintf("typst%s", suffix))
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

func (that *Typst) renameDir() {
	dList, _ := os.ReadDir(config.TypstFilesDir)
	for _, d := range dList {
		if d.IsDir() && strings.Contains(d.Name(), "typst-") {
			os.Rename(filepath.Join(config.TypstFilesDir, d.Name()),
				filepath.Join(config.TypstRootDir))
		}
	}
}

func (that *Typst) Install(force bool) {
	zipFilePath := that.download(force)
	if ok, _ := utils.PathIsExist(config.TypstRootDir); ok && !force {
		tui.PrintInfo("Vlang is already installed.")
		return
	} else {
		os.RemoveAll(config.TypstRootDir)
	}
	if err := archiver.Unarchive(zipFilePath, config.TypstFilesDir); err != nil {
		os.RemoveAll(config.TypstRootDir)
		os.RemoveAll(zipFilePath)
		tui.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
		return
	}
	that.renameDir()
	if ok, _ := utils.PathIsExist(config.TypstRootDir); ok {
		that.CheckAndInitEnv()
	} else {
		tui.PrintError("Install typst failed.")
	}
}

func (that *Typst) CheckAndInitEnv() {
	if runtime.GOOS != utils.Windows {
		typstEnv := fmt.Sprintf(utils.TypstEnv, config.TypstRootDir)
		that.env.UpdateSub(utils.SUB_TYPST, typstEnv)
	} else {
		envList := map[string]string{
			"PATH": config.TypstRootDir,
		}
		that.env.SetEnvForWin(envList)
	}
}
