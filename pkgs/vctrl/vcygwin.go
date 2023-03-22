package vctrl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
)

type CygwinPackages struct {
	Packages []string `koanf:"packages"`
}

type Cygwin struct {
	*downloader.Downloader
	Conf            *config.GVConfig
	env             *utils.EnvsHandler
	k               *koanf.Koanf
	parser          *yaml.YAML
	P               *CygwinPackages
	installerPath   string
	packageFilePath string
}

func NewCygwin() (cy *Cygwin) {
	cy = &Cygwin{
		Downloader:      &downloader.Downloader{},
		Conf:            config.New(),
		env:             utils.NewEnvsHandler(),
		k:               koanf.New("."),
		parser:          yaml.Parser(),
		P:               &CygwinPackages{},
		installerPath:   filepath.Join(config.CygwinFilesDir, config.CygwinInstallerName),
		packageFilePath: filepath.Join(config.CygwinFilesDir, config.CygwinPackageFileName),
	}
	cy.initDir()
	return
}

func (that *Cygwin) initDir() {
	if ok, _ := utils.PathIsExist(config.CygwinRootDir); !ok {
		if err := os.MkdirAll(config.CygwinRootDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
}

func (that *Cygwin) getInstaller() {
	if ok, _ := utils.PathIsExist(that.installerPath); !ok {
		that.Url = that.Conf.Cygwin.InstallerUrl
		if that.Url != "" {
			fmt.Println("[Download] ", that.Url)
			that.Timeout = 10 * time.Minute
			if size := that.GetFile(that.installerPath, os.O_CREATE|os.O_WRONLY, 0777); size == 0 {
				fmt.Println("[Download Cygwin installer failed!]")
				os.RemoveAll(that.installerPath)
			} else {
				if runtime.GOOS == utils.Windows {
					that.env.SetEnvForWin(map[string]string{
						"PATH": fmt.Sprintf("%s;%s;%s", config.CygwinFilesDir, config.CygwinRootDir, config.CygwinBinaryDir),
					})
				}
			}
		}
	}

	if ok, _ := utils.PathIsExist(that.packageFilePath); !ok {
		that.Url = that.Conf.Cygwin.PackagesUrl
		if that.Url == "" {
			fmt.Println("[Download] ", that.Url)
			that.Timeout = 10 * time.Minute
			if size := that.GetFile(that.packageFilePath, os.O_CREATE|os.O_WRONLY, 0644); size == 0 {
				fmt.Println("[Download cygwin package file failed!]")
				os.RemoveAll(that.packageFilePath)
			}
		}
	}
}

// func (that *Cygwin) getPackagesInfo() {
// 	fpath := filepath.Join(config.CygwinFilesDir, config.CygwinPackageFileName)
// 	if ok, _ := utils.PathIsExist(fpath); !ok {
// 		err := that.k.Load(file.Provider(fpath), that.parser)
// 		if err != nil {
// 			fmt.Println("[Config Load Failed] ", err)
// 			return
// 		}
// 		that.k.UnmarshalWithConf("", that.P, koanf.UnmarshalConf{Tag: "koanf"})
// 	} else {
// 		pList := []string{"git", "gcc", "gdb", "gdbm", "clang", "bash", "bashdb", "wget"}
// 		that.P.Packages = append(that.P.Packages, pList...)
// 	}
// }

func (that *Cygwin) InstallByDefault(packInfo string) {
	that.getInstaller()
	if packInfo == "" {
		// that.getPackagesInfo()
		// packInfo = strings.Join(that.P.Packages, ",")
		packInfo = "git,bash,wget,gcc,gdb,clang,openssh,bashdb,gdbm,gcc-fortran,clang-analyzer,clang-doc,bash-completion,bash-devel,bash-completion-cmake"
	}

	fmt.Println("[Install Packages] ", packInfo)

	if ok, _ := utils.PathIsExist(that.installerPath); ok && runtime.GOOS == utils.Windows {
		ePath := os.Getenv("PATH")
		if !strings.Contains(ePath, config.CygwinFilesDir) {
			ePath = fmt.Sprintf("%s;%s", config.CygwinFilesDir, ePath)
			os.Setenv("PATH", ePath)
		}
		cmd := exec.Command(config.CygwinInstallerName, "-q", "-f", "-N", "-O", "-s",
			that.Conf.Cygwin.MirrorUrls[0], "-R", config.CygwinRootDir,
			"-P", packInfo)
		cmd.Env = os.Environ()
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
