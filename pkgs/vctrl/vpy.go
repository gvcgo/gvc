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
	downloader "github.com/moqsien/gvc/pkgs/fetcher"
	"github.com/moqsien/gvc/pkgs/utils"
)

type PyVenv struct {
	*downloader.Downloader
	Conf      *config.GVConfig
	pyenvPath string
	env       *utils.EnvsHandler
}

func NewPyVenv() (py *PyVenv) {
	py = &PyVenv{
		Conf:       config.New(),
		Downloader: &downloader.Downloader{},
		env:        utils.NewEnvsHandler(),
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
	if ok, _ := utils.PathIsExist(config.PyenvInstallDir); !ok {
		if err := os.MkdirAll(config.PyenvInstallDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if runtime.GOOS != utils.Windows {
		if ok, _ := utils.PathIsExist(config.PyenvCacheDir); !ok {
			if err := os.MkdirAll(config.PyenvCacheDir, os.ModePerm); err != nil {
				fmt.Println("[mkdir Failed] ", err)
			}
		}
		if ok, _ := utils.PathIsExist(config.PythonBinaryPath); !ok {
			if err := os.MkdirAll(config.PythonBinaryPath, os.ModePerm); err != nil {
				fmt.Println("[mkdir Failed] ", err)
			}
		}
		if ok, _ := utils.PathIsExist(config.PyenvVersionsPath); !ok {
			if err := os.MkdirAll(config.PyenvVersionsPath, os.ModePerm); err != nil {
				fmt.Println("[mkdir Failed] ", err)
			}
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

		if len(fList) == 1 && fList[0].IsDir() {
			os.Rename(filepath.Join(config.PyenvInstallDir, fList[0].Name()), p)
		}

		if len(fList) >= 2 {
			if ok, _ := utils.PathIsExist(p); ok {
				os.RemoveAll(p)
			}
			for _, f := range fList {
				if f.IsDir() && f.Name() != dirName {
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
	if runtime.GOOS == utils.Windows {
		that.pyenvPath = filepath.Join(config.PyenvInstallDir, "pyenv/pyenv-win/bin")
	}
}

func (that *PyVenv) setEnv() {
	if runtime.GOOS == utils.Windows {
		envList := map[string]string{
			config.PyenvRootName: config.PyenvRootPath,
			"PATH":               fmt.Sprintf("%s;%s", that.pyenvPath, config.PythonBinaryPath),
		}
		that.env.SetEnvForWin(envList)
	} else {
		pyEnv := fmt.Sprintf(utils.PyEnv,
			config.PyenvRootName,
			config.PyenvRootPath,
			that.pyenvPath,
			config.PythonBinaryPath)
		that.env.UpdateSub(utils.SUB_PY, pyEnv)
	}
}

func (that *PyVenv) modifyAccelertion(pyenvDir string) {
	if runtime.GOOS == utils.Windows {
		fpath := filepath.Join(pyenvDir, "pyenv-win/libexec/pyenv-install.vbs")
		utils.ReplaceFileContent(fpath,
			config.PyenvModifyForwin1,
			config.PyenvAfterModifyWin1,
			0777)
		utils.ReplaceFileContent(fpath,
			config.PyenvModifyForwin2,
			config.PyenvAfterModifyWin2,
			0777)
	} else {
		fpath := filepath.Join(pyenvDir, "plugins/python-build/bin/python-build")
		utils.ReplaceFileContent(fpath,
			config.PyenvModifyForUnix,
			config.PyenvAfterModifyUnix,
			0777)
	}
}

func (that *PyVenv) getPyenv(force ...bool) {
	flag := false
	if len(force) > 0 {
		flag = force[0]
	}

	if !flag && that.getExecutablePath() != "" {
		fmt.Println("pyenv already installed.")
		return
	}
	if runtime.GOOS == utils.Windows {
		that.Url = that.Conf.Python.PyenvWin
	} else {
		that.Url = that.Conf.Python.PyenvUnix
	}
	if that.Url != "" {
		if strings.Contains(that.Url, "github.com") {
			that.Url = that.Conf.Github.GetDownUrl(that.Url)
		}
		that.Timeout = 10 * time.Second
		fPath := filepath.Join(config.PythonToolsPath, "pyenv-master.zip")
		if flag {
			os.RemoveAll(fPath)
		}
		if size := that.GetFile(fPath, os.O_CREATE|os.O_WRONLY, 0644); size != 0 {
			if err := archiver.Unarchive(fPath, config.PyenvInstallDir); err != nil {
				fmt.Println("[unarchive pyenv failed.]")
				os.RemoveAll(fPath)
				os.RemoveAll(config.PyenvInstallDir)
				return
			}
			that.handlePyenvUntarfile()
			pDir := filepath.Join(config.PyenvInstallDir, "pyenv")
			that.getPyenvPath(pDir)
			if that.pyenvPath != "" {
				that.setEnv()
			} else {
				fmt.Println("[Cannot set env for Pyenv]")
			}
			that.modifyAccelertion(pDir)
		}
	}
}

func (that *PyVenv) InstallPyenv() {
	that.getPyenv(true)
}

func (that *PyVenv) getExecutablePath() (exePath string) {
	p := filepath.Join(config.PyenvInstallDir, "pyenv")
	that.getPyenvPath(p)
	if that.pyenvPath != "" {
		exePath = filepath.Join(that.pyenvPath, "pyenv")
		if ok, _ := utils.PathIsExist(exePath); !ok {
			exePath = ""
		}
	}
	return
}

func (that *PyVenv) UpdatePyenv() {
	if runtime.GOOS == utils.Windows {
		fmt.Println("This would delete the python versions you have installed, still continue?[Y/N]")
		var r string
		fmt.Scan(&r)
		r = strings.ToLower(r)
		if r != "y" && r != "yes" {
			fmt.Println("Aborted.")
			return
		}
	}
	fmt.Println("Update pyenv...")
	that.getPyenv(true)
}

func (that *PyVenv) setTempEnvs() {
	os.Setenv(config.PyenvRootName, config.PyenvRootPath)
	os.Setenv(config.PyenvMirrorEnvName, that.Conf.Python.PyBuildUrl)
	os.Setenv(config.PyenvMirrorEnabledName, "true")
	if ok, _ := utils.PathIsExist(config.PythonBinaryPath); ok && runtime.GOOS == utils.Windows {
		vPath := os.Getenv("PATH")
		if !strings.Contains(vPath, config.PythonBinaryPath) {
			os.Setenv("PATH", fmt.Sprintf("%s;%s", vPath, config.PythonBinaryPath))
		}
	}
}

func (that *PyVenv) ListRemoteVersions() {
	that.getPyenv()
	that.setTempEnvs()
	utils.ExecuteCommand(that.getExecutablePath(), "install", "--list")
}

func (that *PyVenv) isInstalled(version string) (r bool) {
	cmd := exec.Command(that.getExecutablePath(), "versions")
	cmd.Env = os.Environ()
	output, _ := cmd.CombinedOutput()
	if strings.Contains(string(output), version) {
		r = true
	}
	return
}

func (that *PyVenv) downloadCache(version, cUrl string) {
	name := fmt.Sprintf("Python-%s.tar.xz", version)
	that.Url = fmt.Sprintf("%s%s/%s", cUrl, version, name)
	that.Timeout = 10 * time.Minute
	fpath := filepath.Join(config.PyenvCacheDir, name)
	that.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644)
}

func (that *PyVenv) getReadlineForUnix() {
	if runtime.GOOS == utils.Windows {
		return
	}
	rUrls := that.Conf.Python.PyenvReadline
	if len(rUrls) == 0 {
		return
	}
	for _, rUrl := range rUrls {
		that.Url = rUrl
		that.Timeout = 10 * time.Minute
		sList := strings.Split(that.Url, "/")
		fpath := filepath.Join(config.PyenvCacheDir, sList[len(sList)-1])
		if size := that.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size == 0 {
			os.Remove(fpath)
		}
	}
}

func (that *PyVenv) getInstallNeededForWin(version string) {
	if ok, _ := utils.PathIsExist(config.PyenvRootPath); ok && runtime.GOOS == utils.Windows && runtime.GOARCH == utils.X64 {
		if ok, _ := utils.PathIsExist(config.PyenvCacheDir); !ok {
			os.MkdirAll(config.PyenvCacheDir, 0666)
		}
		fpath := filepath.Join(config.PyenvCacheDir, "needed.zip")
		that.Url = that.Conf.Python.PyenvWinNeeded
		that.Timeout = 6 * time.Minute
		if size := that.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size > 0 {
			untarfilePath := filepath.Join(config.PyenvCacheDir, version)
			if ok, _ := utils.PathIsExist(untarfilePath); !ok {
				os.MkdirAll(untarfilePath, 0666)
			}
			if err := archiver.Unarchive(fpath, untarfilePath); err != nil {
				utils.ClearDir(untarfilePath)
				os.Remove(fpath)
				fmt.Println("[Unarchive failed] ", err)
				return
			}
		}
	}
}

func (that *PyVenv) InstallVersion(version string, useDefault bool) {
	that.getPyenv()
	that.setTempEnvs()
	if !that.isInstalled(version) {
		if runtime.GOOS != utils.Windows && os.Getenv("PYENV_PRE_CACHE") != "" {
			cUrl := that.Conf.Python.PyBuildUrls[0]
			fmt.Println("[**] Download cache file from ", cUrl)
			that.downloadCache(version, cUrl)
		}
		that.getReadlineForUnix()
		if useDefault {
			that.getInstallNeededForWin(version)
		}
		utils.ExecuteCommand(that.getExecutablePath(), "install", version)
	}
	utils.ExecuteCommand(that.getExecutablePath(), "global", version)
	that.setPipAcceleration()
	// that.env.HintsForWin()
}

func (that *PyVenv) RemoveVersion(version string) {
	that.getPyenv()
	that.setTempEnvs()
	utils.ExecuteCommand(that.getExecutablePath(), "uninstall", version)
}

func (that *PyVenv) ShowInstalled() {
	that.getPyenv()
	that.setTempEnvs()
	utils.ExecuteCommand(that.getExecutablePath(), "versions")
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
