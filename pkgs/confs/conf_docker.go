package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type DockerConf struct {
	LinuxDockerInstallShellScript string `koanf:"docker_install_script"`
	MacOSDockerInstallUsingBrew   string `koanf:"mac_docker_cmd"`
	WindowsDockerDownloadUrl      string `koanf:"windows_docker_download"`
	path                          string
}

func NewDockerConf() (d *DockerConf) {
	d = &DockerConf{
		path: DockerFilesDir,
	}
	d.setup()
	return
}

func (that *DockerConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *DockerConf) Reset() {
	that.LinuxDockerInstallShellScript = "https://test.docker.com"
	that.MacOSDockerInstallUsingBrew = `brew install --cask --appdir=/Applications docker`
	that.WindowsDockerDownloadUrl = "https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe"
}
