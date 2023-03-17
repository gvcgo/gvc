package vctrl

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
)

type PyVenv struct {
	*downloader.Downloader
	Conf      *config.GVConfig
	pyenvPath string
	// exePath string
}

func NewPyVenv() (py *PyVenv) {
	py = &PyVenv{
		Conf:       config.New(),
		Downloader: &downloader.Downloader{},
	}
	py.initeDirs()
	return
}

func (that *PyVenv) initeDirs() {
	if ok, _ := utils.PathIsExist(config.PythonToolsPath); !ok {
		if err := os.MkdirAll(config.PythonToolsPath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.PythonRootPath); !ok {
		if err := os.MkdirAll(config.PythonRootPath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.PyenvVersionsPath); !ok {
		if err := os.MkdirAll(config.PyenvVersionsPath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.PyenvInstallDir); !ok {
		if err := os.MkdirAll(config.PyenvInstallDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.PyenvCacheDir); !ok {
		if err := os.MkdirAll(config.PyenvCacheDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
}

func (that *PyVenv) handlePyenvUntarfile() {
	if fList, err := os.ReadDir(config.PyenvInstallDir); err == nil {
		dirName := "pyenv"
		p := filepath.Join(config.PyenvInstallDir, dirName)
		for _, f := range fList {
			if f.IsDir() && f.Name() == dirName {
				os.RemoveAll(p)
			}
		}

		if len(fList) == 1 {
			os.Rename(filepath.Join(config.PyenvInstallDir, fList[0].Name()), p)
		}

		if len(fList) == 2 {
			if ok, _ := utils.PathIsExist(p); ok {
				os.RemoveAll(p)
			}
			for _, f := range fList {
				if f.IsDir() && f.Name() == dirName {
					os.Rename(filepath.Join(config.PyenvInstallDir, f.Name()), p)
				}
			}
		}
	}
}

func (that *PyVenv) getPyenvPath(p string) {
	if fList, err := os.ReadDir(p); err == nil {
		for _, d := range fList {
			if d.IsDir() {
				that.getPyenvPath(filepath.Join(p, d.Name()))
			} else {
				if d.Name() == "pyenv" && strings.Contains(p, "bin") {
					that.pyenvPath = p
				}
			}
		}
	}
}

func (that *PyVenv) setEnv() {
	if runtime.GOOS == utils.Windows {
		utils.SetWinEnv(config.GetPyenvRootEnvName(), config.PyenvRootPath)
		utils.SetWinEnv("Path", fmt.Sprintf("%s;%s",
			that.pyenvPath, config.PythonRootPath))
	} else {
		envars := fmt.Sprintf(config.PythonUnixEnvPattern,
			config.GetPyenvRootEnvName(),
			config.PyenvRootPath,
			that.pyenvPath,
			config.PythonRootPath)
		utils.SetUnixEnv(envars)
	}
}

func (that *PyVenv) getPyenv(force ...bool) {
	flag := false
	if len(force) > 0 {
		flag = force[0]
	}

	if !flag && that.getExecutable() != "" {
		fmt.Println("pyenv already installed.")
		return
	}
	if runtime.GOOS == utils.Windows {
		that.Url = that.Conf.Python.PyenvWin
	} else {
		that.Url = that.Conf.Python.PyenvUnix
	}
	if that.Url != "" {
		that.Url = that.Conf.Github.GetDownUrl(that.Url)
		that.Timeout = 10 * time.Second
		fPath := filepath.Join(config.PythonToolsPath, "pyenv-master.zip")
		if flag {
			os.RemoveAll(fPath)
		}
		if size := that.GetFile(fPath, os.O_CREATE|os.O_WRONLY, 0644); size != 0 {
			if err := archiver.Unarchive(fPath, config.PyenvInstallDir); err != nil {
				os.RemoveAll(fPath)
				os.RemoveAll(config.PyenvInstallDir)
				return
			}
			that.handlePyenvUntarfile()
			that.getPyenvPath(filepath.Join(config.PyenvInstallDir, "pyenv"))
			if that.pyenvPath != "" {
				that.setEnv()
			} else {
				fmt.Println("[Cannot set env for Pyenv]")
			}
		}
	}
}

func (that *PyVenv) InstallPyenv() {
	that.getPyenv(true)
}

func (that *PyVenv) getExecutable() (exePath string) {
	that.getPyenvPath(filepath.Join(config.PyenvInstallDir, "pyenv"))
	if that.pyenvPath != "" {
		exePath = filepath.Join(that.pyenvPath, "pyenv")
	}
	return
}

func (that *PyVenv) UpdatePyenv() {
	that.getPyenv(true)
}

func (that *PyVenv) setTempEnvs() {
	os.Setenv(config.GetPyenvRootEnvName(), config.PyenvRootPath)
	os.Setenv("PYTHON_BUILD_MIRROR_URL", that.Conf.Python.PyBuildUrls[0])
}

func (that *PyVenv) ListRemoteVersions() {
	that.getPyenv()
	that.setTempEnvs()
	utils.ExecuteCommand(that.getExecutable(), "install", "--list")
}

func (that *PyVenv) isInstalled(version string) (r bool) {
	cmd := exec.Command(that.getExecutable(), "versions")
	cmd.Env = os.Environ()
	output, _ := cmd.CombinedOutput()
	if strings.Contains(string(output), version) {
		r = true
	}
	return
}

func (that *PyVenv) downloadCache(version string) {
	name := fmt.Sprintf("Python-%s.tar.xz", version)
	that.Url = fmt.Sprintf("%s%s/%s", that.Conf.Python.PyBuildUrls[0], version, name)
	that.Timeout = 10 * time.Minute
	fpath := filepath.Join(config.PyenvCacheDir, name)
	that.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644)
}

func (that *PyVenv) InstallVersion(version string) {
	that.getPyenv()
	that.setTempEnvs()
	if !that.isInstalled(version) {
		that.downloadCache(version)
		utils.ExecuteCommand(that.getExecutable(), "install", version)
	}
	utils.ExecuteCommand(that.getExecutable(), "global", version)
	that.setPipAcceleration()
}

func (that *PyVenv) RemoveVersion(version string) {
	that.getPyenv()
	that.setTempEnvs()
	utils.ExecuteCommand(that.getExecutable(), "uninstall", version)
}

func (that *PyVenv) ShowInstalled() {
	that.getPyenv()
	that.setTempEnvs()
	utils.ExecuteCommand(that.getExecutable(), "versions")
}

func (that *PyVenv) ShowVersionPath() {
	fmt.Println("Python versions are installed in: ")
	fmt.Println(config.PyenvVersionsPath)
}

func (that *PyVenv) setPipAcceleration() {
	p := config.GetPipConfPath()
	pDir := filepath.Dir(p)
	if ok, _ := utils.PathIsExist(p); !ok {
		if ok, _ := utils.PathIsExist(pDir); !ok {
			if err := os.MkdirAll(pDir, os.ModePerm); err != nil {
				fmt.Println("[mkdir Failed] ", err)
				return
			}
		}
		pUrl := that.Conf.Python.PypiProxies[0]
		parser, _ := url.Parse(pUrl)
		content := fmt.Sprintf(config.PipConfig, pUrl, parser.Host)
		os.WriteFile(p, []byte(content), 0644)
	}
}
