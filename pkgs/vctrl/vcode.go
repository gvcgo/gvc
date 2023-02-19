package vctrl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mholt/archiver/v3"
	"github.com/moqsien/gvc/pkgs/config"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
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
	Conf     *config.Conf
	*downloader.Downloader
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
		Downloader: &downloader.Downloader{
			ManuallyRedirect: true,
		},
	}
	co.initeDirs()
	return
}

func (that *Code) initeDirs() {
	if ok, _ := utils.PathIsExist(config.CodeFileDir); !ok {
		if err := os.MkdirAll(config.CodeFileDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.CodeTarFileDir); !ok {
		if err := os.MkdirAll(config.CodeTarFileDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
}

func (that *Code) getPackages() (r string) {
	that.Url = that.Conf.Config.Code.DownloadUrl
	that.Timeout = 30 * time.Second
	if resp := that.GetUrl(); resp != nil {
		rjson, _ := io.ReadAll(resp.Body)
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
		fmt.Println("[Get vscode package info failed]")
	}
	return
}

func (that *Code) download() (r string) {
	that.getPackages()
	key := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	// key := "windows-arm64"
	if p := that.Packages[key]; p != nil {
		fmt.Println(p.Url)
		var suffix string
		if strings.HasSuffix(p.Url, ".zip") {
			suffix = ".zip"
		} else if strings.HasSuffix(p.Url, ".tar.gz") {
			suffix = ".tar.gz"
		} else {
			fmt.Println("[Unsupported file type] ", p.Url)
			return
		}
		fpath := filepath.Join(config.CodeTarFileDir, fmt.Sprintf("%s-%s%s", key, that.Version, suffix))
		that.Url = strings.Replace(p.Url, that.Conf.Config.Code.StableUrl, that.Conf.Config.Code.CdnUrl, -1)
		that.Timeout = 60 * time.Second
		if size := that.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size > 0 {
			if ok := utils.CheckFile(fpath, p.CheckType, p.CheckSum); ok {
				r = fpath
			} else {
				os.RemoveAll(fpath)
			}
		}
	} else {
		fmt.Println("Cannot find package for ", key)
	}
	if ok, _ := utils.PathIsExist(config.CodeUntarFile); !ok {
		if r != "" {
			if err := archiver.Unarchive(r, config.CodeUntarFile); err != nil {
				os.RemoveAll(config.CodeUntarFile)
				fmt.Println("[Unarchive failed] ", err)
				return
			}
		}
	}
	return
}

func (that *Code) InstallForWin() {
	that.download()
	if codeDir, _ := os.ReadDir(config.CodeUntarFile); len(codeDir) > 0 {
		for _, file := range codeDir {
			if strings.Contains(file.Name(), ".exe") {
				if err := utils.WindowsSetEnv("Path", fmt.Sprintf("%s;%s", "%Path%", config.CodeWinCmdBinaryDir)); err != nil {
					fmt.Println("[Set envs failed] ", err)
					return
				}
				if err := utils.MkSymLink(filepath.Join(config.CodeUntarFile, "Code.exe"), config.CodeWinShortcutPath); err != nil {
					fmt.Println("[Create shortcut failed] ", err)
					return
				}
				break
			}
		}
	}
}

func (that *Code) addEnvForUnix(binaryDir string) {
	shellrc := utils.GetShellRcFile()
	if file, err := os.Open(shellrc); err == nil {
		defer file.Close()
		content, err := io.ReadAll(file)
		if err == nil {
			c := string(content)
			os.WriteFile(fmt.Sprintf("%s.backup", shellrc), content, 0644)
			envir := fmt.Sprintf(config.CodeEnvForUnix, binaryDir)
			if !strings.Contains(c, "# VSCode start") {
				s := fmt.Sprintf("%v\n%s", c, envir)
				os.WriteFile(shellrc, []byte(strings.ReplaceAll(s, utils.GetHomeDir(), "$HOME")), 0644)
			}
		}
	}
}

func (that *Code) InstallForMac() {
	that.download()
	if codeDir, _ := os.ReadDir(config.CodeUntarFile); len(codeDir) > 0 {
		for _, file := range codeDir {
			if strings.Contains(file.Name(), ".app") {
				source := filepath.Join(config.CodeUntarFile, file.Name())
				if ok, _ := utils.PathIsExist(config.CodeMacCmdBinaryDir); !ok {
					if err := utils.CopyFileOnUnixSudo(source, config.CodeMacInstallDir); err != nil {
						fmt.Println("[Install vscode failed] ", err)
					} else {
						os.RemoveAll(config.CodeUntarFile)
					}
				}
			}
		}
	}
	that.addEnvForUnix(config.CodeMacCmdBinaryDir)
}

func (that *Code) InstallForLinux() {
	that.download()
	if codeDir, _ := os.ReadDir(config.CodeUntarFile); len(codeDir) > 0 && len(codeDir) < 3 {
		for _, file := range codeDir {
			if file.IsDir() {
				binaryDir := filepath.Join(config.CodeUntarFile, file.Name(), "bin")
				that.addEnvForUnix(binaryDir)
			}
		}
	}
}

func (that *Code) Install() {
	switch runtime.GOOS {
	case "windows":
		that.InstallForWin()
	case "darwin":
		if ok, _ := utils.PathIsExist(config.CodeMacCmdBinaryDir); !ok {
			that.InstallForMac()
		}
	case "linux":
		that.InstallForLinux()
	}
}
