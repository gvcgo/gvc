package vctrl

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gvcgo/goutils/pkgs/archiver"
	"github.com/gvcgo/goutils/pkgs/crypt"
	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/gtea/input"
	"github.com/gvcgo/goutils/pkgs/koanfer"
	"github.com/gvcgo/goutils/pkgs/request"
	"github.com/gvcgo/goutils/pkgs/storage"
	config "github.com/gvcgo/gvc/pkgs/confs"
	"github.com/gvcgo/gvc/pkgs/utils"
)

/*
Use github/gitee repositary as remote storage for your local configurations.
*/

type RepoType string

type EncryptoType string

const (
	RepoTypeGithub  RepoType     = "github"
	RepoTypeGitee   RepoType     = "gitee"
	RepoName        string       = "gvc_configs"
	StorageConfName string       = ".remote_storage.json"
	EncryptByAES    EncryptoType = "aes"
	EncryptByZip    EncryptoType = "zip"
	EncryptByNone   EncryptoType = "none"
)

/*
New configuration for gvc using github/gitee as remote storage.
*/
type StorageConf struct {
	Type        RepoType `json,koanf:"type"`
	UserName    string   `json,koanf:"username"` // username for github or gitee.
	AccessToken string   `json,koanf:"token"`
	CryptoKey   string   `json,koanf:"crypto_key"` // Key to encrypt your private data like passwords and ssh files, etc.
	ProxyURI    string   `json,koanf:"proxy_uri"`
}

type Synchronizer struct {
	CNF      *StorageConf `json,koanf:"storage"`
	repoName string
	storage  storage.IStorage
	path     string
	koanfer  *koanfer.JsonKoanfer
}

func NewSynchronizer(repoName ...string) (s *Synchronizer) {
	s = &Synchronizer{
		CNF: &StorageConf{},
	}
	if len(repoName) == 0 {
		s.repoName = RepoName
	} else {
		s.repoName = repoName[0]
	}
	s.path = filepath.Join(config.GVCDir, StorageConfName)
	s.koanfer, _ = koanfer.NewKoanfer(s.path)
	s.initiate()
	return
}

func (that *Synchronizer) GetConfPath() string {
	return that.path
}

func (that *Synchronizer) initiate() {
	if that.koanfer == nil {
		gprint.PrintError("nil koanfer.")
		return
	}

	that.koanfer.Load(that.CNF)

	// configs for remote repo.
	if that.CNF.AccessToken == "" || that.CNF.UserName == "" || that.CNF.CryptoKey == "" {
		that.Setup()
	}

	// remote storage.
	switch that.CNF.Type {
	case RepoTypeGithub:
		gh := storage.NewGhStorage(that.CNF.UserName, that.CNF.AccessToken)
		gh.Proxy = that.CNF.ProxyURI
		that.storage = gh
	case RepoTypeGitee:
		that.storage = storage.NewGtStorage(that.CNF.UserName, that.CNF.AccessToken)
	default:
		gprint.PrintError("unsupported repo type.")
		os.Exit(1)
	}
	that.createRepo()
}

func (that *Synchronizer) Setup() {
	if ok, _ := utils.PathIsExist(that.path); ok {
		that.koanfer.Load(that.CNF)
	}
	var (
		repoType  string = "RepoType"
		userName  string = "UserName"
		authToken string = "AuthToken"
		cryptoKey string = "CryptoKey"
		proxyUri  string = "ProxyURI"
	)
	mInput := input.NewMultiInput()
	repoTypeValueList := []string{
		string(RepoTypeGithub),
		string(RepoTypeGitee),
	}

	mInput.AddOneOption(
		repoType,
		repoTypeValueList,
		input.MWithWidth(60),
		input.MWithPlaceholder("choose your remote repository type."),
		input.MWithDefaultValue(string(that.CNF.Type)),
	)

	mInput.AddOneItem(
		userName,
		input.MWithWidth(60),
		input.MWithPlaceholder("your github/gitee username."),
		input.MWithDefaultValue(that.CNF.UserName),
	)

	mInput.AddOneItem(
		authToken,
		input.MWithWidth(60),
		input.MWithPlaceholder("your github/gitee access token."),
		input.MWithDefaultValue(that.CNF.AccessToken),
	)

	mInput.AddOneItem(
		cryptoKey,
		input.MWithWidth(60),
		input.MWithPlaceholder("your crypto key to encrypt private data."),
		input.MWithEchoChar("*"),
		input.MWithEchoMode(textinput.EchoPassword),
		input.MWithDefaultValue(that.CNF.CryptoKey),
	)

	mInput.AddOneItem(
		proxyUri,
		input.MWithWidth(60),
		input.MWithPlaceholder("proxy uri for github."),
		input.MWithDefaultValue(that.CNF.ProxyURI),
	)

	mInput.Run()
	result := mInput.Values()
	if r := result[repoType]; r != "" {
		that.CNF.Type = RepoType(r)
	}
	if r := result[userName]; r != "" {
		that.CNF.UserName = r
	}
	if r := result[authToken]; r != "" {
		that.CNF.AccessToken = r
	}
	if r := result[cryptoKey]; r != "" {
		that.CNF.CryptoKey = r
	}
	if r := result[proxyUri]; r != "" {
		that.CNF.ProxyURI = r
	}
	that.koanfer.Save(that.CNF)
}

// Creates remote repo if needed.
func (that *Synchronizer) createRepo() {
	if that.storage == nil {
		gprint.PrintError("No remote storages found.")
		return
	}
	r := that.storage.GetRepoInfo(that.repoName)
	j := gjson.New(r)
	if j.Get("id").Int64() == 0 {
		gprint.PrintInfo("Create remote repo: %s .", that.CNF.UserName+"/"+that.repoName)
		that.storage.CreateRepo(that.repoName)
	}
}

func (that *Synchronizer) formatKeyForAES() (newKey string) {
	if that.CNF.CryptoKey == "" {
		return
	}
	data := []byte(that.CNF.CryptoKey)
	newKey = fmt.Sprintf("%x", md5.Sum(data))[:16]
	return
}

func (that *Synchronizer) upload(fPath, remoteFileName string) (r []byte) {
	if that.storage == nil {
		gprint.PrintError("No remote storages found.")
		return
	}
	content := that.storage.GetContents(that.repoName, "", remoteFileName)
	shaStr := gjson.New(content).Get("sha").String()
	return that.storage.UploadFile(that.repoName, "", fPath, shaStr)
}

func (that *Synchronizer) UploadFile(fPath, remoteFileName string, et EncryptoType) {
	if ok, _ := utils.PathIsExist(fPath); !ok {
		gprint.PrintError("File not exist: %s", fPath)
		return
	}
	switch et {
	case EncryptByAES:
		if that.CNF.CryptoKey == "" {
			gprint.PrintError("No crypto key found.")
			return
		}
		cc := crypt.NewCrptWithKey([]byte(that.formatKeyForAES()))
		content, err := os.ReadFile(fPath)
		if err != nil {
			gprint.PrintError("Read file error: %+v", err)
			return
		}
		if r, err := cc.AesEncrypt([]byte(content)); err != nil {
			gprint.PrintError("Encrypt file error: %+v", err)
			return
		} else {
			fPath = filepath.Join(config.GVCBackupDir, remoteFileName)
			if err = os.WriteFile(fPath, r, os.ModePerm); err != nil {
				gprint.PrintError("Write file error: %+v", err)
				return
			}
		}
	case EncryptByZip:
		if that.CNF.CryptoKey == "" {
			gprint.PrintError("No crypto key found.")
			return
		}
		if archive, err := archiver.NewArchiver(fPath, config.GVCBackupDir, false); err == nil {
			archive.SetZipName(remoteFileName)
			archive.SetPassword(that.CNF.CryptoKey)
			err = archive.ZipDir()
			if err != nil {
				gprint.PrintError("Zip dir error: %+v", err)
				return
			}
		}
		fPath = filepath.Join(config.GVCBackupDir, remoteFileName)
	default:
		// copy local file to backupdir then upload.
		content, err := os.ReadFile(fPath)
		if err != nil {
			gprint.PrintError("Read file error: %+v", err)
			return
		}
		fPath = filepath.Join(config.GVCBackupDir, remoteFileName)
		if err = os.WriteFile(fPath, content, os.ModePerm); err != nil {
			gprint.PrintError("Write file error: %+v", err)
			return
		}
	}

	r := that.upload(fPath, remoteFileName)
	j := gjson.New(r)
	if j.Get("content.path").String() != "" && j.Get("content.sha").String() != "" {
		gprint.PrintSuccess("uploaded successfully: %s", fPath)
	} else {
		gprint.PrintWarning("error occurred: %s", string(r))
	}
}

func (that *Synchronizer) download(remoteFileName string) (dUrl string) {
	if that.storage == nil {
		gprint.PrintError("No remote storages found.")
		return
	}
	content := that.storage.GetContents(that.repoName, "", remoteFileName)
	dUrl = gjson.New(content).Get("download_url").String()
	if dUrl == "" {
		gprint.PrintWarning("can not find %s in remote repo. %s", remoteFileName, string(content))
	}
	return
}

func (that *Synchronizer) DownloadFile(fPath, remoteFileName string, et EncryptoType) {
	dUrl := that.download(remoteFileName)
	if dUrl == "" {
		return
	}
	// download and deploy files.
	fetcher := request.NewFetcher()
	fetcher.Timeout = time.Minute * 30
	fetcher.SetUrl(dUrl)
	fetcher.Proxy = that.CNF.ProxyURI

	srcPath := filepath.Join(config.GVCBackupDir, remoteFileName)
	if size := fetcher.GetAndSaveFile(srcPath, true); size > 20 {
		switch et {
		case EncryptByAES:
			if that.CNF.CryptoKey == "" {
				gprint.PrintError("No crypto key found.")
				return
			}

			cc := crypt.NewCrptWithKey([]byte(that.formatKeyForAES()))
			content, err := os.ReadFile(srcPath)
			if err != nil {
				gprint.PrintError("Read file failed: %+v", err)
				return
			}
			if r, err := cc.AesDecrypt([]byte(content)); err != nil {
				gprint.PrintError("Decrypt file failed: %+v", err)
				return
			} else {
				// deploy remote file to local.
				if err = os.WriteFile(fPath, r, os.ModePerm); err != nil {
					gprint.PrintError("Write file failed: %+v", err)
					return
				}
			}
		case EncryptByZip:
			if archive, err := archiver.NewArchiver(srcPath, fPath, false); err == nil {
				archive.SetPassword(that.CNF.CryptoKey)
				_, err = archive.UnArchive()
				if err != nil {
					gprint.PrintError("unarchive failed: %+v", err)
					return
				}
				gprint.PrintSuccess("download successed: %s", fPath)
			}
			extraFile := filepath.Join(fPath, filepath.Base(fPath))
			if ok, _ := utils.PathIsExist(extraFile); ok {
				os.RemoveAll(extraFile)
			}
		default:
			content, err := os.ReadFile(srcPath)
			if err != nil {
				gprint.PrintError("Read file failed: %+v", err)
				return
			}
			// deploy remote file to local.
			if err = os.WriteFile(fPath, content, os.ModePerm); err != nil {
				gprint.PrintError("Write file failed: %+v", err)
				return
			}
		}
	} else {
		gprint.PrintError("download failed: %s", remoteFileName)
	}
}
