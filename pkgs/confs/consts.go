package confs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/moqsien/gvc/pkgs/utils"
)

/*
gvc related
*/
const GVCVersion = "v0.1.0"

var (
	GVCWorkDir          = filepath.Join(utils.GetHomeDir(), ".gvc")
	GVCWebdavConfigPath = filepath.Join(GVCWorkDir, "webdav.yml")
	GVCBackupDir        = filepath.Join(GVCWorkDir, "backup")
	GVConfigPath        = filepath.Join(GVCBackupDir, "gvc-config.yml")
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
		return
	}
	r = HostFilePathForNix
	return
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
	CodeWinShortcutPath string = filepath.Join(utils.GetHomeDir(), `Desktop\vscode`)
)

var (
	CodeEnvForUnix string = `# VSCode start
export PATH="%s:$PATH"
# VSCode end`
)

var (
	CodeUserSettingsFilePathForMac string = filepath.Join(utils.GetHomeDir(),
		"Library/Application Support/Code/User/settings.json")
	CodeKeybindingsFilePathForMac string = filepath.Join(utils.GetHomeDir(),
		"Library/Application Support/Code/User/keybindings.json")
	CodeUserSettingsFilePathForWin string = filepath.Join(utils.GetWinAppdataEnv(),
		`Code\User\settings.json`)
	CodeKeybindingsFilePathForWin string = filepath.Join(utils.GetWinAppdataEnv(),
		`Code\User\keybindings.json`)
	CodeUserSettingsFilePathForLinux string = filepath.Join(utils.GetHomeDir(),
		".config/Code/User/settings.json")
	CodeKeybindingsFilePathForLinux string = filepath.Join(utils.GetHomeDir(),
		".config/Code/User/keybindings.json")
	CodeUserSettingsBackupPath = filepath.Join(GVCBackupDir, "vscode-settings.json")
	CodeKeybindingsBackupPath  = filepath.Join(GVCBackupDir, "vscode-keybindings.json")
)

func GetCodeUserSettingsPath() string {
	switch runtime.GOOS {
	case "darwin":
		return CodeUserSettingsFilePathForMac
	case "linux":
		return CodeUserSettingsFilePathForLinux
	case "windows":
		return CodeUserSettingsFilePathForWin
	default:
		return ""
	}
}

func GetCodeKeybindingsPath() string {
	switch runtime.GOOS {
	case "darwin":
		return CodeKeybindingsFilePathForMac
	case "linux":
		return CodeKeybindingsFilePathForLinux
	case "windows":
		return CodeKeybindingsFilePathForWin
	default:
		return ""
	}
}

// shortcut maker for windows.
var WinShortcutCreator = `set WshShell = WScript.CreateObject("WScript.Shell" )
set oShellLink = WshShell.CreateShortcut(Wscript.Arguments.Named("shortcut") & ".lnk")
oShellLink.TargetPath = Wscript.Arguments.Named("target")
oShellLink.WindowStyle = 1
oShellLink.Save`

var (
	WinShortcutCreatorName          = "sc.vbs"
	WinShortcutCreatorPath   string = filepath.Join(GVCWorkDir, WinShortcutCreatorName)
	WinVSCodeShortcutCommand        = []string{
		WinShortcutCreatorPath,
		fmt.Sprintf(`/target:%s`, filepath.Join(CodeUntarFile, "Code.exe")),
		fmt.Sprintf(`/shortcut:%s`, CodeWinShortcutPath),
	}
)

func SaveWinShortcutCreator() {
	if ok, _ := utils.PathIsExist(WinShortcutCreatorPath); !ok {
		os.WriteFile(WinShortcutCreatorPath, []byte(WinShortcutCreator), os.ModePerm)
	}
}

var (
	GVCShortcutCommand = []string{
		WinShortcutCreatorPath,
		fmt.Sprintf(`/target:%s`, filepath.Join(GVCWorkDir, "gvc.exe")),
		fmt.Sprintf(`/shortcut:%s`, filepath.Join(GVCWorkDir, "gvc")),
	}
)

/*
go related
*/
var GoFilesDir = filepath.Join(GVCWorkDir, "go_files")

func getGoPath() string {
	if runtime.GOOS != "windows" {
		return "data/projects/go"
	}
	return `data\projects\go`
}

var (
	DefaultGoRoot    string = filepath.Join(GoFilesDir, "go")
	DefaultGoPath    string = filepath.Join(utils.GetHomeDir(), getGoPath())
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

/*
Neovim related.
*/
var (
	NVimFileDir        string = filepath.Join(GVCWorkDir, "nvim_files")
	NVimWinInitPath    string = filepath.Join(utils.GetHomeDir(), `\AppData\Local\nvim\init.vim`)
	NVimUnixInitPath   string = filepath.Join(utils.GetHomeDir(), ".config/nvim/init.vim")
	NVimInitBackupPath string = filepath.Join(GVCBackupDir, "nvim-init.vim")
)

func GetNVimInitPath() (r string) {
	if runtime.GOOS == "windows" {
		r = NVimWinInitPath
	} else {
		r = NVimUnixInitPath
	}
	dir := filepath.Dir(r)
	if ok, _ := utils.PathIsExist(dir); !ok {
		os.MkdirAll(dir, os.ModePerm)
	}
	return r
}

func GetNVimPlugDir() string {
	return filepath.Dir(GetNVimInitPath())
}

var (
	NVimUnixEnv string = `# nvim start
export PATH="%s:$PATH"
# nvim end`
)

/*
Proxy related
*/
var ProxyFilesDir = filepath.Join(GVCWorkDir, "proxy_files")

/*
Java related
*/
var JavaFilesDir = filepath.Join(GVCWorkDir, "java_files")

var (
	DefaultJavaRoot    = filepath.Join(JavaFilesDir, "java")
	JavaTarFilesPath   = filepath.Join(JavaFilesDir, "downloads")
	JavaUnTarFilesPath = filepath.Join(JavaFilesDir, "versions")
)

var JavaEnvVarPattern string = `# Java Env started
export JAVA_HOME="%s"
export CLASS_PATH="$JAVA_HOME/lib"
export PATH="$JAVA_HOME/bin:$JAVA_HOME/lib/tools.jar:$JAVA_HOME/lib/dt.jar:$PATH"
# Java Env Ended`

/*
Rust related
*/
var (
	RustFilesDir      = filepath.Join(GVCWorkDir, "rust_files")
	DistServerEnvName = "RUSTUP_DIST_SERVER"
	UpdateRootEnvName = "RUSTUP_UPDATE_ROOT"
)

var RustEnvPattern string = `# Rust env start
export %s=%s
export %s=%s
# Rust env end`

/*
Nodejs related
*/
var (
	NodejsFilesDir   = filepath.Join(GVCWorkDir, "nodejs_files")
	NodejsRoot       = filepath.Join(NodejsFilesDir, "nodejs")
	NodejsTarFiles   = filepath.Join(NodejsFilesDir, "downloads")
	NodejsUntarFiles = filepath.Join(NodejsFilesDir, "versions")
	NodejsGlobal     = filepath.Join(NodejsFilesDir, "node_global")
	NodejsCache      = filepath.Join(NodejsFilesDir, "node_cache")
)

var NodejsEnvPattern string = `# Nodejs env start
export NODE_HOME="%s"
export PATH="$NODE_HOME/bin:$PATH"
# Nodejs env end`
