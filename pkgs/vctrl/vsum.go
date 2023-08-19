package vctrl

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

/*
CheckSum for some apps
*/
type SumInfo struct {
	SHA256   string `koanf,json:"SHA256"`
	UpdateAt string `koanf,json:"UpdatAt"`
}

type SumChecker struct {
	InfoList *map[string]*SumInfo `koanf,json:"InfoList"`
	conf     *config.GVConfig
	fetcher  *request.Fetcher
}

func NewSumChecker(cnf *config.GVConfig) (s *SumChecker) {
	s = &SumChecker{InfoList: &map[string]*SumInfo{}, conf: cnf}
	s.fetcher = request.NewFetcher()
	return
}

func (that *SumChecker) LoadInfoList() {
	that.fetcher.SetUrl(that.conf.Sum.SumFileUrls)
	that.fetcher.Timeout = time.Second * 30
	if resp := that.fetcher.Get(); resp != nil {
		content, _ := io.ReadAll(resp.RawResponse.Body)
		if err := json.Unmarshal(content, that.InfoList); err != nil {
			tui.PrintError("Download checksum file failed: ", err, " length: ", len(content))
			os.Exit(1)
		}
	}
}

func (that *SumChecker) parseRemoteFilename(dUrl string) string {
	sList := strings.Split(dUrl, "/")
	return sList[len(sList)-1]
}

func (that *SumChecker) IsUpdated(fPath, dUrl string) bool {
	if ok, _ := utils.PathIsExist(fPath); !ok {
		return true
	}
	rfName := that.parseRemoteFilename(dUrl)
	that.LoadInfoList()
	infoList := *that.InfoList
	if info, ok := infoList[rfName]; !ok {
		return true
	} else {
		return !utils.CheckFile(fPath, "sha256", info.SHA256)
	}
}
