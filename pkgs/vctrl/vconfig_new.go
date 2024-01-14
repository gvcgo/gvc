package vctrl

import (
	"fmt"
	"path/filepath"

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
	if that.CNF.AccessToken == "" {
		var SType RepoType
		fmt.Println("Choose your repo type: ")
		fmt.Println("1. Github. ")
		fmt.Println("2. Gitee. ")
		fmt.Scanln(&SType)
		that.CNF.Type = SType
		// TODO: config uploader.
	}
	that.storage = storage.NewGhStorage(that.CNF.UserName, that.CNF.AccessToken)
}

func (that *Synchronizer) UploadFile() {

}

func (that *Synchronizer) DownloadFile() {

}
