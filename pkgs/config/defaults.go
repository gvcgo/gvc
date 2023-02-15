package config

import (
	"path/filepath"

	"github.com/moqsien/gvc/pkgs/utils"
)

var GVCWorkDir = filepath.Join(utils.GetHomeDir(), ".gvc/")
var GoFilesDir = filepath.Join(GVCWorkDir, "go_files")

var DefaultConfigPath = filepath.Join(GVCWorkDir, "config.yml")

var (
	DefaultGoRoot  string = filepath.Join(GoFilesDir, "go")
	DefaultGoPath  string = filepath.Join(utils.GetHomeDir(), "go")
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
