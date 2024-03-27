package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/gvc/conf"
	"github.com/gvcgo/gvc/pkg/asciinema"
	"github.com/gvcgo/gvc/utils"
)

/*
Upload file/dir to Repo.
*/
func UploadToRepo(repoType RepoType, encryptEnabled bool, remoteFileName, localFilePath string) {
	cfg := conf.NewGVConfig()
	cfg.Load()
	repoName := cfg.GetBackupRepo()
	repo := NewRepo(repoType, encryptEnabled)
	if err := repo.Upload(repoName, remoteFileName, localFilePath); err != nil {
		gprint.PrintError("upload file failed: %+v", err)
	}
}

/*
Download file/dir from Repo.
*/
func DownloadFromRepo(repoType RepoType, encryptEnabled bool, remoteFileName, localFilePath string) {
	cfg := conf.NewGVConfig()
	cfg.Load()
	repoName := cfg.GetBackupRepo()
	repo := NewRepo(repoType, encryptEnabled)

	backupFileName := fmt.Sprintf("%s.old", localFilePath)
	if ok, _ := gutils.PathIsExist(localFilePath); ok {
		fmt.Println(gprint.CyanStr("File or directory already exists: %s", localFilePath))
		fmt.Println(gprint.YellowStr("Backup the old files or not?[y/N]"))
		var okStr string
		fmt.Scanln(&okStr)
		if strings.ToLower(okStr) == "y" {
			os.RemoveAll(backupFileName)
			os.Rename(localFilePath, backupFileName)
		} else {
			if utils.PathIsDir(localFilePath) {
				os.RemoveAll(localFilePath)
			}
		}
	}

	err := repo.Download(repoName, remoteFileName, localFilePath)
	if err != nil {
		gprint.PrintError("download file failed: %+v", err)
		// recover from backuped files.
		if ok, _ := gutils.PathIsExist(backupFileName); ok {
			os.Rename(backupFileName, localFilePath)
		}
	}
}

/*
VSCode user-settings.json keymap.json extension_list.json
*/
func getVSCodeSettingsFile() (fPath string) {
	switch runtime.GOOS {
	case gutils.Windows:
		appDataDir, _ := os.UserConfigDir()
		fPath = filepath.Join(appDataDir, "Code", "User", "settings.json")
	case gutils.Linux:
		homeDir, _ := os.UserHomeDir()
		fPath = filepath.Join(homeDir, ".config", "Code", "User", "settings.json")
	case gutils.Darwin:
		homeDir, _ := os.UserHomeDir()
		fPath = filepath.Join(homeDir, "Library", "Application Support", "Code", "User", "settings.json")
	default:
	}
	return
}

func getVSCodeKeybindingFile() (fPath string) {
	switch runtime.GOOS {
	case gutils.Windows:
		appDataDir, _ := os.UserConfigDir()
		fPath = filepath.Join(appDataDir, "Code", "User", "keybindings.json")
	case gutils.Linux:
		homeDir, _ := os.UserHomeDir()
		fPath = filepath.Join(homeDir, ".config", "Code", "User", "keybindings.json")
	case gutils.Darwin:
		homeDir, _ := os.UserHomeDir()
		fPath = filepath.Join(homeDir, "Library", `Application Support`, "Code", "User", "keybindings.json")
	default:
	}
	return
}

func getVSCodeExtensionsFile() (fPath string) {
	vscodeDir := filepath.Join(conf.GetGVCWorkDir(), "vscode_data")
	os.MkdirAll(vscodeDir, os.ModePerm)
	return filepath.Join(vscodeDir, "vscode_extensions.txt")
}

func isCodeInstalled() bool {
	_, err := gutils.ExecuteSysCommand(true, "", "code", "--help")
	return err == nil
}

func getCodeBinPath() (r string) {
	switch runtime.GOOS {
	case gutils.Windows:
		path1 := `C:\Program Files\Microsoft VS Code\Code.exe`
		if ok, _ := gutils.PathIsExist(path1); ok {
			return path1
		}
		appDataDir, _ := os.UserConfigDir()
		path2 := filepath.Join(appDataDir, "AppData", "Local", "Programs", "Code.exe")
		if ok, _ := gutils.PathIsExist(path2); ok {
			return path2
		}
	case gutils.Darwin:
		path1 := `/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code`
		if ok, _ := gutils.PathIsExist(path1); ok {
			return path1
		}
	case gutils.Linux:
		path1 := `/usr/share/code`
		if ok, _ := gutils.PathIsExist(path1); ok {
			return path1
		}
		path2 := `/usr/bin/code`
		if ok, _ := gutils.PathIsExist(path2); ok {
			return path2
		}
	default:
	}
	return
}

func collectVSCodeExtensions() {
	cmdStr := "code"
	if !isCodeInstalled() {
		cmdStr = getCodeBinPath()
	}
	if cmdStr != "" {
		b, err := gutils.ExecuteSysCommand(true, "", cmdStr, "--list-extensions")
		if err == nil {
			content := b.String()
			if len(content) > 0 {
				os.WriteFile(getVSCodeExtensionsFile(), []byte(content), os.ModePerm)
			}
		}
	}
}

func installVSCodeExtensions() {
	content, _ := os.ReadFile(getVSCodeExtensionsFile())
	extensions := strings.Split(string(content), "\n")
	cmdStr := "code"
	if !isCodeInstalled() {
		cmdStr = getCodeBinPath()
	}
	if cmdStr != "" {
		for _, extId := range extensions {
			gutils.ExecuteSysCommand(false, "", cmdStr, "--install-extension", extId)
		}
	}
}

func UploadVSCodeFiles(repoType RepoType) {
	// Upload settings.json
	settingsF := getVSCodeSettingsFile()
	remoteFileName := filepath.Base(settingsF)
	UploadToRepo(repoType, false, remoteFileName, settingsF)

	// Upload keybindings.json
	keybindingsF := getVSCodeKeybindingFile()
	remoteFileName = fmt.Sprintf("%s_%s", runtime.GOOS, filepath.Base(keybindingsF))
	UploadToRepo(repoType, false, remoteFileName, keybindingsF)

	// Upload extensions
	collectVSCodeExtensions()
	remoteFileName = filepath.Base(getVSCodeExtensionsFile())
	UploadToRepo(repoType, false, remoteFileName, getVSCodeExtensionsFile())
}

func DownloadVSCodeFiles(repoType RepoType) {
	// Download settings.json
	settingsF := getVSCodeSettingsFile()
	remoteFileName := fmt.Sprintf("%s_%s", runtime.GOOS, filepath.Base(settingsF))
	DownloadFromRepo(repoType, false, remoteFileName, settingsF)

	// Download keybindings.json
	keybindingsF := getVSCodeKeybindingFile()
	remoteFileName = fmt.Sprintf("%s_%s", runtime.GOOS, filepath.Base(keybindingsF))
	DownloadFromRepo(repoType, false, remoteFileName, keybindingsF)

	// Download extensions
	remoteFileName = filepath.Base(getVSCodeExtensionsFile())
	DownloadFromRepo(repoType, false, remoteFileName, getVSCodeExtensionsFile())
	installVSCodeExtensions()
}

/*
.ssh dir
*/
const (
	dotSSHRemoteFileName string = "dotssh.zip"
)

func getDotSSHDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".ssh")
}

func UploadSSHFiles(repoType RepoType) {
	UploadToRepo(repoType, true, dotSSHRemoteFileName, getDotSSHDir())
}

func DownloadSSHFiles(repoType RepoType) {
	DownloadFromRepo(repoType, true, dotSSHRemoteFileName, getDotSSHDir())
	if runtime.GOOS != gutils.Windows {
		gutils.ExecuteSysCommand(false, "", "chmod", "-R", "700", getDotSSHDir())
	}
}

/*
asciinema id
*/
const (
	AsciinemaIDFileName string = "aciinema.conf"
)

func getAsciinemaIDFile() string {
	return filepath.Join(asciinema.GetAsciinemaWorkDir(), AsciinemaIDFileName)
}

func UploadAsciinemaID(repoType RepoType) {
	UploadToRepo(repoType, true, AsciinemaIDFileName, getAsciinemaIDFile())
}

func DownloadAsciinemaID(repoType RepoType) {
	DownloadFromRepo(repoType, true, AsciinemaIDFileName, getAsciinemaIDFile())
}

/*
Neobox config dir
*/
const (
	NeoboxRemoteFileName string = "neobox.zip"
)

func getNeoboxConfigDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".neobox")
}

func UploadNeoboxConfig(repoType RepoType) {
	UploadToRepo(repoType, true, NeoboxRemoteFileName, getNeoboxConfigDir())
}

func DownloadNeoboxConfig(repoType RepoType) {
	DownloadFromRepo(repoType, true, NeoboxRemoteFileName, getNeoboxConfigDir())
}
