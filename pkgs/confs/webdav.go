package confs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/mholt/archiver/v3"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/studio-b12/gowebdav"
)

func init() {
	if ok, _ := utils.PathIsExist(GVCWorkDir); !ok {
		if err := os.MkdirAll(GVCWorkDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", GVCWorkDir, err)
		}
	}
}

type WebdavConf struct {
	Host            string `koanf:"url"`
	Username        string `koanf:"username"`
	Password        string `koanf:"password"`
	RemoteDir       string `koanf:"remote_dir"`
	RocalDir        string `koanf:"local_dir"`
	DefaultFilesUrl string `koanf:"default_files"`
	k               *koanf.Koanf
	parser          *yaml.YAML
	client          *gowebdav.Client
	d               *downloader.Downloader
}

func NewWebdav() (r *WebdavConf) {
	r = &WebdavConf{
		RemoteDir: "/gvc_backups",
		k:         koanf.New("."),
		parser:    yaml.Parser(),
		d:         &downloader.Downloader{},
	}
	r.setupW()
	return
}

func (that *WebdavConf) setupW() {
	if ok, _ := utils.PathIsExist(GVCWebdavConfigPath); !ok {
		that.Reset()
	} else {
		that.Reload()
	}
	if ok, _ := utils.PathIsExist(GVCBackupDir); !ok {
		if err := os.MkdirAll(GVCBackupDir, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", GVCBackupDir, err)
		}
	}
}

func (that *WebdavConf) set() {
	that.k.Load(structs.Provider(that, "koanf"), nil)
	if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
		os.WriteFile(GVCWebdavConfigPath, b, 0666)
	}
}

func (that *WebdavConf) Reset() {
	that.RocalDir = GVCBackupDir
	that.Host = "https://dav.jianguoyun.com/dav/"
	that.DefaultFilesUrl = "https://gitee.com/moqsien/gvc/releases/download/v1/misc-all.zip"
	that.set()
}

func (that *WebdavConf) Reload() {
	err := that.k.Load(file.Provider(GVCWebdavConfigPath), that.parser)
	if err != nil {
		fmt.Println("[Config Load Failed] ", err)
		return
	}
	that.k.UnmarshalWithConf("", that, koanf.UnmarshalConf{Tag: "koanf"})
	if that.Password != "" && that.Username != "" {
		that.client = gowebdav.NewClient(that.Host, that.Username, that.Password)
		if err := that.client.Connect(); err != nil {
			that.client = nil
			fmt.Println("[Webdav Connect Failed] ", err)
		}
	}
}

func (that *WebdavConf) ShowDavConfigPath() {
	fmt.Println(GVCWebdavConfigPath)
}

func (that *WebdavConf) SetConf() {
	var (
		wUrl string
		name string
		pass string
	)
	fmt.Println("Please enter your webdav host uri: ")
	fmt.Scanln(&wUrl)
	fmt.Println("Please enter your webdav username: ")
	fmt.Scanln(&name)
	fmt.Println("Please enter your webdav password: ")
	fmt.Scanln(&pass)
	wUrl = strings.Trim(wUrl, " ")
	name = strings.Trim(name, " ")
	pass = strings.Trim(pass, " ")
	if utils.VerifyUrls(wUrl) {
		that.Host = wUrl
	}
	if name != "" {
		that.Username = name
	}
	if pass != "" {
		that.Password = pass
	}
	that.set()
}

func (that *WebdavConf) GetDefaultFiles() {
	that.d.Url = that.DefaultFilesUrl
	that.d.Timeout = 60 * time.Second
	fpath := filepath.Join(GVCWorkDir, "all.zip")
	if size := that.d.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size > 0 {
		if l, _ := os.ReadDir(that.RocalDir); len(l) == 0 {
			if err := archiver.Unarchive(fpath, that.RocalDir); err != nil {
				fmt.Println("[Unarchive file failed] ", err)
			}
		} else {
			fmt.Println("[Local dir is not empty]")
		}
	}
	os.RemoveAll(fpath)
}

func (that *WebdavConf) Pull() {
	if that.client != nil {
		iList, err := that.client.ReadDir(that.RemoteDir)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				if err := that.client.MkdirAll(that.RemoteDir, 0644); err != nil {
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
					b, _ := that.client.Read(filepath.Join(that.RemoteDir, info.Name()))
					os.WriteFile(filepath.Join(that.RocalDir, info.Name()), b, 0644)
				}
			}
		} else if that.DefaultFilesUrl != "" {
			that.GetDefaultFiles()
		}
	} else {
		fmt.Println("Please set your correct webdav info.")
	}
}

func (that *WebdavConf) Push() {
	if that.client != nil {
		_, err := that.client.ReadDir(that.RemoteDir)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				if err := that.client.MkdirAll(that.RemoteDir, 0644); err != nil {
					fmt.Println("Create a new dir for webdav failed! ", err)
					return
				}
			}
			fmt.Println(err)
			return
		}
		if iList, _ := os.ReadDir(that.RocalDir); len(iList) > 0 {
			for _, info := range iList {
				if !info.IsDir() {
					b, _ := os.ReadFile(filepath.Join(that.RocalDir, info.Name()))
					that.client.Write(filepath.Join(that.RemoteDir, info.Name()), b, 0644)
				}
			}
		} else {
			that.GetDefaultFiles()
		}
	} else {
		fmt.Println("Please set your correct webdav info.")
	}
}
