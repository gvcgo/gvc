package vctrl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mholt/archiver/v3"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

var (
	CygwinInstallerName string = "cygwin-installer.exe"
	Msys2InstallerName  string = "msys2_installer.exe"
	VCpkgZipName        string = "vcpkig.zip"
	// .\msys2-x86_64-latest.exe in --confirm-command --accept-messages --root C:/msys64
	Msys2Args []string = []string{
		"in",
		"--confirm-command",
		"--accept-messages",
	}
)

type CppManager struct {
	Conf    *config.GVConfig
	Doc     *goquery.Document
	env     *utils.EnvsHandler
	fetcher *request.Fetcher
}

func NewCppManager() (cm *CppManager) {
	cm = &CppManager{
		fetcher: request.NewFetcher(),
		Conf:    config.New(),
		env:     utils.NewEnvsHandler(),
	}
	cm.initDirs()
	cm.env.SetWinWorkDir(config.GVCDir)
	return
}

func (that *CppManager) initDirs() {
	utils.MakeDirs(config.CppFilesDir, config.Msys2Dir, config.VCpkgDir, config.CppDownloadDir, config.CygwinRootDir)
}

func (that *CppManager) getMsys2Installer() (fPath string) {
	fPath = filepath.Join(config.CppDownloadDir, Msys2InstallerName)
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(that.Conf.Cpp.MsysInstallerUrl)
	that.fetcher.Timeout = 10 * time.Minute
	that.fetcher.GetAndSaveFile(fPath)
	return
}

func (that *CppManager) writeScript() (scriptPath string) {
	fPath := filepath.Join(config.CppDownloadDir, "execute_installer.bat")
	if ok, _ := utils.PathIsExist(fPath); !ok {
		content := fmt.Sprintf(`%s in --confirm-command --accept-messages --root %s`, Msys2InstallerName, config.Msys2Dir)
		os.WriteFile(fPath, []byte(content), 0777)
	}
	return fPath
}

func (that *CppManager) InstallMsys2() {
	if runtime.GOOS != utils.Windows {
		return
	}
	fPath := that.getMsys2Installer()
	if ok, _ := utils.PathIsExist(fPath); ok {
		os.Setenv("PATH", fmt.Sprintf("%s;%s", config.CppDownloadDir, os.Getenv("PATH")))
		batPath := that.writeScript()
		c := exec.Command(batPath)
		c.Env = os.Environ()
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		if err := c.Run(); err != nil {
			gprint.PrintError("%+v", err)
			return
		}
		binPath := filepath.Join(config.Msys2Dir, "usr", "bin")
		if ok, _ := utils.PathIsExist(binPath); ok {
			winEnv := map[string]string{
				"PATH": binPath,
			}
			that.env.SetEnvForWin(winEnv)
			winEnv = map[string]string{
				"PATH": config.CppDownloadDir,
			}
			that.env.SetEnvForWin(winEnv)
		}
	}
}

func (that *CppManager) UninstallMsys2() {
	uninstallExe := filepath.Join(config.Msys2Dir, "uninstall.exe")
	if ok, _ := utils.PathIsExist(uninstallExe); ok {
		c := exec.Command(uninstallExe, "--purge")
		c.Env = os.Environ()
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		if err := c.Run(); err != nil {
			gprint.PrintError("%+v", err)
			return
		}
	}
}

/*
use Msys2/Cygwin git for vscode.
Need some fixes.
*/
func (that *CppManager) RepairGitForVSCode() {
	if runtime.GOOS != utils.Windows {
		return
	}
	bPath := filepath.Join(config.CppFilesDir, "mgit.bat")
	if ok, _ := utils.PathIsExist(bPath); !ok {
		os.WriteFile(bPath, []byte(config.Msys2CygwinGitFixBat), 0777)
	}
	bPath = strings.ReplaceAll(bPath, `\`, `\\`)
	cnf := NewGVCWebdav()
	filesToSync := cnf.GetFilesToSync()
	vscodeSettingsPath := filesToSync[config.CodeUserSettingsBackupFileName]
	utils.AddNewlineToVscodeSettings("git.path", bPath, vscodeSettingsPath)
}

func (that *CppManager) getCygwinInstaller() (fPath string) {
	fPath = filepath.Join(config.CppDownloadDir, CygwinInstallerName)
	if ok, _ := utils.PathIsExist(fPath); !ok {
		that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(that.Conf.Cpp.CygwinInstallerUrl)
		if that.fetcher.Url != "" {
			that.fetcher.Timeout = 10 * time.Minute
			if size := that.fetcher.GetAndSaveFile(fPath); size == 0 {
				gprint.PrintError("Download Cygwin installer failed!")
				os.RemoveAll(fPath)
			} else {
				if runtime.GOOS == utils.Windows {
					that.env.SetEnvForWin(map[string]string{
						"PATH": config.CppDownloadDir,
					})
				}
			}
		}
	}
	return
}

func (that *CppManager) InstallCygwin(packInfo string) {
	installerPath := that.getCygwinInstaller()
	if packInfo == "" {
		packInfo = "git,bash,wget,gcc,gdb,clang,openssh,bashdb,gdbm,gcc-fortran,clang-analyzer,clang-doc,bash-completion,bash-devel,bash-completion-cmake"
	}
	gprint.PrintInfo(fmt.Sprintf("Install Packages: %s", packInfo))
	if ok, _ := utils.PathIsExist(installerPath); ok && runtime.GOOS == utils.Windows {
		ePath := os.Getenv("PATH")
		if !strings.Contains(ePath, config.CppDownloadDir) {
			ePath = fmt.Sprintf("%s;%s", config.CppDownloadDir, ePath)
			os.Setenv("PATH", ePath)
		}
		cmd := exec.Command(CygwinInstallerName, "-q", "-f", "-N", "-O", "-s",
			that.Conf.Cpp.CygwinMirrorUrls[0], "-R", config.CygwinRootDir,
			"-P", packInfo)
		cmd.Env = os.Environ()
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			gprint.PrintError("%+v", err)
			return
		}

		if ok, _ := utils.PathIsExist(config.CygwinBinaryDir); !ok {
			that.env.SetEnvForWin(map[string]string{
				"PATH": config.CygwinBinaryDir,
			})
			that.env.SetEnvForWin(map[string]string{
				"PATH": config.CygwinRootDir,
			})
		}
	}
}

var (
	VCpkgFilename     string = "vcpkg.zip"
	VCpkgToolFilename string = "vcpkg-tool.zip"
)

func (that *CppManager) getVCPkg() string {
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(that.Conf.Cpp.VCpkgUrl)
	fPath := filepath.Join(config.CppDownloadDir, VCpkgFilename)
	if size := that.fetcher.GetAndSaveFile(fPath); size != 0 {
		return fPath
	} else {
		os.RemoveAll(fPath)
	}
	return ""
}

func (that *CppManager) getVCPkgTool() string {
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(that.Conf.Cpp.VCpkgToolUrl)
	fPath := filepath.Join(config.CppDownloadDir, VCpkgToolFilename)
	if size := that.fetcher.GetAndSaveFile(fPath); size != 0 {
		return fPath
	} else {
		os.RemoveAll(fPath)
	}
	return ""
}

func (that *CppManager) writeCompileScript(buildPath, srcPath string) (cmd, fPath string) {
	if runtime.GOOS != utils.Windows {
		script := fmt.Sprintf(config.VCPkgScript, buildPath,
			"g++", srcPath, buildPath)
		fPath = filepath.Join(config.CppDownloadDir, "compile_vcpkg.sh")
		os.WriteFile(fPath, []byte(script), 0777)
		cmd = "sh"
	}
	return
}

func (that *CppManager) checkVcpkgCompilationEnv() (hasCompiler, hasCmake bool) {
	if runtime.GOOS == utils.Windows {
		return true, true
	}

	if ok, _ := utils.PathIsExist("/usr/bin/g++"); ok {
		hasCompiler = true
	}
	if ok, _ := utils.PathIsExist("/usr/local/bin/g++"); ok {
		hasCompiler = true
	}

	if ok, _ := utils.PathIsExist("/usr/bin/cmake"); ok {
		hasCmake = true
	}
	if ok, _ := utils.PathIsExist("/usr/local/bin/cmake"); ok {
		hasCmake = true
	}
	if !hasCompiler {
		gprint.Yellow("Please install g++ compiler.")
	}
	if !hasCmake {
		gprint.Yellow("Please install cmake.")
	}
	return
}

func (that *CppManager) InstallVCPkg() {
	hasCompiler, hasCmake := that.checkVcpkgCompilationEnv()
	if !(hasCompiler && hasCmake) {
		return
	}
	fPath := that.getVCPkg()
	if ok, _ := utils.PathIsExist(fPath); ok {
		if err := archiver.Unarchive(fPath, config.CppDownloadDir); err != nil {
			gprint.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
			return
		}
		dirList, _ := os.ReadDir(config.CppDownloadDir)
		os.RemoveAll(config.VCpkgDir)
		for _, d := range dirList {
			if d.IsDir() && (strings.Contains(d.Name(), "vcpkg") && !strings.Contains(d.Name(), "tool")) {
				os.Rename(filepath.Join(config.CppDownloadDir, d.Name()), config.VCpkgDir)
				break
			}
		}
		if ok, _ = utils.PathIsExist(config.VCpkgDir); !ok {
			return
		}

		if runtime.GOOS != utils.Windows {
			fPath = that.getVCPkgTool()
			if ok, _ := utils.PathIsExist(fPath); !ok {
				return
			}
			basePath := filepath.Join(config.VCpkgDir, "buildtrees", "_vcpkg")
			buildPath := filepath.Join(basePath, "build")
			srcPath := filepath.Join(basePath, "src")
			os.MkdirAll(buildPath, os.ModePerm)

			if err := archiver.Unarchive(fPath, config.CppDownloadDir); err != nil {
				gprint.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
				return
			}
			dirList, _ = os.ReadDir(config.CppDownloadDir)
			for _, d := range dirList {
				if d.IsDir() && (strings.Contains(d.Name(), "vcpkg") && strings.Contains(d.Name(), "tool")) {
					os.Rename(filepath.Join(config.CppDownloadDir, d.Name()), srcPath)
					break
				}
			}

			if ok, _ = utils.PathIsExist(srcPath); ok {
				cmdName, scriptPath := that.writeCompileScript(buildPath, srcPath)
				if scriptPath != "" {
					cmd := exec.Command(cmdName, scriptPath)
					cmd.Env = os.Environ()
					cmd.Stderr = os.Stderr
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					if err := cmd.Run(); err != nil {
						gprint.PrintError(fmt.Sprintf("Execute Compilation Script Failed: %+v", err))
						return
					}
				}
			}
			var name string = "vcpkg"
			vcpkgBinary := filepath.Join(buildPath, name)
			if ok, _ := utils.PathIsExist(vcpkgBinary); ok {
				os.Rename(vcpkgBinary, filepath.Join(config.VCpkgDir, name))
				that.setEnvForVcpkg()
			}
			os.RemoveAll(filepath.Join(config.VCpkgDir, "buildtrees"))
		} else {
			fPath := filepath.Join(config.VCpkgDir, "vcpkg.exe")
			that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(that.Conf.Cpp.WinVCpkgToolUrls[runtime.GOARCH])
			if that.fetcher.Url != "" {
				if size := that.fetcher.GetAndSaveFile(fPath); size == 0 {
					os.RemoveAll(fPath)
				} else {
					that.setEnvForVcpkg()
				}
			}
		}
	}
}

func (that *CppManager) setEnvForVcpkg() {
	if runtime.GOOS == utils.Windows {
		that.env.SetEnvForWin(map[string]string{
			"PATH": config.VCpkgDir,
		})
	} else {
		vcpkgEnv := fmt.Sprintf(utils.VcpkgEnv, config.VCpkgDir)
		that.env.UpdateSub(utils.SUB_VCPKG, vcpkgEnv)
	}
}
