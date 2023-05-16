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

	color "github.com/TwiN/go-color"
	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/query"
	"github.com/moqsien/gvc/pkgs/utils"
)

type PyVenv struct {
	Conf      *config.GVConfig
	pyenvPath string
	env       *utils.EnvsHandler
	fetcher   *query.Fetcher
}

func NewPyVenv() (py *PyVenv) {
	py = &PyVenv{
		Conf:    config.New(),
		fetcher: query.NewFetcher(),
		env:     utils.NewEnvsHandler(),
	}
	py.initeDirs()
	py.env.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *PyVenv) initeDirs() {
	if ok, _ := utils.PathIsExist(config.PythonToolsPath); !ok {
		if err := os.MkdirAll(config.PythonToolsPath, os.ModePerm); err != nil {
			fmt.Println(color.InRed("[mkdir Failed] "), err)
		}
	}
	if ok, _ := utils.PathIsExist(config.PyenvInstallDir); !ok {
		if err := os.MkdirAll(config.PyenvInstallDir, os.ModePerm); err != nil {
			fmt.Println(color.InRed("[mkdir Failed] "), err)
		}
	}
	if runtime.GOOS != utils.Windows {
		if ok, _ := utils.PathIsExist(config.PyenvCacheDir); !ok {
			if err := os.MkdirAll(config.PyenvCacheDir, os.ModePerm); err != nil {
				fmt.Println(color.InRed("[mkdir Failed] "), err)
			}
		}
		if ok, _ := utils.PathIsExist(config.PythonBinaryPath); !ok {
			if err := os.MkdirAll(config.PythonBinaryPath, os.ModePerm); err != nil {
				fmt.Println(color.InRed("[mkdir Failed] "), err)
			}
		}
		if ok, _ := utils.PathIsExist(config.PyenvVersionsPath); !ok {
			if err := os.MkdirAll(config.PyenvVersionsPath, os.ModePerm); err != nil {
				fmt.Println(color.InRed("[mkdir Failed] "), err)
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
		that.fetcher.Url = that.Conf.Python.PyenvWin
	} else {
		that.fetcher.Url = that.Conf.Python.PyenvUnix
	}
	if that.fetcher.Url != "" {
		if strings.Contains(that.fetcher.Url, "github.com") {
			that.fetcher.Url = that.Conf.Github.GetDownUrl(that.fetcher.Url)
		}
		that.fetcher.Timeout = 10 * time.Second
		fPath := filepath.Join(config.PythonToolsPath, "pyenv-master.zip")
		if flag {
			os.RemoveAll(fPath)
		}
		if size := that.fetcher.GetAndSaveFile(fPath); size != 0 {
			if err := archiver.Unarchive(fPath, config.PyenvInstallDir); err != nil {
				fmt.Println(color.InRed("[unarchive pyenv failed.]"))
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
				fmt.Println(color.InRed("[Cannot set env for Pyenv]"))
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
		fmt.Println(color.InYellow("This would delete the python versions you have installed, still continue? [Y/N]"))
		var r string
		fmt.Scan(&r)
		r = strings.ToLower(r)
		if r != "y" && r != "yes" {
			fmt.Println(color.InGreen("Aborted."))
			return
		}
	}
	fmt.Println(color.InYellow("Update pyenv..."))
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
	if output, err := utils.ExecuteSysCommand(true, that.getExecutablePath(), "install", "--list"); err == nil {
		result := output.String()
		if result == "" {
			return
		}
		var rList []string
		if strings.Contains(result, "\r") {
			rList = strings.Split(result, "\r")
		} else {
			rList = strings.Split(result, "\n")
		}
		newList := []string{}
		for _, v := range rList {
			v = strings.Trim(strings.Trim(v, "\n"), "\r")
			if strings.Contains(v, ":") || strings.Contains(v, "[") {
				continue
			}
			newList = append(newList, v)
		}
		fmt.Println(color.InGreen(strings.Join(newList, "  ")))
	}
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
	that.fetcher.Url = fmt.Sprintf("%s%s/%s", cUrl, version, name)
	that.fetcher.Timeout = 15 * time.Minute
	fpath := filepath.Join(config.PyenvCacheDir, name)
	that.fetcher.GetAndSaveFile(fpath)
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
		that.fetcher.Url = rUrl
		that.fetcher.Timeout = 15 * time.Minute
		sList := strings.Split(that.fetcher.Url, "/")
		fpath := filepath.Join(config.PyenvCacheDir, sList[len(sList)-1])
		if size := that.fetcher.GetAndSaveFile(fpath); size == 0 {
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
		that.fetcher.Url = that.Conf.Python.PyenvWinNeeded
		that.fetcher.Timeout = 15 * time.Minute
		if size := that.fetcher.GetAndSaveFile(fpath); size > 0 {
			untarfilePath := filepath.Join(config.PyenvCacheDir, version)
			if ok, _ := utils.PathIsExist(untarfilePath); !ok {
				os.MkdirAll(untarfilePath, 0666)
			}
			if err := archiver.Unarchive(fpath, untarfilePath); err != nil {
				utils.ClearDir(untarfilePath)
				os.Remove(fpath)
				fmt.Println(color.InRed("[Unarchive failed] "), err)
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
			fmt.Println(color.InGreen(fmt.Sprintf("[**] Download cache file from %s", cUrl)))
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
	fmt.Println(color.InGreen("Python versions are installed in: "))
	fmt.Println(color.InYellow(config.PyenvVersionsPath))
}

func (that *PyVenv) setPipAcceleration() {
	p := config.GetPipConfPath()
	pDir := filepath.Dir(p)
	if ok, _ := utils.PathIsExist(p); !ok {
		if ok, _ := utils.PathIsExist(pDir); !ok {
			if err := os.MkdirAll(pDir, os.ModePerm); err != nil {
				fmt.Println(color.InRed("[mkdir Failed] "), err)
				return
			}
		}
		pUrl := that.Conf.Python.PypiProxies[0]
		parser, _ := url.Parse(pUrl)
		content := fmt.Sprintf(config.PipConfig, pUrl, parser.Host)
		os.WriteFile(p, []byte(content), 0644)
	}
}
