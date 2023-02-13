package config

import (
	"path/filepath"

	"github.com/moqsien/gvc/pkgs/utils"
)

var DefaultConfigPath = filepath.Join(utils.GetHomeDir(), ".gvc/config.yml")

var (
	DefaultGoRoot  string = filepath.Join(utils.GetHomeDir(), ".gvc/go/")
	DefaultGoPath  string = filepath.Join(utils.GetHomeDir(), "go")
	DefaultGoProxy string = "https://goproxy.cn,direct"
)

var (
	GoTarFilesPath   string = filepath.Join(utils.GetHomeDir(), ".gvc/tarfiles/")
	GoUnTarFilesPath string = filepath.Join(utils.GetHomeDir(), ".gvc/untarfiles/")
	GoInstalled      string = filepath.Join(utils.GetHomeDir(), ".gvc/version.yml")
)

const (
	HostFilePathForNix = "/etc/hosts"
	HostFilePathForWin = `C:\Windows\System32\drivers\etc\hosts`
)
