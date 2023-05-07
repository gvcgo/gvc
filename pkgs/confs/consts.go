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
const (
	GVCVersion = "v0.1.0"
)

var (
	GVCWorkDir          = filepath.Join(utils.GetHomeDir(), ".gvc")
	GVCWebdavConfigPath = filepath.Join(GVCWorkDir, "webdav.json")
	GVCBackupDir        = filepath.Join(GVCWorkDir, "backup")
	GVConfigPath        = filepath.Join(GVCBackupDir, "gvc-config.json")
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
	if runtime.GOOS == utils.Windows {
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

// var (
// 	CodeUserSettingsFilePathForMac string = filepath.Join(utils.GetHomeDir(),
// 		"Library/Application Support/Code/User/settings.json")
// 	CodeKeybindingsFilePathForMac string = filepath.Join(utils.GetHomeDir(),
// 		"Library/Application Support/Code/User/keybindings.json")
// 	CodeUserSettingsFilePathForWin string = filepath.Join(utils.GetWinAppdataEnv(),
// 		`Code\User\settings.json`)
// 	CodeKeybindingsFilePathForWin string = filepath.Join(utils.GetWinAppdataEnv(),
// 		`Code\User\keybindings.json`)
// 	CodeUserSettingsFilePathForLinux string = filepath.Join(utils.GetHomeDir(),
// 		".config/Code/User/settings.json")
// 	CodeKeybindingsFilePathForLinux string = filepath.Join(utils.GetHomeDir(),
// 		".config/Code/User/keybindings.json")
// 	CodeUserSettingsBackupPath = GetUserSettingsBackupPath()
// 	CodeKeybindingsBackupPath  = GetCodeKeybindingsBackupPath()
// )

var (
	CodeUserSettingsBackupFileName = "vscode-user-settings.json"
	CodeKeybindingsBackupFileName  = "vscode-keybindings.json"
	CodeExtensionsBackupFileName   = "vscode-extensions.yml"
)

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
		fmt.Sprintf(`/shortcut:%s`, filepath.Join(GVCWorkDir, "g")),
	}
)

/*
go related
*/
var GoFilesDir = filepath.Join(GVCWorkDir, "go_files")

func getGoPath() string {
	if runtime.GOOS != utils.Windows {
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
	NVimFileDir            string = filepath.Join(GVCWorkDir, "nvim_files")
	NVimWinInitPath        string = filepath.Join(utils.GetHomeDir(), `\AppData\Local\nvim\init.vim`)
	NVimUnixInitPath       string = filepath.Join(utils.GetHomeDir(), ".config/nvim/init.vim")
	NVimInitBackupPath     string = filepath.Join(GVCBackupDir, "nvim-init.vim")
	NVimInitBackupFileName string = "nvim-init.vim"
)

func GetNVimInitPath() (r string) {
	if runtime.GOOS == utils.Windows {
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
var (
	ProxyFilesDir     = filepath.Join(GVCWorkDir, "proxy_files")
	ProxyListFilePath = filepath.Join(ProxyFilesDir, "proxy_list.json")
)

var ProxyXrayShellScript = `#!/bin/sh
export PATH=$PATH:~/.gvc/
nohup gvc xray -s > /dev/null 2>&1 &`

var ProxyXrayKeepRunningScript = `
#!/bin/sh
export PATH=$PATH:~/.gvc/
nohup gvc xray -k > /dev/null 2>&1 &`

// // Start-Process "C:\Program Files\Prometheus.io\prometheus.exe" -WorkingDirectory "C:\Program Files\Prometheus.io" -WindowStyle Hidden
var ProxyXrayBatScript = `Start-Process "%s" xray -s -WorkingDirectory "%s" -WindowStyle Hidden`

var ProxyXrayKeepRunningBat = `Start-Process "%s" xray -k -WorkingDirectory "%s" -WindowStyle Hidden`

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

var JavaLocalRepoPath string = filepath.Join(JavaFilesDir, "respository")

/*
Gradle related
*/
var (
	GradleRoot          = filepath.Join(JavaFilesDir, "gradle")
	GradleTarFilePath   = JavaTarFilesPath
	GradleUntarFilePath = JavaUnTarFilesPath
	GradleInitFilePath  = filepath.Join(JavaLocalRepoPath, ".gradle")
)

var GradleInitFileContent = `allprojects {
	group "org.springframework.boot"
 
	repositories {
		// 本地仓库
		mavenLocal()
		// 阿里公共仓库
		maven {
			url 'https://maven.aliyun.com/repository/public/'
		}
		// 阿里-谷歌
		maven {
			url 'https://maven.aliyun.com/repository/google'
		}
		// 阿里-gradle插件
		maven{
			url 'https://maven.aliyun.com/repository/gradle-plugin'
		}
		// 阿里-spring
		maven {
			url 'https://maven.aliyun.com/repository/spring/'
		}
		// 阿里-grails
		maven {
			url 'https://maven.aliyun.com/repository/grails-core'
		}
		// maven仓库
		mavenCentral()
	}
 
	configurations.all {
		resolutionStrategy.cacheChangingModulesFor 0, "minutes"
	}
}`

/*
Maven related
*/
var (
	MavenRoot            = filepath.Join(JavaFilesDir, "maven")
	MavenTarFilePath     = JavaTarFilesPath
	MavenUntarFilePath   = JavaUnTarFilesPath
	MavenSettingsFileDir = filepath.Join(MavenRoot, "conf")
)

var MavenSettingsPattern = `<?xml version="1.0" encoding="UTF-8"?>
<settings xmlns="http://maven.apache.org/SETTINGS/1.2.0"
          xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
          xsi:schemaLocation="http://maven.apache.org/SETTINGS/1.2.0 https://maven.apache.org/xsd/settings-1.2.0.xsd">
  <localRepository>%s</localRepository>
  <pluginGroups>
  </pluginGroups>

  <proxies>
  </proxies>

  <servers>
  </servers>

  <mirrors>
    <mirror>
		<id>nexus-tencentyun</id>
		<mirrorOf>*</mirrorOf>
		<name>Nexus tencentyun</name>
		<url>http://mirrors.cloud.tencent.com/nexus/repository/maven-public/</url>
    </mirror>
	<mirror>
		<id>nju_mirror</id>
		<mirrorOf>*</mirrorOf>
		<name>Nexus nju</name>
		<url>https://repo.nju.edu.cn/repository/maven-public/</url>
    </mirror>
	<mirror>
		<id>aliyunmaven</id>
		<mirrorOf>*</mirrorOf>
		<name>Nexus aliyun</name>
		<url>https://maven.aliyun.com/repository/public</url>
	</mirror>
	<mirror>
		<id>alimaven</id>
		<mirrorOf>central</mirrorOf>
		<name>aliyun maven</name>
		<url>https://maven.aliyun.com/repository/central</url>
	</mirror>
  </mirrors>

  <profiles>
  </profiles>
</settings>`

var MavenSettings = fmt.Sprintf(MavenSettingsPattern, JavaLocalRepoPath)

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

/*
Python related
*/
var (
	PythonFilesDir         string = filepath.Join(GVCWorkDir, "py_files")
	PythonToolsPath        string = filepath.Join(PythonFilesDir, "tools")
	PyenvInstallDir        string = filepath.Join(PythonToolsPath, "pyenv")
	PyenvRootPath          string = GetPyenvRootPath()
	PyenvRootName          string = "PYENV_ROOT"
	PyenvVersionsPath      string = filepath.Join(PyenvRootPath, "versions")
	PyenvCacheDir          string = GetPyenvCachePath()
	PythonBinaryPath       string = filepath.Join(PyenvRootPath, "shims")
	PyenvMirrorEnvName     string = "PYTHON_BUILD_MIRROR_URL"
	PyenvMirrorEnabledName string = "PYTHON_BUILD_MIRROR_URL_SKIP_CHECKSUM"
)

func GetPyenvRootPath() (r string) {
	if runtime.GOOS == utils.Windows {
		r = filepath.Join(PyenvInstallDir, "pyenv/pyenv-win")
	} else {
		r = PythonFilesDir
	}
	return
}

func GetPyenvCachePath() (r string) {
	if runtime.GOOS == utils.Windows {
		r = filepath.Join(GetPyenvRootPath(), "install_cache")
	} else {
		r = filepath.Join(GetPyenvRootPath(), "cache")
	}
	return
}

var (
	PythonUnixEnvPattern string = `# python env start
# pyenv root
export %s=%s
# pyenv & python executable path
export PATH=%s:%s:$PATH
# python env end`

	PipConfig = `[global]
timeout = 6000 
index-url = %s
trusted-host = %s`
)

func GetPipConfPath() (r string) {
	if runtime.GOOS != utils.Windows {
		r = filepath.Join(utils.GetHomeDir(), ".pip/pip.conf")
	} else {
		appdata := os.Getenv("APPDATA")
		if ok, _ := utils.PathIsExist(appdata); ok {
			r = filepath.Join(appdata, "pip/pip.ini")
		} else {
			fmt.Println("Cannot find appdata dir.")
		}
	}
	return
}

var (
	PyenvModifyForUnix   string = `    verify_checksum "$package_filename" "$checksum" >&4 2>&1 || return 1`
	PyenvAfterModifyUnix string = `    #verify_checksum "$package_filename" "$checksum" >&4 2>&1 || return 1
	echo "download completed."`
	PyenvModifyForwin1   string = `verDef = versions(version)`
	PyenvAfterModifyWin1 string = `verDef = versions(version)
	Dim list
	Dim url
	Dim mirror
	mirror = objws.Environment("Process")("PYTHON_BUILD_MIRROR_URL")
	If mirror = "" Then mirror = "https://www.python.org/ftp/python"
	list = split(verDef(LV_URL), "/ftp/python/")
	url = mirror+list(1)
	WScript.Echo url`
	PyenvModifyForwin2   string = `verDef(LV_URL), _`
	PyenvAfterModifyWin2 string = `url, _`
)

/*
C/C++ related
*/
var (
	CppFilesDir            = filepath.Join(GVCWorkDir, "cpp_files")
	Msys2Dir               = filepath.Join(CppFilesDir, "msys2")
	CygwinRootDir   string = filepath.Join(CppFilesDir, "cygwin")
	CygwinBinaryDir string = filepath.Join(CygwinRootDir, "bin")
	VCpkgDir               = filepath.Join(CppFilesDir, "vcpkg")
	CppDownloadDir         = filepath.Join(CppFilesDir, "download")
)

/*
Homebrew related
*/
var HomebrewFileDir string = filepath.Join(GVCWorkDir, "homebrew_files")

/*
Vlang related
*/
var (
	VlangFilesDir string = filepath.Join(GVCWorkDir, "vlang_files")
	VlangRootDir  string = filepath.Join(VlangFilesDir, "v")
)

/*
Flutter related
*/
var (
	FlutterFilesDir      string = filepath.Join(GVCWorkDir, "flutter_files")
	FlutterRootDir       string = filepath.Join(FlutterFilesDir, "flutter")
	FlutterTarFilePath   string = filepath.Join(FlutterFilesDir, "downloads")
	FlutterUntarFilePath string = filepath.Join(FlutterFilesDir, "versions")
)

/*
Julia related
*/
var (
	JuliaFilesDir      string = filepath.Join(GVCWorkDir, "julia_files")
	JuliaRootDir       string = filepath.Join(JuliaFilesDir, "julia")
	JuliaTarFilePath   string = filepath.Join(JuliaFilesDir, "downloads")
	JuliaUntarFilePath string = filepath.Join(JuliaFilesDir, "versions")
)

/*
Typst related
*/
var (
	TypstFilesDir string = filepath.Join(GVCWorkDir, "typst_files")
	TypstRootDir  string = filepath.Join(TypstFilesDir, "typst")
)

/*
Chatgpt related
*/
var ChatgptFilesDir string = GVCBackupDir

const (
	ChatgptConversationFileName string = "chatgpt_conversation.json"
	ChatgptConfigFileName       string = "chatgpt_config.yml"
)
