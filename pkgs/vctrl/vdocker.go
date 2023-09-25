package vctrl

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type VDocker struct {
	Conf    *config.GVConfig
	fetcher *request.Fetcher
}

func NewVDocker() (vd *VDocker) {
	vd = &VDocker{
		Conf:    config.New(),
		fetcher: request.NewFetcher(),
	}
	return
}

func (that *VDocker) installDockerForWindows() {
	that.fetcher.SetUrl(that.Conf.Docker.WindowsDockerDownloadUrl)
	that.fetcher.Timeout = time.Minute * 30
	that.fetcher.SetThreadNum(4)
	fPath := filepath.Join(config.DockerFilesDir, "docker-desktop-installer.exe")
	os.RemoveAll(fPath)
	that.fetcher.GetAndSaveFile(fPath, true)
	if ok, _ := utils.PathIsExist(fPath); ok {
		if ok, _ := utils.PathIsExist(config.DockerWindowsInstallationDir); ok {
			os.RemoveAll(config.DockerWindowsInstallationDir)
		}
		_, err := utils.ExecuteSysCommand(false,
			fPath,
			"install",
			"--quiet",
			"--accept-license",
			"--backend=wsl-2",
			fmt.Sprintf("--installation-dir=%s", config.DockerWindowsInstallationDir),
		)
		if err != nil {
			gprint.PrintError("%+v", err)
		} else {
			u, _ := user.Current()
			userNameList := strings.Split(u.Username, `\`)
			if len(userNameList) == 0 {
				return
			}
			uname := userNameList[len(userNameList)-1]
			_, err = utils.ExecuteSysCommand(false,
				"net",
				"localgroup",
				"docker-users",
				uname,
				"/add",
			)
			if err != nil {
				gprint.PrintWarning("< net localgroup docker-users <user> /add > errored: %+v", err)
			}
		}
		os.RemoveAll(fPath)
	}
}

func (that *VDocker) installDockerForLinux() {
	fPath := filepath.Join(config.DockerFilesDir, "install-docker.sh")
	that.fetcher.SetUrl(that.Conf.Docker.LinuxDockerInstallShellScript)
	that.fetcher.Timeout = 3 * time.Minute
	that.fetcher.GetAndSaveFile(fPath, true)
	if ok, _ := utils.PathIsExist(fPath); ok {
		_, err := utils.ExecuteSysCommand(false, "sudo", "sh", fPath)
		if err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *VDocker) installDockerForMacOS() {
	cmdArgs := strings.Split(that.Conf.Docker.MacOSDockerInstallUsingBrew, " ")
	_, err := utils.ExecuteSysCommand(false, cmdArgs...)
	if err != nil {
		gprint.PrintError("%+v", err)
	}
}

func (that *VDocker) Install() {
	if runtime.GOOS == utils.MacOS {
		that.installDockerForMacOS()
	} else if runtime.GOARCH == utils.Linux {
		that.installDockerForLinux()
	} else {
		that.installDockerForWindows()
	}
}

func (that *VDocker) ShowRegistryMirrorInChina() {
	mirrors := []string{
		"https://docker.mirrors.ustc.edu.cn",
		"https://registry.docker-cn.com",
		"http://hub-mirror.c.163.com",
		"https://mirror.ccs.tencentyun.com",
	}
	for idx, mirror := range mirrors {
		fmt.Printf("%d. %s\n", idx, mirror)
	}
}
