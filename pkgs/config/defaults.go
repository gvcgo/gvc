package config

import (
	"fmt"
	"path/filepath"

	"github.com/moqsien/gvc/pkgs/utils"
)

var GVCWorkDir = filepath.Join(utils.GetHomeDir(), ".gvc/")
var GoFilesDir = filepath.Join(GVCWorkDir, "go_files")

var DefaultConfigPath = filepath.Join(GVCWorkDir, "config.yml")

var (
	DefaultGoRoot  string = filepath.Join(GoFilesDir, "go")
	DefaultGoPath  string = filepath.Join(utils.GetHomeDir(), "data/projects/go")
	DefaultGoProxy string = "https://goproxy.cn,direct"
)

var (
	GoTarFilesPath   string = filepath.Join(GoFilesDir, "downloads")
	GoUnTarFilesPath string = filepath.Join(GoFilesDir, "versions")
)

const (
	HostFilePathForNix = "/etc/hosts"
	HostFilePathForWin = `C:\Windows\System32\drivers\etc\hosts`
)

var (
	GoEnvsPattern string = `# Golang Start
export GOROOT="%s"
export GOPATH="%s"
export GOBIN="%s"
export GOPROXY="%s"
export PATH="%s"
# Golang End`
	GoEnvsReg string = `# Golang Start[\s\S]# Golang End`
	EnvPath          = "$PATH:%s"
)

var GoEnv string = fmt.Sprintf(GoEnvsPattern,
	DefaultGoRoot,
	DefaultGoPath,
	filepath.Join(DefaultGoPath, "bin"),
	`%s`,
	`%s`)
