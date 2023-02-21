package vctrl

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gogf/gf/os/genv"
	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
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
	Conf     *config.GVConfig
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
	that.Url = that.Conf.Code.DownloadUrl
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
		that.Url = strings.Replace(p.Url, that.Conf.Code.StableUrl, that.Conf.Code.CdnUrl, -1)
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
				if err := utils.SetWinEnv("PATH", config.CodeWinCmdBinaryDir); err != nil {
					fmt.Println("[Set envs failed] ", err)
					return
				}
				// Automatically create shortcut.
				that.GenerateShortcut()
				break
			}
		}
	}
}

func (that *Code) GenerateShortcut() error {
	if ok, _ := utils.PathIsExist(config.WinShortcutCreatorPath); !ok {
		if err := os.WriteFile(config.WinShortcutCreatorPath, []byte(config.WinShortcutCreator), os.ModePerm); err != nil {
			fmt.Println("[Generate shortcut failed] ", err)
			return err
		}
	}
	return exec.Command("wscript", config.WinVSCodeShortcutCommand...).Run()
}

func (that *Code) addEnvForUnix(binaryDir string) {
	utils.SetUnixEnv(fmt.Sprintf(config.CodeEnvForUnix, binaryDir))
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

func (that *Code) installExtension(extName string) error {
	cmd := exec.Command("code", "--install-extension", extName)
	cmd.Env = genv.All()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (that *Code) InstallExts() {
	for _, extName := range that.Conf.Code.ExtIdentifiers {
		that.installExtension(extName)
	}
}

func (that *Code) SyncInstalledExts() {
	cmd := exec.Command("code", "--list-extensions")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %sn", err)
		return
	}
	iNameList := strings.Split(string(out), "\n")
	if len(iNameList) > 0 {
		newList := []string{}
		fmt.Println("Local installed vscode extensions: ")
		for _, iName := range iNameList {
			if strings.Contains(iName, ".") && len(iName) > 3 {
				newList = append(newList, iName)
				fmt.Println(iName)
			}
		}
		if len(newList) > 0 {
			that.Conf.Code.ExtIdentifiers = newList
		}
		that.Conf.Restore()
		that.Conf.Push()
	}
}

func (that *Code) GetSettings() {
	// get vscode settings from remote webdav.
	if ok, _ := utils.PathIsExist(config.CodeUserSettingsBackupPath); ok {
		if ok, _ := utils.PathIsExist(filepath.Dir(config.GetCodeUserSettingsPath())); ok {
			utils.CopyFile(config.CodeUserSettingsBackupPath, config.GetCodeUserSettingsPath())
		}
	}

	if ok, _ := utils.PathIsExist(config.CodeKeybindingsBackupPath); ok {
		if ok, _ := utils.PathIsExist(filepath.Dir(config.GetCodeKeybindingsPath())); ok {
			utils.CopyFile(config.CodeKeybindingsBackupPath, config.GetCodeKeybindingsPath())
		}
	}
}

func (that *Code) SyncSettings() {
	// push vscode settings to remote webdav.
	if ok, _ := utils.PathIsExist(config.GetCodeUserSettingsPath()); ok {
		if ok, _ := utils.PathIsExist(filepath.Dir(config.CodeUserSettingsBackupPath)); ok {
			utils.CopyFile(config.GetCodeUserSettingsPath(), config.CodeUserSettingsBackupPath)
			that.Conf.Push()
		}
	}

	if ok, _ := utils.PathIsExist(config.GetCodeKeybindingsPath()); ok {
		if ok, _ := utils.PathIsExist(filepath.Dir(config.CodeKeybindingsBackupPath)); ok {
			utils.CopyFile(config.GetCodeKeybindingsPath(), config.CodeKeybindingsBackupPath)
			that.Conf.Push()
		}
	}
}
