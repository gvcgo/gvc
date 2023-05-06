package vctrl

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
)

var (
	Msys2InstallerName string = "msys2_installer.exe"
	// "msys2_installer.exe in --confirm-command --accept-messages --root C:/msys64"
	Msys2InstallCmd string = fmt.Sprintf("%s in --confirm-command --accept-messages --root", Msys2InstallerName)
)

type CppManager struct {
	*downloader.Downloader
	c    *colly.Collector
	Conf *config.GVConfig
	Doc  *goquery.Document
	env  *utils.EnvsHandler
}

func NewCppManager() (cm *CppManager) {
	cm = &CppManager{
		Downloader: &downloader.Downloader{},
		c:          colly.NewCollector(),
		Conf:       config.New(),
		env:        utils.NewEnvsHandler(),
	}
	cm.initDirs()
	return
}

func (that *CppManager) initDirs() {
	if ok, _ := utils.PathIsExist(config.CppFilesDir); !ok {
		if err := os.MkdirAll(config.CppFilesDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.Msys2Dir); !ok {
		if err := os.MkdirAll(config.Msys2Dir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.VCpkgDir); !ok {
		if err := os.MkdirAll(config.VCpkgDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.CppDownloadDir); !ok {
		if err := os.MkdirAll(config.CppDownloadDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
}

func (that *CppManager) getDoc() {
	mUrl := that.Conf.Cpp.MsysInstallerUrl
	if !utils.VerifyUrls(mUrl) {
		return
	}
	that.c.OnResponse(func(r *colly.Response) {
		that.Doc, _ = goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
	})
	that.c.Visit(mUrl)
}

func (that *CppManager) getInstaller() (fPath string) {
	if that.Doc == nil {
		that.getDoc()
	}
	if that.Doc != nil {
		var exeUrl string
		maxIdx := that.Doc.Find("table#list").Find("tr").Last().Index()
		for i := maxIdx; i >= 0; i-- {
			exeUrl = that.Doc.Find("table#list").Find("tr").Eq(i).Find("a").AttrOr("href", "")
			if strings.HasSuffix(exeUrl, ".exe") {
				break
			}
		}

		if exeUrl != "" {
			if !strings.HasPrefix(exeUrl, "http://") {
				exeUrl, _ = url.JoinPath(that.Conf.Cpp.MsysInstallerUrl, exeUrl)
			}
			fPath = filepath.Join(config.CppDownloadDir, Msys2InstallerName)
			that.Url = exeUrl
			that.GetFile(fPath, os.O_CREATE|os.O_WRONLY, 0777)
		}
	}
	return
}

func (that *CppManager) InstallMsys2() {
	if runtime.GOOS != utils.Windows {
		return
	}
	fPath := that.getInstaller()
	if ok, _ := utils.PathIsExist(fPath); ok {
		os.Setenv("PATH", fmt.Sprintf("%s;%s", config.CppDownloadDir, os.Getenv("PATH")))
		c := exec.Command(fmt.Sprintf("%s %s", Msys2InstallCmd, config.CppDownloadDir))
		c.Env = os.Environ()
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		if err := c.Run(); err != nil {
			fmt.Println("Execute Msys2Installer Failed: ", err)
			return
		}
		binPath := filepath.Join(config.Msys2Dir, "usr", "bin")
		if ok, _ := utils.PathIsExist(binPath); ok {
			winEnv := map[string]string{
				"PATH": fmt.Sprintf("%s:%s", binPath, config.CppDownloadDir),
			}
			if !strings.Contains(os.Getenv("PATH"), binPath) {
				that.env.SetEnvForWin(winEnv)
			}
		}
	}
}
