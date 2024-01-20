package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/gtea/input"
	"github.com/moqsien/goutils/pkgs/storage"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

/*
Use github or gitee as picture repo.
Especially for markdown.
*/
type PicRepo struct {
	PicRepoName string
	Syncer      *Synchronizer
	confPath    string
}

func NewPicRepo() (pr *PicRepo) {
	pr = &PicRepo{
		confPath: filepath.Join(config.GVCDir, ".pic_repo"),
	}
	pr.initiate()
	return
}

func (that *PicRepo) initiate() {
	if ok, _ := utils.PathIsExist(that.confPath); !ok {
		that.SetPicRepoName()
	} else {
		content, _ := os.ReadFile(that.confPath)
		that.PicRepoName = string(content)
		if that.PicRepoName == "" {
			that.SetPicRepoName()
		}
	}

	if that.PicRepoName != "" {
		that.Syncer = NewSynchronizer(that.PicRepoName)
	}
}

func (that *PicRepo) SetPicRepoName() {
	ipt := input.NewInput(input.WithPrompt("RepoName: "), input.WithPlaceholder("remote repositary name."))
	ipt.Run()

	value := ipt.Value()
	if value != "" {
		that.PicRepoName = value
		os.WriteFile(that.confPath, []byte(value), os.ModePerm)
	} else {
		gprint.PrintError("No repo name specified.")
		os.Exit(1)
	}
}

func (that *PicRepo) UploadPic(fPath string) {
	if that.Syncer == nil {
		gprint.PrintError("Nil syncer.")
		return
	}

	if ok, _ := utils.PathIsExist(fPath); !ok {
		gprint.PrintError("File does not exist: %s.", fPath)
		return
	}

	allowedExt := []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".svg"}

	for _, ext := range allowedExt {
		if strings.HasSuffix(fPath, ext) {
			name := filepath.Base(fPath)
			that.Syncer.UploadFile(
				fPath,
				name,
				EncryptByNone,
			)

			st := that.Syncer.storage
			// https://github.com/moqsien/neobox_resources/raw/main/raw_domains.txt
			// https://cdn.jsdelivr.net/gh/moqsien/neobox_resources@main/domains.txt
			if ghs, ok := st.(*storage.GhStorage); ok {
				rawUrl := fmt.Sprintf(
					"https://github.com/%s/%s/raw/main/%s",
					ghs.UserName,
					that.PicRepoName,
					name,
				)
				gprint.PrintInfo("rawcontent: %s", rawUrl)
				jsdelvrUrl := fmt.Sprintf(
					"https://cdn.jsdelivr.net/gh/%s/%s@main/%s",
					ghs.UserName,
					that.PicRepoName,
					name,
				)
				gprint.PrintInfo("jsdelivr: %s", jsdelvrUrl)
			} else {
				// https://gitee.com/moqsien/gvc/raw/master/homebrew.sh
				if ght, ok := st.(*storage.GtStorage); ok {
					rawUrl := fmt.Sprintf(
						"https://gitee.com/%s/%s/raw/master/%s",
						ght.UserName,
						that.PicRepoName,
						name,
					)
					gprint.PrintInfo("rawcontent: %s", rawUrl)
				}
			}
			return
		}
	}
	gprint.PrintError("unsurpported format.")
}
