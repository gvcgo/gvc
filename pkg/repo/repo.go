package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gvcgo/goutils/pkgs/archiver"
	"github.com/gvcgo/goutils/pkgs/crypt"
	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/goutils/pkgs/request"
	"github.com/gvcgo/goutils/pkgs/storage"
	"github.com/gvcgo/gvc/conf"
	"github.com/gvcgo/gvc/utils"
)

type RepoType int

const (
	RepoGithub RepoType = iota
	RepoGitee
)

/*
1. Backups local files to github/gitee repo.
2. Encrypts txt files with AES.
3. Zip dirs with password.
*/
type Repo struct {
	Storage        storage.IStorage
	Type           RepoType
	EncryptEnabled bool
	cfg            *conf.GVConfig
	username       string
}

func NewRepo(repoType RepoType, encryptEnabled bool) (r *Repo) {
	r = &Repo{
		Type:           repoType,
		EncryptEnabled: encryptEnabled,
		cfg:            &conf.GVConfig{},
	}
	r.cfg.Load()
	return
}

func (r *Repo) getStorage() {
	switch r.Type {
	case RepoGithub:
		r.username = r.cfg.GetGitUserName()
		gh := storage.NewGhStorage(r.username, r.cfg.GetGitToken())
		gh.Proxy = r.cfg.GetLocalProxy()
		r.Storage = gh
	case RepoGitee:
		r.username = r.cfg.GetGiteeUserName()
		r.Storage = storage.NewGtStorage(r.username, r.cfg.GetGiteeToken())
	default:
		gprint.PrintError("unsupported repository.")
		os.Exit(1)
	}
}

func (r *Repo) doesRepoExist(repoName string) (ok bool) {
	if r.Storage == nil {
		r.getStorage()
	}
	resp := r.Storage.GetRepoInfo(repoName)
	j := gjson.New(resp)
	if j.Get("id").Int64() == 0 {
		return
	}
	return true
}

// Creates remote repo.
func (r *Repo) Create(repoName string) {
	if ok := r.doesRepoExist(repoName); !ok {
		gprint.PrintInfo("Create remote repo: %s .", r.username+"/"+repoName)
		r.Storage.CreateRepo(repoName)
	}
}

// Uploads local file to remote repo.
func (r *Repo) Upload(repoName, remoteFileName, localFilePath string) (err error) {
	if ok, _ := gutils.PathIsExist(localFilePath); !ok {
		gprint.PrintError("file not found: %s", localFilePath)
		return fmt.Errorf("file not found: %s", localFilePath)
	}
	if remoteFileName == "" {
		remoteFileName = filepath.Base(localFilePath)
	}
	var fPath string
	if r.EncryptEnabled || utils.PathIsDir(localFilePath) {
		password := r.cfg.GetPassword()
		if password == "" {
			return fmt.Errorf("password not found")
		}
		if utils.PathIsDir(localFilePath) {
			// zip with password
			if !strings.HasSuffix(remoteFileName, ".zip") {
				remoteFileName += ".zip"
			}
			if archive, err1 := archiver.NewArchiver(localFilePath, conf.GetGVCWorkDir(), false); err1 == nil {
				archive.SetZipName(remoteFileName)
				archive.SetPassword(password)
				err = archive.ZipDir()
				if err != nil {
					return fmt.Errorf("create zip failed: %+v", err)
				}
			} else {
				return fmt.Errorf("create zip failed: %+v", err1)
			}
			fPath = filepath.Join(conf.GetGVCWorkDir(), remoteFileName)
		} else {
			// encrypt content with password
			cc := crypt.NewCrptWithKey([]byte(password))
			var content []byte
			content, err = os.ReadFile(fPath)
			if err != nil {
				return fmt.Errorf("read file failed: %+v", err)
			}

			if content, err = cc.AesEncrypt([]byte(content)); err != nil {
				return fmt.Errorf("encrypt file failed: %+v", err)
			} else {
				fPath = filepath.Join(conf.GetGVCWorkDir(), remoteFileName)
				if err = os.WriteFile(fPath, content, os.ModePerm); err != nil {
					return fmt.Errorf("write file failed: %+v", err)
				}
			}
		}
	} else {
		// no encryption
		fPath = filepath.Join(conf.GetGVCWorkDir(), remoteFileName)
		if err = gutils.CopyAFile(localFilePath, fPath); err != nil {
			return fmt.Errorf("copy file failed: %+v", err)
		}
	}
	r.Create(repoName)
	if r.Storage != nil {
		content := r.Storage.GetContents(repoName, "", remoteFileName)
		shaStr := gjson.New(content).Get("sha").String()
		resp := r.Storage.UploadFile(repoName, "", fPath, shaStr)
		j := gjson.New(resp)
		if j.Get("content.path").String() != "" && j.Get("content.sha").String() != "" {
			err = nil
		} else {
			err = fmt.Errorf("upload file failed: %s", string(resp))
		}
	}
	os.RemoveAll(fPath)
	return
}

// Downloads file from remote repo to local disk.
func (r *Repo) Download(repoName, remoteFileName, localFilePath string) (err error) {
	if ok := r.doesRepoExist(repoName); !ok {
		return fmt.Errorf("can not find remote repo: %s", repoName)
	}
	r.getStorage()
	if r.Storage == nil {
		return fmt.Errorf("cannot get remote repo")
	}
	content := r.Storage.GetContents(repoName, "", remoteFileName)
	dUrl := gjson.New(content).Get("download_url").String()
	if dUrl == "" {
		return fmt.Errorf("cannot find file: %s in %s", remoteFileName, repoName)
	}

	// download and deploy files.
	fetcher := request.NewFetcher()
	fetcher.Timeout = time.Minute * 30
	fetcher.SetUrl(dUrl)
	if r.Type == RepoGithub {
		fetcher.Proxy = r.cfg.GetLocalProxy()
	}

	fPath := filepath.Join(conf.GetGVCWorkDir(), remoteFileName)

	size := fetcher.GetAndSaveFile(fPath, true)
	if size < 20 {
		return fmt.Errorf("download file failed: %s", remoteFileName)
	}

	if r.EncryptEnabled || strings.HasSuffix(dUrl, ".zip") {
		password := r.cfg.GetPassword()
		if password == "" {
			return fmt.Errorf("password not found")
		}
		if utils.PathIsDir(localFilePath) {
			// zip file
			if archive, err1 := archiver.NewArchiver(fPath, localFilePath, false); err1 == nil {
				archive.SetPassword(password)
				_, err = archive.UnArchive()
				if err != nil {
					return fmt.Errorf("unarchive failed: %+v", err)
				}
				gprint.PrintSuccess("download successed: %s", fPath)
			} else {
				return fmt.Errorf("unarchive failed: %+v", err1)
			}
			extraFile := filepath.Join(localFilePath, filepath.Base(fPath))
			if ok, _ := gutils.PathIsExist(extraFile); ok {
				os.RemoveAll(extraFile)
			}
		} else {
			// encrypted file
			cc := crypt.NewCrptWithKey([]byte(password))
			var content []byte
			content, err = os.ReadFile(fPath)
			if err != nil {
				return fmt.Errorf("read file failed: %+v", err)
			}
			if content, err = cc.AesDecrypt([]byte(content)); err != nil {
				return fmt.Errorf("decrypt file failed: %+v", err)
			} else {
				// deploy remote file to local.
				if err = os.WriteFile(localFilePath, content, os.ModePerm); err != nil {
					return fmt.Errorf("write file failed: %+v", err)
				}
			}
		}
	} else {
		// no encryption
		if err = gutils.CopyAFile(fPath, localFilePath); err != nil {
			return fmt.Errorf("copy file failed: %+v", err)
		}
	}
	os.RemoveAll(fPath)
	return
}

// Delete file from remote repo.
func (r *Repo) Delete(repoName, remoteFileName string) (err error) {
	if ok := r.doesRepoExist(repoName); !ok {
		return fmt.Errorf("can not find remote repo: %s", repoName)
	}
	r.getStorage()
	if r.Storage == nil {
		return fmt.Errorf("cannot get remote repo")
	}
	content := r.Storage.GetContents(repoName, "", remoteFileName)
	j := gjson.New(content)
	dUrl := j.Get("download_url").String()
	shaStr := j.Get("sha").String()
	if dUrl == "" {
		return fmt.Errorf("cannot find file: %s in %s", remoteFileName, repoName)
	}
	remotePath := filepath.Dir(remoteFileName)
	rFileName := filepath.Base(remoteFileName)
	r.Storage.DeleteFile(repoName, remotePath, rFileName, shaStr)
	return
}
