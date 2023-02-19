package confs

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/moqsien/gvc/pkgs/utils"
)

/*
gvc related
*/
const GVCVersion = "v0.1.0"

var (
	GVCWorkDir          = filepath.Join(utils.GetHomeDir(), ".gvc/")
	GVCWebdavConfigPath = filepath.Join(GVCWorkDir, "webdav.yml")
	GVCBackupDir        = filepath.Join(GVCWorkDir, "backup")
	DefaultConfigPath   = filepath.Join(GVCWorkDir, "config.yml")
	RealConfigPath      = filepath.Join(GVCBackupDir, "gvc-config.yml")
)

/*
hosts related
*/
const (
	HostFilePathForNix = "/etc/hosts"
	HostFilePathForWin = `C:\Windows\System32\drivers\etc\hosts`
)

var TempHostsFilePath = filepath.Join(GVCWorkDir, "/temp_hosts.txt")

func GetHostsFilePath() (r string) {
	if runtime.GOOS == "windows" {
		r = HostFilePathForWin
	}
	r = HostFilePathForNix
	return r
}

/*
vscode related
*/
var (
	CodeFileDir         string = filepath.Join(GVCWorkDir, "vscode_file")
	CodeTarFileDir      string = filepath.Join(CodeFileDir, "downloads")
	CodeUntarFile       string = filepath.Join(CodeFileDir, "vscode")
	CodeMacInstallDir   string = "/Applications/"
	CodeMacCmdBinaryDir string = filepath.Join(CodeMacInstallDir, "Visual Studio Code.app/Contents/Resources/app/bin")
	CodeWinCmdBinaryDir string = filepath.Join(CodeUntarFile, "bin")
	CodeWinShortcutPath string = filepath.Join(utils.GetHomeDir(), "Desktop/", "Visual Studio Code")
)

var (
	CodeEnvForUnix string = `# VSCode start
export PATH="%s:$PATH"
# VSCode end`
)

/*
go related
*/
var GoFilesDir = filepath.Join(GVCWorkDir, "go_files")

var (
	DefaultGoRoot    string = filepath.Join(GoFilesDir, "go")
	DefaultGoPath    string = filepath.Join(utils.GetHomeDir(), "data/projects/go")
	DefaultGoProxy   string = "https://goproxy.cn,direct"
	GoTarFilesPath   string = filepath.Join(GoFilesDir, "downloads")
	GoUnTarFilesPath string = filepath.Join(GoFilesDir, "versions")
)

var (
	GoUnixEnvsPattern string = `# Golang Start
export GOROOT="%s"
export GOPATH="%s"
export GOBIN="%s"
export GOPROXY="%s"
export PATH="%s"
# Golang End`
	GoUnixEnv string = fmt.Sprintf(GoUnixEnvsPattern,
		DefaultGoRoot,
		DefaultGoPath,
		filepath.Join(DefaultGoPath, "bin"),
		`%s`,
		`%s`)
)

var (
	GoWinBatPattern string = `@echo off
setx "GOROOT" "%s"
setx "GOPATH" "%s"
setx "GORIN" "%s"
setx "GOPROXY" "%s"
setx Path "%s"
@echo on
`
	GoWinBatPath string = filepath.Join(GoFilesDir, "genv.bat")
	GoWinEnv     string = fmt.Sprintf(GoWinBatPattern,
		DefaultGoRoot,
		DefaultGoPath,
		filepath.Join(DefaultGoPath, "bin"),
		`%s`,
		`%s`)
)
