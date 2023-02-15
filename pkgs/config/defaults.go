package config

import (
	"path/filepath"

	"github.com/moqsien/gvc/pkgs/utils"
)

var GVCWorkDir = filepath.Join(utils.GetHomeDir(), ".gvc/")

var DefaultConfigPath = filepath.Join(GVCWorkDir, "config.yml")

var (
	DefaultGoRoot  string = filepath.Join(GVCWorkDir, "go_files/go")
	DefaultGoPath  string = filepath.Join(utils.GetHomeDir(), "go")
	DefaultGoProxy string = "https://goproxy.cn,direct"
)

var (
	GoTarFilesPath   string = filepath.Join(GVCWorkDir, "go_files/downloads")
	GoUnTarFilesPath string = filepath.Join(GVCWorkDir, "go_files/versions")
)

const (
	HostFilePathForNix = "/etc/hosts"
	HostFilePathForWin = `C:\Windows\System32\drivers\etc\hosts`
)
