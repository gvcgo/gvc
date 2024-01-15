package vctrl

import (
	"fmt"
	"path/filepath"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/koanfer"
	"github.com/moqsien/goutils/pkgs/storage"
	"github.com/moqsien/gvc/pkgs/confs"
)

type RepoType string

const (
	RepoTypeGithub  RepoType = "github"
	RepoTypeGitee   RepoType = "gitee"
	RepoName        string   = "gvc_configs"
	StorageConfName string   = ".remote_storage.json"
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
	CNF     *StorageConf
	storage storage.IStorage
	path    string
	koanfer *koanfer.JsonKoanfer
}

func NewSynchronizer() (s *Synchronizer) {
	s = &Synchronizer{}
	s.path = filepath.Join(confs.GetGVCWorkDir(), StorageConfName)
	s.koanfer, _ = koanfer.NewKoanfer(s.path)
	s.initiate()
	return
}

func (that *Synchronizer) initiate() {
	if that.koanfer == nil {
		gprint.PrintError("nil koanfer.")
		return
	}
	if that.CNF.AccessToken == "" {
		var SType RepoType
		fmt.Println("Choose your repo type: ")
		fmt.Println("1. Github. ")
		fmt.Println("2. Gitee. ")
		fmt.Scanln(&SType)

		var username string
		fmt.Println("Enter your username: ")
		fmt.Scanln(&username)

		var token string
		fmt.Println("Enter your access token: ")
		fmt.Scanln(&token)

		var key string
		fmt.Println("Enter your crypto key: ")
		fmt.Scanln(&key)

		var proxyUri string
		fmt.Println("Enter your proxy uri: ")
		fmt.Scanln(&proxyUri)

		that.CNF.Type = SType
		that.CNF.UserName = username
		that.CNF.AccessToken = token
		that.CNF.CryptoKey = key
		that.CNF.ProxyURI = proxyUri
		that.koanfer.Save(that.CNF)
	}

	that.koanfer.Load(that.CNF)
	that.storage = storage.NewGhStorage(that.CNF.UserName, that.CNF.AccessToken)
}

func (that *Synchronizer) UploadFile() {

}

func (that *Synchronizer) DownloadFile() {

}
