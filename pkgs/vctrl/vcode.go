package vctrl

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mholt/archiver/v3"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/pterm/pterm"
	"github.com/tidwall/gjson"
)

type CodePackage struct {
	OsArchName string
	Url        string
	CheckSum   string
	CheckType  string
}

type Code struct {
	Version  string
	Packages map[string]*CodePackage
	Conf     *config.GVConfig
	env      *utils.EnvsHandler
	fetcher  *request.Fetcher
}

type typeMap map[string]string

var codeType typeMap = typeMap{
	"win32-x64-archive":   "windows-amd64",
	"win32-arm64-archive": "windows-arm64",
	"linux-x64":           "linux-amd64",
	"linux-arm64":         "linux-arm64",
	"darwin":              "darwin-amd64",
	"darwin-arm64":        "darwin-arm64",
}

func NewCode() (co *Code) {
	co = &Code{
		Packages: make(map[string]*CodePackage),
		Conf:     config.New(),
		fetcher:  request.NewFetcher(),
		env:      utils.NewEnvsHandler(),
	}
	co.fetcher.NoRedirect = true
	co.initeDirs()
	co.env.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *Code) initeDirs() {
	utils.MakeDirs(config.CodeFileDir, config.CodeTarFileDir)
}

func (that *Code) getPackages() (r string) {
	that.fetcher.Url = that.Conf.Code.DownloadUrl
	that.fetcher.Timeout = 60 * time.Second
	if resp := that.fetcher.Get(); resp != nil {
		defer resp.RawBody().Close()
		rjson, _ := io.ReadAll(resp.RawBody())
		products := gjson.Get(string(rjson), "products")
		for _, product := range products.Array() {
			if that.Version == "" {
				that.Version = product.Get("productVersion").String()
			}
			osArch := product.Get("platform.os")
			if localOsArch, ok := codeType[osArch.String()]; ok {
				that.Packages[localOsArch] = &CodePackage{
					OsArchName: osArch.String(),
					Url:        product.Get("url").String(),
					CheckSum:   product.Get("sha256hash").String(),
					CheckType:  "sha256",
				}
			}
		}
	} else {
		gprint.PrintError("Get vscode package info failed.")
	}
	return
}

func (that *Code) download() (r string) {
	that.getPackages()
	key := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	if p := that.Packages[key]; p != nil {
		fmt.Println(p.Url)
		var suffix string
		if strings.HasSuffix(p.Url, ".zip") {
			suffix = ".zip"
		} else if strings.HasSuffix(p.Url, ".tar.gz") {
			suffix = ".tar.gz"
		} else {
			gprint.PrintError(fmt.Sprintf("Unsupported file type: %s", p.Url))
			return
		}
		fpath := filepath.Join(config.CodeTarFileDir, fmt.Sprintf("%s-%s%s", key, that.Version, suffix))
		pterm.Println(pterm.Green("Use CDN to accelerate download or not?"))
		result, _ := pterm.DefaultInteractiveConfirm.Show()
		pterm.Println()
		if result {
			that.fetcher.Url = strings.Replace(p.Url, that.Conf.Code.StableUrl, that.Conf.Code.CdnUrl, -1)
		}
		that.fetcher.Timeout = 600 * time.Second
		that.fetcher.SetThreadNum(8)
		if size := that.fetcher.GetAndSaveFile(fpath); size > 0 {
			if ok := utils.CheckFile(fpath, p.CheckType, p.CheckSum); ok {
				r = fpath
			} else {
				os.RemoveAll(fpath)
			}
		}
	} else {
		gprint.PrintError(fmt.Sprintf("Cannot find package for %s", key))
	}

	if ok, _ := utils.PathIsExist(config.CodeUntarFile); !ok {
		that.Unarchive(r)
	} else {
		if runtime.GOOS == utils.Windows || runtime.GOOS == utils.Linux {
			os.RemoveAll(config.CodeUntarFile)
			that.Unarchive(r)
		}
	}
	return
}

func (that *Code) Unarchive(fPath string) {
	if fPath != "" {
		if err := archiver.Unarchive(fPath, config.CodeUntarFile); err != nil {
			os.RemoveAll(config.CodeUntarFile)
			gprint.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
			return
		}
	}
}

func (that *Code) InstallForWin() {
	that.download()
	if codeDir, _ := os.ReadDir(config.CodeUntarFile); len(codeDir) > 0 {
		for _, file := range codeDir {
			if strings.Contains(file.Name(), ".exe") {
				if !strings.Contains(os.Getenv("PATH"), config.CodeWinCmdBinaryDir) {
					that.env.SetEnvForWin(map[string]string{
						"PATH": config.CodeWinCmdBinaryDir,
					})
				}
				// Automatically create shortcut.
				that.GenerateShortcut()
				break
			}
		}
	}
}

func (that *Code) GenerateShortcut() error {
	config.SaveWinShortcutCreator()
	if ok, _ := utils.PathIsExist(config.WinShortcutCreatorPath); ok {
		args := append([]string{"wscript"}, config.WinVSCodeShortcutCommand...)
		_, err := utils.ExecuteSysCommand(false, args...)
		return err
	}
	return errors.New("shortcut script not found")
}

func (that *Code) InstallForMac() {
	that.download()
	if codeDir, _ := os.ReadDir(config.CodeUntarFile); len(codeDir) > 0 {
		for _, file := range codeDir {
			if strings.Contains(file.Name(), ".app") {
				source := filepath.Join(config.CodeUntarFile, file.Name())
				if ok, _ := utils.PathIsExist(config.CodeMacCmdBinaryDir); !ok {
					if err := utils.CopyFileOnUnixSudo(source, config.CodeMacInstallDir); err != nil {
						gprint.PrintError(fmt.Sprintf("Install vscode failed: %+v", err))
					} else {
						os.RemoveAll(config.CodeUntarFile)
					}
				}
			}
		}
	}
	that.env.UpdateSub(utils.SUB_CODE, config.CodeMacCmdBinaryDir)
}

func (that *Code) InstallForLinux() {
	that.download()
	if codeDir, _ := os.ReadDir(config.CodeUntarFile); len(codeDir) > 0 && len(codeDir) < 3 {
		for _, file := range codeDir {
			if file.IsDir() {
				binaryDir := filepath.Join(config.CodeUntarFile, file.Name(), "bin")
				that.env.UpdateSub(utils.SUB_CODE, binaryDir)
			}
		}
	}
}

func (that *Code) Install() {
	switch runtime.GOOS {
	case utils.Windows:
		that.InstallForWin()
	case utils.MacOS:
		if ok, _ := utils.PathIsExist(config.CodeMacCmdBinaryDir); !ok {
			that.InstallForMac()
		}
	case utils.Linux:
		that.InstallForLinux()
	}
}
