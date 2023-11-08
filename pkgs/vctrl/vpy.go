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
	myArchiver "github.com/moqsien/goutils/pkgs/archiver"
	"github.com/moqsien/goutils/pkgs/gtea/confirm"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type PyVenv struct {
	Conf      *config.GVConfig
	pyenvPath string
	env       *utils.EnvsHandler
	fetcher   *request.Fetcher
}

func NewPyVenv() (py *PyVenv) {
	py = &PyVenv{
		Conf:    config.New(),
		fetcher: request.NewFetcher(),
		env:     utils.NewEnvsHandler(),
	}
	py.initeDirs()
	py.env.SetWinWorkDir(config.GVCDir)
	return
}

func (that *PyVenv) initeDirs() {
	utils.MakeDirs(config.PythonToolsPath, config.PyenvInstallDir)
	if runtime.GOOS != utils.Windows {
		utils.MakeDirs(config.PyenvCacheDir, config.PythonBinaryPath, config.PyenvVersionsPath)
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
	cfm := confirm.NewConfirm(confirm.WithTitle("Set download accelerations in China or not?"))
	cfm.Run()
	result := cfm.Result()
	if !result {
		return
	}
	if runtime.GOOS == utils.Windows {
		rootDir := config.GetPyenvRootPath()
		versionFilePath := filepath.Join(rootDir, ".versions_cache.xml")
		if ok, _ := utils.PathIsExist(versionFilePath); ok {
			content, _ := os.ReadFile(versionFilePath)
			newStr := strings.ReplaceAll(string(content), config.PyenvWinOriginalPyUrl, config.PyenvWinTaobaoPyUrl)
			os.WriteFile(versionFilePath, []byte(newStr), os.ModePerm)
		}
		fpath := filepath.Join(pyenvDir, "pyenv-win", "libexec", "pyenv-install.vbs")
		utils.ReplaceFileContent(fpath,
			config.PyenvWinBeforeFixed,
			config.PyenvWinAfterFixed,
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
		gprint.PrintInfo("Pyenv already installed.")
		return
	}
	if runtime.GOOS == utils.Windows {
		that.fetcher.Url = that.Conf.Python.PyenvWin
	} else {
		that.fetcher.Url = that.Conf.Python.PyenvUnix
	}
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(that.fetcher.Url)

	if that.fetcher.Url != "" {
		that.fetcher.Timeout = 20 * time.Minute
		fPath := filepath.Join(config.PythonToolsPath, "pyenv-master.zip")
		if flag {
			os.RemoveAll(fPath)
		}
		if size := that.fetcher.GetAndSaveFile(fPath); size != 0 {
			if err := archiver.Unarchive(fPath, config.PyenvInstallDir); err != nil {
				os.RemoveAll(fPath)
				os.RemoveAll(config.PyenvInstallDir)
				gprint.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
				return
			}
			that.handlePyenvUntarfile()
			pDir := filepath.Join(config.PyenvInstallDir, "pyenv")
			that.getPyenvPath(pDir)
			if that.pyenvPath != "" {
				that.setEnv()
			} else {
				gprint.PrintError("Cannot set env for Pyenv.")
			}
			that.modifyAccelertion(pDir)
			that.setAssetDirForPyenvWin()
		}
	}
}

func (that *PyVenv) setAssetDirForPyenvWin() {
	if runtime.GOOS == utils.Windows {
		fpath := filepath.Join(config.GetPyenvRootPath(), "libexec", "libs", "pyenv-lib.vbs")

		newCachePath := filepath.Join(config.PythonFilesDir, "install_cache")
		os.MkdirAll(newCachePath, os.ModePerm)
		utils.ReplaceFileContent(fpath,
			config.PyenvWinOriginalCacheDir,
			fmt.Sprintf(config.PyenvWinNewCacheDir, newCachePath),
			0777)

		newVersionPath := filepath.Join(config.PythonFilesDir, "versions")
		os.MkdirAll(newVersionPath, os.ModePerm)
		utils.ReplaceFileContent(fpath,
			config.PyenvWinOriginalVersionsDir,
			fmt.Sprintf(config.PyenvWinNewVersionsDir, newVersionPath),
			0777)

		batPath := filepath.Join(config.GetPyenvRootPath(), "bin", "pyenv.bat")
		utils.ReplaceFileContent(batPath,
			config.PyenvWinBatBeforeFixed,
			strings.ReplaceAll(config.PyenvWinBatAfterFixed, "$$$", newVersionPath),
			0777)
	}
}

func (that *PyVenv) InstallPyenv() {
	shimsDir := filepath.Join(config.GetPyenvRootPath(), "shims")
	zipDir := config.PythonFilesDir
	zipPath := filepath.Join(config.PythonFilesDir, "old_shims.zip")
	if ok, _ := utils.PathIsExist(shimsDir); ok {
		a, _ := myArchiver.NewArchiver(shimsDir, zipDir)
		a.ZipDir()
	}
	that.getPyenv(true)
	if ok, _ := utils.PathIsExist(zipPath); ok {
		ok, _ = utils.PathIsExist(shimsDir)
		if !ok {
			a, _ := myArchiver.NewArchiver(zipDir, shimsDir)
			a.UnArchive()
		}
	}
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
	gprint.PrintInfo("Updating pyenv...")
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
		fc := gprint.NewFadeColors(newList)
		fc.Println()
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
			os.MkdirAll(config.PyenvCacheDir, os.ModePerm)
		}
		fpath := filepath.Join(config.PyenvCacheDir, "needed.zip")
		that.fetcher.Url = that.Conf.Python.PyenvWinNeeded
		that.fetcher.Timeout = 15 * time.Minute
		if size := that.fetcher.GetAndSaveFile(fpath); size > 0 {
			untarfilePath := filepath.Join(config.PyenvCacheDir, version)
			if ok, _ := utils.PathIsExist(untarfilePath); !ok {
				os.MkdirAll(untarfilePath, os.ModePerm)
			}
			if err := archiver.Unarchive(fpath, untarfilePath); err != nil {
				utils.ClearDir(untarfilePath)
				os.Remove(fpath)
				gprint.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
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
			gprint.PrintInfo(fmt.Sprintf("Download cache file from %s", cUrl))
			that.downloadCache(version, cUrl)
		}
		that.getReadlineForUnix()
		if useDefault {
			that.getInstallNeededForWin(version)
		}
		utils.ExecuteSysCommand(false, that.getExecutablePath(), "install", version)
	}
	utils.ExecuteSysCommand(false, that.getExecutablePath(), "global", version)
	that.setPipAcceleration()
}

func (that *PyVenv) RemoveVersion(version string) {
	that.getPyenv()
	that.setTempEnvs()
	utils.ExecuteSysCommand(false, that.getExecutablePath(), "uninstall", version)
}

func (that *PyVenv) ShowInstalled() {
	that.getPyenv()
	that.setTempEnvs()
	utils.ExecuteSysCommand(false, that.getExecutablePath(), "versions")
}

func (that *PyVenv) ShowVersionPath() {
	fc := gprint.NewFadeColors(fmt.Sprintf("Python versions are installed in: %s", config.PyenvVersionsPath))
	fc.Println()
}

func (that *PyVenv) setPipAcceleration() {
	p := config.GetPipConfPath()
	pDir := filepath.Dir(p)
	if ok, _ := utils.PathIsExist(p); !ok {
		if ok, _ := utils.PathIsExist(pDir); !ok {
			if err := os.MkdirAll(pDir, os.ModePerm); err != nil {
				gprint.PrintError("%+v", err)
				return
			}
		}
		pUrl := that.Conf.Python.PypiProxies[0]
		parser, _ := url.Parse(pUrl)
		content := fmt.Sprintf(config.PipConfig, pUrl, parser.Host)
		os.WriteFile(p, []byte(content), 0644)
	}
}

func (that *PyVenv) FixSystemGenerationsForWin() {
	if runtime.GOOS == utils.Windows {
		if homeDir, err := os.UserHomeDir(); err == nil {
			targetDir := filepath.Join(homeDir, `AppData\Local\Microsoft\WindowsApps`)
			pathList := []string{
				filepath.Join(targetDir, "python"),
				filepath.Join(targetDir, "python3"),
				filepath.Join(targetDir, "python.exe"),
				filepath.Join(targetDir, "python3.exe"),
			}
			for _, p := range pathList {
				os.RemoveAll(p)
			}
		}
	}
}
