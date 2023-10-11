package confs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/moqsien/goutils/pkgs/gtea/input"
	"github.com/moqsien/gvc/pkgs/utils"
)

/*
gvc related
*/
const (
	GVCVersion = "v0.1.0"
)

var (
	GVCDir              = filepath.Join(utils.GetHomeDir(), ".gvc")
	GVCInstallDir       = GetGVCWorkDir()
	GVCWebdavConfigPath = filepath.Join(GVCDir, "webdav.json")
	GVCBackupDir        = filepath.Join(GVCDir, "backup")
	GVConfigPath        = filepath.Join(GVCBackupDir, "gvc-config.json")
	GVCBinTempDir       = filepath.Join(GVCDir, "bin_temp")
)

func GetGVCWorkDir() string {
	installPathConfig := filepath.Join(GVCDir, "pkg_install_path.conf")
	if ok, _ := utils.PathIsExist(installPathConfig); ok {
		content, _ := os.ReadFile(installPathConfig)
		return string(content)
	}
	ipt := input.NewInput(input.WithPlaceholder(`set where to install packages; default: "$HomeDir/.gvc/"`), input.WithWidth(100))
	ipt.Run()
	d := ipt.Value()
	if ok, _ := utils.PathIsExist(d); ok {
		os.WriteFile(installPathConfig, []byte(d), os.ModePerm)
		return d
	} else if err := os.MkdirAll(d, os.ModePerm); err == nil {
		os.WriteFile(installPathConfig, []byte(d), os.ModePerm)
		return d
	} else {
		os.WriteFile(installPathConfig, []byte(GVCDir), os.ModePerm)
		return GVCDir
	}
}

/*
windows gsudo
*/
var GsudoFilePath = filepath.Join(GVCInstallDir, "gsudo_files")

/*
hosts related
*/
const (
	HostFilePathForNix = "/etc/hosts"
	HostFilePathForWin = `C:\Windows\System32\drivers\etc\hosts`
)

var TempHostsFilePath = filepath.Join(GVCDir, "temp_hosts.txt")

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
	CodeFileDir           string = filepath.Join(GVCInstallDir, "vscode_file")
	CodeTarFileDir        string = filepath.Join(CodeFileDir, "downloads")
	CodeUntarFile         string = filepath.Join(CodeFileDir, "vscode")
	CodeMacInstallDir     string = "/Applications/"
	CodeMacCmdBinaryDir   string = filepath.Join(CodeMacInstallDir, "Visual Studio Code.app", "Contents", "Resources", "app", "bin")
	CodeWinInstallDir     string = filepath.Join(utils.GetHomeDir(), "AppData", "Local", "Programs", "Microsoft VS Code")
	CodeWinCmdBinaryDir   string = filepath.Join(CodeWinInstallDir, "bin")
	CodeLinuxInstallDir   string = "/usr/share/code"
	CodeLinuxCmdBinaryDir string = filepath.Join(CodeLinuxInstallDir, "bin")
	CodeWinShortcutPath   string = filepath.Join(utils.GetHomeDir(), "Desktop", "VSCode")
)

var (
	CodeEnvForUnix string = `# VSCode start
export PATH="%s:$PATH"
# VSCode end`
)

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
	WinShortcutCreatorName        = "sc.vbs"
	WinShortcutCreatorPath string = filepath.Join(GVCDir, WinShortcutCreatorName)
)

func SaveWinShortcutCreator() {
	if ok, _ := utils.PathIsExist(WinShortcutCreatorPath); !ok {
		os.WriteFile(WinShortcutCreatorPath, []byte(WinShortcutCreator), os.ModePerm)
	}
}

var (
	GVCShortcutCommand = []string{
		WinShortcutCreatorPath,
		fmt.Sprintf(`/target:%s`, filepath.Join(GVCDir, "gvc.exe")),
		fmt.Sprintf(`/shortcut:%s`, filepath.Join(GVCDir, "g")),
	}
)

func CreateShortCut(targetPath, shortcutPath string) error {
	if runtime.GOOS != utils.Windows {
		return os.Symlink(targetPath, shortcutPath)
	} else {
		WinVSCodeShortcutCommand := []string{
			WinShortcutCreatorPath,
			fmt.Sprintf(`/target:%s`, targetPath),
			fmt.Sprintf(`/shortcut:%s`, shortcutPath),
		}
		args := append([]string{"wscript"}, WinVSCodeShortcutCommand...)
		_, err := utils.ExecuteSysCommand(false, args...)
		return err
	}
}

/*
go related
*/
var GoFilesDir = filepath.Join(GVCInstallDir, "go_files")

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
Protobuf related
*/
var (
	ProtobufDir string = filepath.Join(GVCInstallDir, "protobuf_files")
)

/*
Neovim related.
*/
var (
	NVimFileDir            string = filepath.Join(GVCInstallDir, "nvim_files")
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
	ProxyFilesDir     = filepath.Join(GVCInstallDir, "proxy_files")
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
var JavaFilesDir = filepath.Join(GVCInstallDir, "java_files")

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
	RustFilesDir      = filepath.Join(GVCInstallDir, "rust_files")
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
	NodejsFilesDir   = filepath.Join(GVCInstallDir, "nodejs_files")
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
	PythonFilesDir         string = filepath.Join(GVCInstallDir, "py_files")
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
		r = filepath.Join(PyenvInstallDir, "pyenv", "pyenv-win")
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

	PyenvWinOriginalCacheDir    string = `strDirCache  = strPyenvHome & "\install_cache"`
	PyenvWinNewCacheDir         string = `strDirCache  = "%s"`
	PyenvWinOriginalVersionsDir string = `strDirVers   = strPyenvHome & "\versions"`
	PyenvWinNewVersionsDir      string = `strDirVers   = "%s"`
	PyenvWinOriginalShimsDir    string = `strDirShims  = strPyenvHome & "\shims"`
	PyenvWinNewShimsDir         string = `strDirShims  = "%s"`

	PyenvWinBeforeFixed string = `deepExtract = objws.Run("msiexec /quiet /a """& file &""" TargetDir="""& installPath & """", 0, True)
        If deepExtract Then
            WScript.Echo ":: [Error] :: error installing """& baseName &""" component MSI."
            Exit Function
        End If`
	PyenvWinAfterFixed string = `deepExtract = objws.Run("msiexec /quiet /a """& file &""" TargetDir="""& installPath & """", 0, True)
        If deepExtract Then
            deepExtract = objws.Run("msiexec """& file &""" TargetDir="""& installPath & """", 0, True)
        End If
	If deepExtract Then
            WScript.Echo ":: [Error] :: error installing """& baseName &""" component MSI."
            Exit Function
        End If`

	PyenvWinBatBeforeFixed string = `for /f "%skip_arg%delims=" %%i in ('%pyenv% vname') do call :extrapath "%~dp0..\versions\%%i"`
	PyenvWinBatAfterFixed  string = `for /f "%skip_arg%delims=" %%i in ('%pyenv% vname') do call :extrapath "%~dp0..\versions\%%i"

for %%f in ("%bindir%") do set "pyversion=%%~nxf"
set "bindir=$$$\%pyversion%"
set "extrapaths=$$$\%pyversion%;$$$\%pyversion%\Scripts;"`

	PyenvWinOriginalPyUrl string = "https://www.python.org/ftp"
	PyenvWinTaobaoPyUrl   string = "https://registry.npmmirror.com/-/binary"
)

/*
C/C++ related
*/
var (
	CppFilesDir            = filepath.Join(GVCInstallDir, "cpp_files")
	Msys2Dir               = filepath.Join(CppFilesDir, "msys2")
	CygwinRootDir   string = filepath.Join(CppFilesDir, "cygwin")
	CygwinBinaryDir string = filepath.Join(CygwinRootDir, "bin")
	VCpkgDir               = filepath.Join(CppFilesDir, "vcpkg")
	CppDownloadDir         = filepath.Join(CppFilesDir, "download")
)

// -G 'Ninja'
var VCPkgScript string = `(cd %s && CXX="%s" eval cmake %s "-DCMAKE_BUILD_TYPE=Release -DVCPKG_DEVELOPMENT_WARNINGS=OFF") || exit 1
(cd %s && cmake --build .) || exit 1`

var Msys2CygwinGitFixBat = `@echo off
setlocal

if "%1" equ "rev-parse" goto rev_parse
git %*
goto :eof
:rev_parse
for /f %%1 in ('git %*') do cygpath -w %%1`

/*
Homebrew related
*/
var HomebrewFileDir string = filepath.Join(GVCInstallDir, "homebrew_files")

/*
Vlang related
*/
var (
	VlangFilesDir string = filepath.Join(GVCInstallDir, "vlang_files")
	VlangRootDir  string = filepath.Join(VlangFilesDir, "v")
)

/*
Flutter related
*/
var (
	FlutterFilesDir             string = filepath.Join(GVCInstallDir, "flutter_files")
	FlutterRootDir              string = filepath.Join(FlutterFilesDir, "flutter")
	FlutterTarFilePath          string = filepath.Join(FlutterFilesDir, "downloads")
	FlutterUntarFilePath        string = filepath.Join(FlutterFilesDir, "versions")
	FlutterAndroidToolDownloads string = filepath.Join(FlutterFilesDir, "android_tools")
	FlutterAndroidHomeDir       string = filepath.Join(FlutterFilesDir, "android_home")
)

/*
Julia related
*/
var (
	JuliaFilesDir      string = filepath.Join(GVCInstallDir, "julia_files")
	JuliaRootDir       string = filepath.Join(JuliaFilesDir, "julia")
	JuliaTarFilePath   string = filepath.Join(JuliaFilesDir, "downloads")
	JuliaUntarFilePath string = filepath.Join(JuliaFilesDir, "versions")
)

/*
Typst related
*/
var (
	TypstFilesDir string = filepath.Join(GVCInstallDir, "typst_files")
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

/*
Docker related
*/
var (
	DockerFilesDir               string = filepath.Join(GVCInstallDir, "docker_files")
	DockerWindowsInstallationDir string = filepath.Join(DockerFilesDir, "docker_installation")
)

/*
git
*/
var (
	GitFileDir                string = filepath.Join(GVCInstallDir, "git_files")
	GitWindowsInstallationDir string = filepath.Join(GitFileDir, "git_installation")
)
