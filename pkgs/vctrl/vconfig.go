package vctrl

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gogf/gf/os/genv"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/studio-b12/gowebdav"
)

type WebdavConf struct {
	Host            string `koanf:"url"`
	Username        string `koanf:"username"`
	Password        string `koanf:"password"`
	RemoteDir       string `koanf:"remote_dir"`
	LocalDir        string `koanf:"local_dir"`
	DefaultFilesUrl string `koanf:"default_files"`
}

type VSCodeExtIds struct {
	VSCodeExts []string `koanf:"vscode_exts"`
}

type GVCWebdav struct {
	DavConf    *WebdavConf
	conf       *config.GVConfig
	vscodeExts *VSCodeExtIds
	k          *koanf.Koanf
	parser     *yaml.YAML
	client     *gowebdav.Client
}

func NewGVCWebdav() (gw *GVCWebdav) {
	gw = &GVCWebdav{
		DavConf: &WebdavConf{
			LocalDir: config.GVCBackupDir,
		},
		conf:       config.New(),
		vscodeExts: &VSCodeExtIds{},
		k:          koanf.New("."),
		parser:     yaml.Parser(),
	}
	gw.initeDirs()
	return
}

func (that *GVCWebdav) initeDirs() {
	if ok, _ := utils.PathIsExist(config.GVCBackupDir); config.GVCBackupDir != "" && !ok {
		if err := os.MkdirAll(config.GVCBackupDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", config.GVCBackupDir, err)
		}
	}
	that.Reload()
}

func (that *GVCWebdav) resetKoanf() {
	that.k = koanf.New(".")
	that.parser = yaml.Parser()
}

// save webdav configurations to yml file.
func (that *GVCWebdav) Restore() {
	that.k.Load(structs.Provider(that.DavConf, "koanf"), nil)
	if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
		os.WriteFile(config.GVCWebdavConfigPath, b, 0666)
	}
}

func (that *GVCWebdav) RestoreDefaultGVConf() {
	that.conf.SetDefault()
	that.conf.Restore()
}

func (that *GVCWebdav) ShowConfigPath() {
	fmt.Println("GVC Config File Path: ", config.GVConfigPath)
	fmt.Println("Webdav Config File Path: ", config.GVCWebdavConfigPath)
}

func (that *GVCWebdav) Reload() {
	if ok, _ := utils.PathIsExist(config.GVCWebdavConfigPath); !ok {
		fmt.Println("[Warning] It seems that you have not set up your webdav account.")
		return
	}
	err := that.k.Load(file.Provider(config.GVCWebdavConfigPath), that.parser)
	if err != nil {
		fmt.Println("[Config Load Failed] ", err)
		return
	}
	that.k.UnmarshalWithConf("", that.DavConf, koanf.UnmarshalConf{Tag: "koanf"})
	that.getClient(true)
}

func (that *GVCWebdav) SetAccount() {
	var (
		wUrl string
		name string
		pass string
	)
	fmt.Println("Please enter your webdav host uri,\n[https://dav.jianguoyun.com/dav/]by default: ")
	fmt.Println("How to get your webdav? Please see https://github.com/moqsien/easynotes/blob/main/usage.md.")
	fmt.Scanln(&wUrl)
	fmt.Println("Please enter your webdav username: ")
	fmt.Scanln(&name)
	fmt.Println("Please enter your webdav password: ")
	fmt.Scanln(&pass)
	wUrl = strings.Trim(wUrl, " ")
	name = strings.Trim(name, " ")
	pass = strings.Trim(pass, " ")
	if utils.VerifyUrls(wUrl) {
		that.DavConf.Host = wUrl
	} else if wUrl == "" {
		defaultUrl := "https://dav.jianguoyun.com/dav/"
		if that.conf.Webdav.DefaultWebdavHost != "" {
			defaultUrl = that.conf.Webdav.DefaultWebdavHost
		}
		fmt.Println("Use default webdav url: ", defaultUrl)
		that.DavConf.Host = defaultUrl
	}
	if name != "" {
		that.DavConf.Username = name
	} else {
		fmt.Println("[Warning] Your username is empty!")
	}
	if pass != "" {
		that.DavConf.Password = pass
	} else {
		fmt.Println("[Warning] Your password is empty!")
	}
	if that.conf.Webdav.DefaultWebdavRemoteDir != "" {
		that.DavConf.RemoteDir = that.conf.Webdav.DefaultWebdavRemoteDir
	}
	if that.conf.Webdav.DefaultWebdavLocalDir != "" {
		that.DavConf.RemoteDir = that.conf.Webdav.DefaultWebdavLocalDir
	}
	that.Restore()
}

func (that *GVCWebdav) getClient(force bool) {
	if that.DavConf.Host == "" || that.DavConf.Username == "" || that.DavConf.Password == "" {
		fmt.Println("[Webdav account info missing]")
		fmt.Println("Please set your webdav account info.")
		fmt.Println("Do you want to set your webdav account now? [Y/N]")
		var flag string
		fmt.Scan(&flag)
		if strings.HasPrefix(strings.ToLower(flag), "y") {
			that.SetAccount()
			that.getClient(force)
		}
		return
	}
	if that.client == nil || force {
		that.client = gowebdav.NewClient(that.DavConf.Host,
			that.DavConf.Username, that.DavConf.Password)
		if err := that.client.Connect(); err != nil {
			that.client = nil
			fmt.Println("[Webdav connecting failed] ", err)
			fmt.Println("Please check your webdav account info or network.")
		}
	}
}

func (that *GVCWebdav) Pull() {
	that.getClient(true)
	if that.client == nil {
		return
	}
	iList, err := that.client.ReadDir(that.DavConf.RemoteDir)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			if err := that.client.MkdirAll(that.DavConf.RemoteDir, os.ModePerm); err != nil {
				fmt.Println("Create a new dir for webdav failed! ", err)
				return
			}
		} else {
			fmt.Println("[Get files from webdav failed] ", err)
			return
		}
	}
	if len(iList) > 0 {
		for _, info := range iList {
			if !info.IsDir() {
				rPath := utils.JoinUnixFilePath(that.DavConf.RemoteDir, info.Name())
				b, _ := that.client.Read(rPath)
				if len(b) == 0 {
					r, _ := that.client.ReadStream(rPath)
					file, _ := os.OpenFile(filepath.Join(that.DavConf.LocalDir, info.Name()), os.O_CREATE|os.O_WRONLY, 0666)
					io.Copy(file, r)
					return
				}
				os.WriteFile(filepath.Join(that.DavConf.LocalDir, info.Name()), b, os.ModePerm)
			}
		}
	}
}

func (that *GVCWebdav) Push() {
	that.getClient(true)
	if that.client == nil {
		return
	}
	_, err := that.client.ReadDir(that.DavConf.RemoteDir)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			if err := that.client.MkdirAll(that.DavConf.RemoteDir, os.ModePerm); err != nil {
				fmt.Println("Create a new dir for webdav failed! ", err)
				return
			}
		}
		fmt.Println(err)
		return
	}
	if iList, _ := os.ReadDir(that.DavConf.LocalDir); len(iList) > 0 {
		for _, info := range iList {
			if !info.IsDir() {
				b, _ := os.ReadFile(filepath.Join(that.DavConf.LocalDir, info.Name()))
				rPath := utils.JoinUnixFilePath(that.DavConf.RemoteDir, info.Name())
				that.client.Write(rPath, b, os.ModePerm)
			}
		}
	}
}

func (that *GVCWebdav) getFilesToSync() (fm config.Filemap) {
	if len(that.conf.Webdav.FilesToSync) > 0 {
		fm = that.conf.Webdav.FilesToSync[runtime.GOOS]
		for k, v := range fm {
			if strings.Contains(v, "$home$") {
				fm[k] = filepath.Join(utils.GetHomeDir(), strings.ReplaceAll(v, "$home$", ""))
			} else if strings.Contains(v, "$appdata$") {
				fm[k] = filepath.Join(utils.GetWinAppdataEnv(), strings.ReplaceAll(v, "$appdata$", ""))
			}
		}
	}
	return
}

// https://code.visualstudio.com/docs/getstarted/keybindings
func (that *GVCWebdav) modifyKeybindings(backupPath string) {
	if ok, _ := utils.PathIsExist(backupPath); !ok {
		return
	}
	content, _ := os.ReadFile(backupPath)
	switch runtime.GOOS {
	case utils.MacOS:
		cStr := utils.BatchReplaceAll(string(content), map[string]string{
			"win+":  "cmd+",
			"Win+":  "cmd+",
			"meta+": "cmd+",
			"Meta+": "cmd+",
		})
		os.WriteFile(backupPath, []byte(cStr), 0666)
	case utils.Windows:
		cStr := utils.BatchReplaceAll(string(content), map[string]string{
			"cmd+":  "win+",
			"Cmd+":  "win+",
			"meta+": "win+",
			"Meta+": "win+",
		})
		os.WriteFile(backupPath, []byte(cStr), 0666)
	default:
		cStr := utils.BatchReplaceAll(string(content), map[string]string{
			"cmd+": "meta+",
			"Cmd+": "meta+",
			"win+": "meta+",
			"Win+": "meta+",
		})
		os.WriteFile(backupPath, []byte(cStr), 0666)
	}
}

func (that *GVCWebdav) FetchAndApplySettings() {
	fmt.Println("Fetching config files from webdav...")
	that.Pull()
	for backupName, fpath := range that.getFilesToSync() {
		if fpath == "" {
			continue
		}
		backupPath := filepath.Join(that.DavConf.LocalDir, backupName)
		if ok, _ := utils.PathIsExist(backupPath); ok {
			fmt.Println("[set config files] ", backupPath)
			if content, _ := os.ReadFile(backupPath); len(content) == 0 {
				continue
			}
			if backupName == config.CodeKeybindingsBackupFileName {
				that.modifyKeybindings(backupPath)
			}
			utils.CopyFile(backupPath, fpath)
		}
	}
}

// install vscode extensions
func (that *GVCWebdav) InstallVSCodeExts(backupPath string) {
	that.Pull()
	if backupPath == "" {
		backupPath = filepath.Join(that.DavConf.LocalDir, config.CodeExtensionsBackupFileName)
	}
	if ok, _ := utils.PathIsExist(backupPath); ok {
		err := that.k.Load(file.Provider(backupPath), that.parser)
		if err != nil {
			fmt.Println("[Config Load Failed] ", err)
			return
		}
		that.k.UnmarshalWithConf("", that.vscodeExts, koanf.UnmarshalConf{Tag: "koanf"})
	} else {
		fmt.Println("[Can not find extensions backup file] ", backupPath)
		return
	}
	if len(that.vscodeExts.VSCodeExts) == 0 {
		that.vscodeExts.VSCodeExts = that.conf.Code.ExtIdentifiers
	}
	for _, extName := range that.vscodeExts.VSCodeExts {
		cmd := exec.Command("code", "--install-extension", extName)
		cmd.Env = genv.All()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Run()
	}
}

// gather version extensions info
func (that *GVCWebdav) gatherVSCodeExtsions() {
	cmd := exec.Command("code", "--list-extensions")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %sn", err)
		return
	}
	iNameList := strings.Split(string(out), "\n")
	if len(iNameList) > 0 {
		newList := []string{}
		fmt.Println("Local installed vscode extensions: ")
		for _, iName := range iNameList {
			if strings.Contains(iName, ".") && len(iName) > 3 {
				newList = append(newList, iName)
				fmt.Println(iName)
			}
		}
		if len(newList) > 0 {
			that.vscodeExts.VSCodeExts = newList
			that.resetKoanf()
			that.k.Load(structs.Provider(that.vscodeExts, "koanf"), nil)
			if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
				fpath := filepath.Join(that.DavConf.LocalDir, config.CodeExtensionsBackupFileName)
				os.WriteFile(fpath, b, 0666)
			}
		}
	}
}

func (that *GVCWebdav) GatherAndPushSettings() {
	if that.DavConf.LocalDir == "" {
		that.DavConf.LocalDir = config.GVCBackupDir
	}
	that.gatherVSCodeExtsions()
	if ok, _ := utils.PathIsExist(that.DavConf.LocalDir); !ok {
		if err := os.MkdirAll(that.DavConf.LocalDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", that.DavConf.LocalDir, err)
			return
		}
	}
	for backupName, fpath := range that.getFilesToSync() {
		if ok, _ := utils.PathIsExist(fpath); ok {
			fmt.Println("[gathering file] ", backupName)
			utils.CopyFile(fpath, filepath.Join(that.DavConf.LocalDir, backupName))
		}
	}
	fmt.Println("Pushing config files to webdav...")
	that.Push()
}
