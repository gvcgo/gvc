package vctrl

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	config "github.com/moqsien/gvc/pkgs/confs"
	downloader "github.com/moqsien/gvc/pkgs/fetcher"
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
	*downloader.Downloader
	Conf *config.GVConfig
	Doc  *goquery.Document
	env  *utils.EnvsHandler
}

func NewCppManager() (cm *CppManager) {
	cm = &CppManager{
		Downloader: &downloader.Downloader{},
		Conf:       config.New(),
		env:        utils.NewEnvsHandler(),
	}
	cm.initDirs()
	return
}

func (that *CppManager) initDirs() {
	if ok, _ := utils.PathIsExist(config.CppFilesDir); !ok {
		if err := os.MkdirAll(config.CppFilesDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.Msys2Dir); !ok {
		if err := os.MkdirAll(config.Msys2Dir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.VCpkgDir); !ok {
		if err := os.MkdirAll(config.VCpkgDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.CppDownloadDir); !ok {
		if err := os.MkdirAll(config.CppDownloadDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.CygwinRootDir); !ok {
		if err := os.MkdirAll(config.CygwinRootDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
}

func (that *CppManager) getDoc() {
	that.Url = that.Conf.Cpp.MsysInstallerUrl
	if !utils.VerifyUrls(that.Url) {
		return
	}
	that.Doc, _ = goquery.NewDocumentFromReader(bytes.NewBuffer(that.GetWithColly()))
}

func (that *CppManager) getMsys2Installer() (fPath string) {
	if that.Doc == nil {
		that.getDoc()
	}
	if that.Doc != nil {
		var exeUrl string
		maxIdx := that.Doc.Find("table#list").Find("tr").Last().Index()
		for i := maxIdx; i >= 0; i-- {
			exeUrl = that.Doc.Find("table#list").Find("tr").Eq(i).Find("a").AttrOr("href", "")
			if strings.HasSuffix(exeUrl, ".exe") {
				break
			}
		}

		if exeUrl != "" {
			if !strings.HasPrefix(exeUrl, "http://") {
				exeUrl, _ = url.JoinPath(that.Conf.Cpp.MsysInstallerUrl, exeUrl)
			}
			fPath = filepath.Join(config.CppDownloadDir, Msys2InstallerName)
			that.Url = exeUrl
			that.GetFile(fPath, os.O_CREATE|os.O_WRONLY, 0777)
		}
	}
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
			fmt.Println("Execute Msys2Installer Failed: ", err)
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
			fmt.Println("Execute uninstall.exe Failed: ", err)
			return
		}
	}
}

func (that *CppManager) getCygwinInstaller() (fPath string) {
	fPath = filepath.Join(config.CppDownloadDir, CygwinInstallerName)
	if ok, _ := utils.PathIsExist(fPath); !ok {
		that.Url = that.Conf.Cpp.CygwinInstallerUrl
		if that.Url != "" {
			that.Timeout = 10 * time.Minute
			if size := that.GetFile(fPath, os.O_CREATE|os.O_WRONLY, 0777); size == 0 {
				fmt.Println("[Download Cygwin installer failed!]")
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

	fmt.Println("[Install Packages] ", packInfo)

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
			fmt.Println("Execute CygwinInstaller Failed: ", err)
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

func (that *CppManager) GetVCPkg() {
	that.Url = that.Conf.Cpp.VCpkgUrl
	resp := that.GetUrl()
	fmt.Println(resp.Request.URL)
}
