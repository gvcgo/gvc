package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
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
	SUB_PY      = "python"
	SUB_NODE    = "nodejs"
	SUB_RUST    = "rust"
	SUB_CODE    = "vscode"
	SUB_NVIM    = "neovim"
	SUB_BREW    = "homebrew"
	SUB_VLANG   = "vlang"
	SUB_FLUTTER = "flutter"
	SUB_LUA     = "lua"
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

type EnvsHandler struct {
	shellName  string
	rcFilePath string
	oldContent []byte
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
	return
}

func getShellTypeForUnix() (st string) {
	s := os.Getenv("SHELL")
	if strings.Contains(s, Zsh) {
		st = Zsh
	} else if strings.Contains(s, Bash) {
		st = Bash
	} else {
		panic("[unsupported shell] please use zsh or bash.")
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
var (
	PwSetPathEnv  string = `[Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";%s", "User")`
	PwSetOtherEnv string = `[Environment]::SetEnvironmentVariable("%s", "%s", "User")`
)

func (that *EnvsHandler) flushEnvForWin(key string) {
	cmd := exec.Command(fmt.Sprintf("$env:%s", key))
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (that *EnvsHandler) setOneEnvForWin(key, value string) {
	var arg string
	_key := strings.ToLower(key)
	if strings.HasSuffix(_key, "path") && strings.Contains(os.Getenv("PATH"), value) {
		fmt.Printf("[%s] Already exists in Path.\n", value)
		return
	}
	if strings.HasPrefix(_key, "path") && !strings.Contains(_key, "_") {
		arg = fmt.Sprintf(PwSetPathEnv, value)
	} else {
		arg = fmt.Sprintf(PwSetOtherEnv, key, value)
	}
	cmd := exec.Command("powershell", arg)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("[Set env failed]", err, key)
		return
	}
	that.flushEnvForWin(key)
}

func (that *EnvsHandler) SetEnvForWin(envList map[string]string) {
	for key, value := range envList {
		that.setOneEnvForWin(key, value)
	}
	time.Sleep(time.Second * 3)
	that.HintsForWin()
}

var HintStr string = ` +-++-++-++-++-++-++-+
|W||a||r||n||i||n||g|
+-++-++-++-++-++-++-+`

func (that *EnvsHandler) HintsForWin(flag ...int) {
	if runtime.GOOS == Windows {
		fmt.Println(HintStr)
		if len(flag) > 0 {
			fmt.Println("[**WARNING**] Make sure you have PowerShell installed.")
			fmt.Println("[**注意**] Windows用户需要使用PowerShell! 请自行检查系统是否自带或者手动安装PowerShell。")
		} else {
			fmt.Println("[**WARNING**] You have to exit current PowerShell and enter another one to make envs work properly.")
			fmt.Println("[**注意**] Windows用户需要关闭当前PowerShell, 然后重开一个, 环境变量才能生效!否则当前设置过的[Path环境变量会被后面的设置操作覆盖]!!! ")
		}
	}
}
