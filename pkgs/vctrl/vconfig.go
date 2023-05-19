package vctrl

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/gogf/gf/os/genv"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/tui"
	xutils "github.com/moqsien/xtray/pkgs/utils"
	"github.com/pterm/pterm"
	"github.com/studio-b12/gowebdav"
)

type WebdavConf struct {
	Host            string `koanf:"url"`
	Username        string `koanf:"username"`
	Password        string `koanf:"password"`
	RemoteDir       string `koanf:"remote_dir"`
	LocalDir        string `koanf:"local_dir"`
	DefaultFilesUrl string `koanf:"default_files"`
	EncryptPass     string `koanf:"encrypt_pass"`
}

const (
	defaultEncryptPass = "^*wgvc$@]}"
)

type VSCodeExtIds struct {
	VSCodeExts []string `koanf:"vscode_exts"`
}

type GVCWebdav struct {
	DavConf    *WebdavConf
	AESCrypt   *utils.AesCrypt
	conf       *config.GVConfig
	vscodeExts *VSCodeExtIds
	client     *gowebdav.Client
	koanfer    *xutils.Koanfer
	k          *koanf.Koanf
	parser     *yaml.YAML
}

func NewGVCWebdav() (gw *GVCWebdav) {
	gw = &GVCWebdav{
		DavConf: &WebdavConf{
			LocalDir: config.GVCBackupDir,
		},
		conf:       config.New(),
		vscodeExts: &VSCodeExtIds{},
		koanfer:    xutils.NewKoanfer(config.GVCWebdavConfigPath),
		k:          koanf.New("."),
		parser:     yaml.Parser(),
	}
	gw.initeDirs()
	gw.initeAES()
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

func (that *GVCWebdav) initeAES() {
	if that.DavConf.EncryptPass != "" {
		that.AESCrypt = utils.NewCrypt(that.DavConf.EncryptPass)
	} else {
		fmt.Println(color.InYellow("[Warning] use default encryption password."))
		that.AESCrypt = utils.NewCrypt(defaultEncryptPass)
	}
}

func (that *GVCWebdav) resetKoanf() {
	that.k = koanf.New(".")
	that.parser = yaml.Parser()
}

// save webdav configurations to json file.
func (that *GVCWebdav) Restore() {
	that.koanfer.Save(that.DavConf)
}

func (that *GVCWebdav) RestoreDefaultGVConf() {
	that.conf.SetDefault()
	that.conf.Restore()
}

func (that *GVCWebdav) Reload() {
	if ok, _ := utils.PathIsExist(config.GVCWebdavConfigPath); !ok {
		fmt.Println("[Warning] It seems that you have not set up your webdav account.")
		return
	}
	that.koanfer.Load(that.DavConf)
	that.getClient(true)
}

func (that *GVCWebdav) SetWebdavAccount() {
	var (
		wHost      string = "Webdav Host"
		wUname     string = "Username"
		wPass      string = "Password"
		wEncrypter string = "Encrypt Password"
	)
	inputItems := []*tui.InputItem{
		{Title: wHost, Default: "https://dav.jianguoyun.com/dav/", Must: true},
		{Title: wUname, Must: true},
		{Title: wPass, Must: true},
		{Title: wEncrypter, Must: true},
	}

	iput := tui.NewInput(inputItems)
	iput.Render()

	for _, item := range inputItems {
		v := item.String()
		switch item.Title {
		case wHost:
			that.DavConf.Host = v
		case wUname:
			that.DavConf.Username = v
		case wPass:
			that.DavConf.Password = v
		case wEncrypter:
			that.DavConf.EncryptPass = v
		default:
			fmt.Println(pterm.Yellow("unknown input"))
		}
	}

	if that.conf.Webdav.DefaultWebdavRemoteDir != "" {
		that.DavConf.RemoteDir = that.conf.Webdav.DefaultWebdavRemoteDir
	}
	if that.conf.Webdav.DefaultWebdavLocalDir != "" {
		that.DavConf.LocalDir = that.conf.Webdav.DefaultWebdavLocalDir
	}
	that.Restore()
}

func (that *GVCWebdav) getClient(force bool) {
	if that.DavConf.Host == "" || that.DavConf.Username == "" || that.DavConf.Password == "" {
		fmt.Println(pterm.Yellow("It seems that you haven't set a webdav account."))
		ok, _ := pterm.DefaultInteractiveConfirm.WithConfirmText("Set your webdav account right now?").Show()
		pterm.Println()
		if ok {
			that.SetWebdavAccount()
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

func (that *GVCWebdav) decryptFile(fpath string) {
	if ok, _ := utils.PathIsExist(fpath); ok {
		b, _ := os.ReadFile(fpath)
		if len(b) > 0 && that.AESCrypt != nil {
			var err error
			if b, err = that.AESCrypt.AesDecrypt(b); err == nil {
				os.WriteFile(fpath, b, 0666)
			}
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
				// decrypt private info.
				if that.AESCrypt != nil && len(b) > 0 && (strings.Contains(info.Name(), "password") || strings.Contains(info.Name(), "credit")) {
					b, _ = that.AESCrypt.AesDecrypt(b)
				}
				fmt.Println(info.Name())
				if len(b) == 0 {
					r, _ := that.client.ReadStream(rPath)
					fpath := filepath.Join(that.DavConf.LocalDir, info.Name())
					file, _ := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0666)
					io.Copy(file, r)
					file.Close()
					// decrypt private info.
					that.decryptFile(fpath)
					continue
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
				// encrypt private info.
				if that.AESCrypt != nil && (strings.Contains(info.Name(), "password") || strings.Contains(info.Name(), "credit")) {
					b, _ = that.AESCrypt.AesEncrypt(b)
				}
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
