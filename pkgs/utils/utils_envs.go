package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/koanfer"
)

/*
Environment consts
*/
const (
	GVC_BLOCK_START = "# ==GVC== block start\n"
	GVC_BLOCK_END   = "# ==GVC== block end"
	SUB_BLOCK_START = "# sub block @%s start\n"
	SUB_BLOCK_END   = "# sub block @%s end"
)

/*
Sub block name
*/
const (
	SUB_GVC     = "gvc"
	SUB_GO      = "go"
	SUB_JDK     = "jdk"
	SUB_GRADLE  = "gradle"
	SUB_MAVEN   = "maven"
	SUB_PY      = "python"
	SUB_NODE    = "nodejs"
	SUB_RUST    = "rust"
	SUB_CODE    = "vscode"
	SUB_NVIM    = "neovim"
	SUB_BREW    = "homebrew"
	SUB_VLANG   = "vlang"
	SUB_FLUTTER = "flutter"
	SUB_ANDROID = "android"
	SUB_LUA     = "lua"
	SUB_JULIA   = "julia"
	SUB_TYPST   = "typst"
	SUB_VCPKG   = "vcpkg"
	SUB_PROTOC  = "protoc"
)

/*
Block Patterns
*/
const (
	BLOCK_PATTERN = "%s%s\n%s"
)

/*
gvc Envs
*/
var GvcEnv string = `export  PATH="$PATH:%s"`

/*
Go Envs
*/
var GoEnv string = `export GOROOT="%s"
export GOPATH="%s"
export GOBIN="%s"
export GOPROXY="%s"
export PATH="%s"`

/*
Protobuf Envs
*/
var ProtoEnv string = `export PATH="%s:$PATH"`

/*
Java Envs
*/
var JavaEnv string = `export JAVA_HOME="%s"
export CLASS_PATH="$JAVA_HOME/lib"
export PATH="$JAVA_HOME/bin:$JAVA_HOME/lib/tools.jar:$JAVA_HOME/lib/dt.jar:$PATH"`

/*
Gradle Envs
*/
var GradleEnv string = `export GRADLE_HOME="%s"
export PATH=$GRADLE_HOME/bin:$PATH
export GRADLE_USER_HOME=%s`

/*
Maven Envs
*/
var MavenEnv string = `export MAVEN_HOME="%s"
export PATH=$MAVEN_HOME/bin:$PATH`

/*
Nodejs Envs
*/
var NodeEnv string = `export NODE_HOME="%s"
export PATH="$NODE_HOME/bin:$PATH"`

/*
Python & Pyenv Envs
*/
var PyEnv string = `# pyenv root
export %s=%s
# pyenv & python executable path
export PATH=%s:%s:$PATH`

/*
Rust Envs for acceleration
*/
var RustEnv string = `export %s=%s
export %s=%s`

/*
Neovim Envs
*/
var NVimEnv string = `export PATH="%s:$PATH"`

/*
VSCode Envs
*/
var VSCodeEnv string = `export PATH="%s:$PATH"`

/*
Homebrew Envs
*/
var HOMEbrewEnv string = `export HOMEBREW_API_DOMAIN="%s"
export HOMEBREW_BOTTLE_DOMAIN="%s"
export HOMEBREW_BREW_GIT_REMOTE="%s"
export HOMEBREW_CORE_GIT_REMOTE="%s"
export HOMEBREW_PIP_INDEX_URL="%s"`

/*
Vlang Envs
*/
var VlangEnv string = `export PATH="%s:$PATH"`

/*
Flutter Envs
*/
var FlutterEnv string = `export FLUTTER_ROOT="%s"
export PUB_HOSTED_URL=%s
export FLUTTER_STORAGE_BASE_URL=%s
export FLUTTER_GIT_URL=%s
export PATH="$FLUTTER_ROOT/bin:$PATH"`

/*
Android cmdline tools Envs
*/
var AndroidEnv string = `export ANDROID_HOME="%s"
export PATH="%s:$PATH"`

/*
Julia Envs
*/
var JuliaEnv string = `export JULIA_ROOT="%s"
export JULIA_PKG_SERVER="%s"
export PATH="$JULIA_ROOT/bin:$PATH"`

/*
Typst Envs
*/
var TypstEnv string = `export PATH="%s:$PATH"`

/*
VCPKG Envs
*/
var VcpkgEnv string = `export PATH="%s:$PATH"`

type WinPathEnvTemp struct {
	PathList []string `koanf,json:"path_list"`
}

type EnvsHandler struct {
	shellName   string
	rcFilePath  string
	oldContent  []byte
	winPathTemp *WinPathEnvTemp
	koanfer     *koanfer.JsonKoanfer
	winWorkdir  string
}

func NewEnvsHandler() (e *EnvsHandler) {
	e = &EnvsHandler{}
	if runtime.GOOS == Windows {
		e.shellName = PowerShell
		e.getRcFilePath()
		e.getOldContents()
	} else {
		e.shellName = getShellTypeForUnix()
		e.getRcFilePath()
		e.getOldContents()
	}
	e.winPathTemp = &WinPathEnvTemp{PathList: []string{}}
	return
}

func getShellTypeForUnix() (st string) {
	s := os.Getenv("SHELL")
	if strings.Contains(s, Zsh) {
		st = Zsh
	} else if strings.Contains(s, Bash) {
		st = Bash
	} else {
		gprint.PrintError("Please use zsh or bash.")
		os.Exit(1)
	}
	return
}

func (that *EnvsHandler) getRcFilePath() {
	switch that.shellName {
	case Zsh:
		that.rcFilePath = filepath.Join(GetHomeDir(), ".zshrc")
	case Bash:
		that.rcFilePath = filepath.Join(GetHomeDir(), ".bashrc")
	case PowerShell:
		that.rcFilePath = ""
	default:
		that.rcFilePath = ""
	}
}

func (that *EnvsHandler) flushEnvs() {
	if runtime.GOOS != Windows {
		cmd := exec.Command("source", that.rcFilePath)
		cmd.Env = os.Environ()
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func (that *EnvsHandler) getOldContents() {
	if runtime.GOOS == Windows {
		return
	}
	if ok, _ := PathIsExist(that.rcFilePath); ok {
		that.oldContent, _ = os.ReadFile(that.rcFilePath)
	} else {
		os.Create(that.rcFilePath)
	}
}

func (that *EnvsHandler) getBlockValue(start, end string) (r []byte) {
	reg := regexp.MustCompile(fmt.Sprintf(`%s([\s\S]*)%s`, start, end))
	cList := reg.FindAllSubmatch(that.oldContent, -1)
	if len(cList) > 0 && len(cList[0]) > 1 {
		r = cList[0][1]
	}
	return
}

func (that *EnvsHandler) updateRcFile(newContent string) error {
	return os.WriteFile(that.rcFilePath, []byte(newContent), os.ModePerm)
}

func (that *EnvsHandler) replaceBlockContent(start, end, value string) {
	reg := regexp.MustCompile(fmt.Sprintf(`%s([\s\S]*)%s`, start, end))
	newSub := ""
	if value != "" {
		value = strings.ReplaceAll(value, "$", "￥")
		newSub = fmt.Sprintf(BLOCK_PATTERN, start, value, end)
	}
	newContent := reg.ReplaceAllString(string(that.oldContent), newSub)
	newContent = strings.ReplaceAll(newContent, "￥", "$")
	that.updateRcFile(newContent)
}

func (that *EnvsHandler) UpdateSub(subname, value string) {
	value = strings.ReplaceAll(value, GetHomeDir(), `$HOME`)
	sub_start := fmt.Sprintf(SUB_BLOCK_START, subname)
	sub_end := fmt.Sprintf(SUB_BLOCK_END, subname)
	oCon := string(that.oldContent)
	if strings.Contains(oCon, sub_start) {
		that.replaceBlockContent(sub_start, sub_end, value)
		that.getOldContents()
	} else if strings.Contains(oCon, GVC_BLOCK_START) {
		start := GVC_BLOCK_START
		end := GVC_BLOCK_END
		gvc_content := that.getBlockValue(start, end)
		new_value := fmt.Sprintf("%s%s", gvc_content,
			fmt.Sprintf(BLOCK_PATTERN, sub_start, value, sub_end))
		that.replaceBlockContent(start, end, new_value)
		that.getOldContents()
	} else {
		new_value := fmt.Sprintf(BLOCK_PATTERN,
			GVC_BLOCK_START,
			fmt.Sprintf(BLOCK_PATTERN, sub_start, value, sub_end),
			GVC_BLOCK_END)
		newContent := fmt.Sprintf("%s\n%s", string(that.oldContent), new_value)
		that.updateRcFile(newContent)
		that.getOldContents()
	}
	that.flushEnvs()
}

func (that *EnvsHandler) RemoveSub(subname string) {
	sub_start := fmt.Sprintf(SUB_BLOCK_START, subname)
	sub_end := fmt.Sprintf(SUB_BLOCK_END, subname)
	if strings.Contains(string(that.oldContent), sub_start) {
		that.replaceBlockContent(sub_start, sub_end, "")
	}
	that.flushEnvs()
}

func (that *EnvsHandler) RemoveSubs() {
	if runtime.GOOS != Windows {
		if strings.Contains(string(that.oldContent), GVC_BLOCK_START) {
			that.replaceBlockContent(GVC_BLOCK_START, GVC_BLOCK_END, "")
		}
		that.flushEnvs()
	}
}

func (that *EnvsHandler) DoesEnvExist(subname string) bool {
	if runtime.GOOS != Windows {
		return strings.Contains(string(that.oldContent), fmt.Sprintf(SUB_BLOCK_START, subname))
	} else {
		return false
	}
}

/*
Windows PowerShell Environment Settings
*/
const (
	TempName string = "windows_env_temp.json"
)

func (that *EnvsHandler) SetWinWorkDir(dirPath string) {
	if ok, _ := PathIsExist(dirPath); ok {
		that.winWorkdir = dirPath
	} else {
		gprint.PrintError(fmt.Sprintf("[%s] does not exist.", dirPath))
	}
}

func (that *EnvsHandler) loadTemp() {
	if ok, _ := PathIsExist(that.winWorkdir); ok {
		fPath := filepath.Join(that.winWorkdir, TempName)
		if ok, _ := PathIsExist(fPath); !ok {
			that.saveTemp()
		}
		if that.koanfer == nil {
			that.koanfer, _ = koanfer.NewKoanfer(fPath)
		}
		if that.koanfer != nil {
			that.koanfer.Load(that.winPathTemp)
		}
	}
}

func (that *EnvsHandler) saveTemp() {
	if ok, _ := PathIsExist(that.winWorkdir); ok {
		fPath := filepath.Join(that.winWorkdir, TempName)
		if that.koanfer == nil {
			that.koanfer, _ = koanfer.NewKoanfer(fPath)
		}
		if that.koanfer != nil {
			that.koanfer.Save(*that.winPathTemp)
		}
	}
}

func (that *EnvsHandler) preparePathToSet(value string) string {
	that.loadTemp()
	pathEnv := os.Getenv("PATH")
	for _, v := range that.winPathTemp.PathList {
		if !strings.Contains(pathEnv, v) {
			value = fmt.Sprintf("%s;%s", value, v)
		}
	}
	return value
}

func (that *EnvsHandler) addToTemp(value string) {
	that.loadTemp()
	for _, v := range that.winPathTemp.PathList {
		if v == value {
			return
		}
	}
	that.winPathTemp.PathList = append(that.winPathTemp.PathList, value)
	that.saveTemp()
}

var (
	PwSetPathEnv  string = `[Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";%s", "User")`
	PwSetOtherEnv string = `[Environment]::SetEnvironmentVariable("%s", "%s", "User")`
)

func (that *EnvsHandler) setOneEnvForWin(key, value string) {
	var arg string
	_key := strings.ToLower(key)
	if _key == "path" && strings.Contains(os.Getenv("PATH"), value) {
		gprint.PrintInfo(fmt.Sprintf("[%s] Already exists in Path.", value))
		return
	}
	originValue := value
	if _key == "path" {
		value = that.preparePathToSet(value)
		arg = fmt.Sprintf(PwSetPathEnv, value)
	} else {
		arg = fmt.Sprintf(PwSetOtherEnv, key, value)
	}
	cmd := exec.Command(PowerShell, arg)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		gprint.PrintError(fmt.Sprintf("Set env [%s] failed: %+v.", key, err))
		return
	}

	if _key == "path" {
		that.addToTemp(originValue)
	}
}

func (that *EnvsHandler) SetEnvForWin(envList map[string]string) {
	for key, value := range envList {
		that.setOneEnvForWin(key, value)
	}
}
