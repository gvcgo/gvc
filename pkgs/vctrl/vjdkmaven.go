package vctrl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mholt/archiver/v3"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/sorts"
)

type MavenPackage struct {
	Version     string
	Url         string
	ChecksumUrl string
	FileName    string
}

type MavenVersion struct {
	Versions map[string]*MavenPackage
	Doc      *goquery.Document
	Conf     *config.GVConfig
	fetcher  *request.Fetcher
	env      *utils.EnvsHandler
}

func NewMavenVersion() (mv *MavenVersion) {
	mv = &MavenVersion{
		Versions: make(map[string]*MavenPackage, 20),
		Conf:     config.New(),
		fetcher:  request.NewFetcher(),
		env:      utils.NewEnvsHandler(),
	}
	mv.initeDirs()
	mv.env.SetWinWorkDir(config.GVCDir)
	return
}

func (that *MavenVersion) initeDirs() {
	utils.MakeDirs(config.MavenRoot, config.MavenTarFilePath, config.MavenUntarFilePath, config.JavaLocalRepoPath)
}

func (that *MavenVersion) getVs(vn string) {
	var mUrl string
	switch vn {
	case "4.":
		mUrl = that.Conf.Maven.ApacheUrl4
	default:
		mUrl = that.Conf.Maven.ApacheUrl3
	}
	if !utils.VerifyUrls(mUrl) {
		return
	}
	that.Doc = nil
	that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(mUrl)
	if resp := that.fetcher.Get(); resp != nil {
		that.Doc, _ = goquery.NewDocumentFromReader(resp.RawBody())
	}
	if that.Doc != nil {
		that.Doc.Find("a").Each(func(i int, s *goquery.Selection) {
			link := s.AttrOr("href", "")
			if strings.HasPrefix(link, vn) {
				p := &MavenPackage{}
				p.Version = strings.ReplaceAll(link, "/", "")
				p.Url = fmt.Sprintf(that.Conf.Maven.UrlPattern,
					mUrl, p.Version, p.Version)
				p.ChecksumUrl = fmt.Sprintf(that.Conf.Maven.ShaUrlPattern,
					mUrl, p.Version, p.Version)
				p.FileName = fmt.Sprintf("maven-%s-bin.tar.gz", p.Version)
				that.Versions[p.Version] = p
			}
		})
	}
}

func (that *MavenVersion) getVersions() {
	if len(that.Versions) > 0 {
		return
	}
	vnList := []string{"3.", "4."}
	for _, vn := range vnList {
		that.getVs(vn)
	}
}

func (that *MavenVersion) getSha(p *MavenPackage) (shaCode string) {
	if utils.VerifyUrls(p.ChecksumUrl) {
		that.fetcher.Url = p.ChecksumUrl
		if resp := that.fetcher.Get(); resp != nil {
			content, _ := io.ReadAll(resp.RawBody())
			shaCode = string(content)
		}
	}
	return
}

func (that *MavenVersion) ShowVersions() {
	that.getVersions()
	vList := []string{}
	for k := range that.Versions {
		vList = append(vList, k)
	}
	if len(vList) > 0 {
		res := sorts.SortGoVersion(vList)
		fc := gprint.NewFadeColors(res)
		fc.Println()
	}
}

func (that *MavenVersion) download(version string) (r string) {
	that.getVersions()
	if p, ok := that.Versions[version]; ok {
		that.fetcher.Url = that.Conf.GVCProxy.WrapUrl(p.Url)
		that.fetcher.Timeout = 900 * time.Minute
		that.fetcher.SetThreadNum(8)
		fpath := filepath.Join(config.MavenTarFilePath, p.FileName)
		if size := that.fetcher.GetAndSaveFile(fpath); size > 0 {
			if ok := utils.CheckFile(fpath, "sha512", that.getSha(p)); ok {
				return fpath
			} else {
				os.RemoveAll(fpath)
			}
		} else {
			os.RemoveAll(fpath)
		}
	}
	return
}

func (that *MavenVersion) UseVersion(version string) {
	untarfile := filepath.Join(config.MavenUntarFilePath, fmt.Sprintf("maven-%s", version))
	if ok, _ := utils.PathIsExist(untarfile); !ok {
		if tarfile := that.download(version); tarfile != "" {
			if err := archiver.Unarchive(tarfile, untarfile); err != nil {
				os.RemoveAll(untarfile)
				gprint.PrintError(fmt.Sprintf("Unarchive failed: %+v", err))
				return
			}
		}
	}
	if ok, _ := utils.PathIsExist(config.MavenRoot); ok {
		os.RemoveAll(config.MavenRoot)
	}
	finder := utils.NewBinaryFinder(untarfile, "bin")
	dir := finder.String()
	if dir != "" {
		if err := utils.MkSymLink(dir, config.MavenRoot); err != nil {
			gprint.PrintError(fmt.Sprintf("Create link failed: %+v", err))
			return
		}
		if !that.env.DoesEnvExist(utils.SUB_MAVEN) {
			that.CheckAndInitEnv()
		}
		utils.RecordVersion(version, dir)
		gprint.PrintSuccess(fmt.Sprintf("Use %s succeeded!", version))
	}
}

func (that *MavenVersion) CheckAndInitEnv() {
	if runtime.GOOS != utils.Windows {
		mavenEnv := fmt.Sprintf(utils.MavenEnv,
			config.MavenRoot)
		that.env.UpdateSub(utils.SUB_MAVEN, mavenEnv)
	} else {
		envList := map[string]string{
			"MAVEN_HOME": config.MavenRoot,
			"PATH":       filepath.Join(config.MavenRoot, "bin"),
		}
		that.env.SetEnvForWin(envList)
	}
}

func (that *MavenVersion) GenSettingsFile() {
	sf := filepath.Join(config.MavenSettingsFileDir, "settings.xml")
	osf := filepath.Join(config.MavenSettingsFileDir, "settings.xml.origin")
	if ok, _ := utils.PathIsExist(config.MavenSettingsFileDir); ok {
		if ok1, _ := utils.PathIsExist(osf); !ok1 {
			utils.CopyFile(sf, osf)
		}
		os.WriteFile(sf, []byte(config.MavenSettings), 0644)
	}
}

func (that *MavenVersion) ShowInstalled() {
	if ok, _ := utils.PathIsExist(config.MavenUntarFilePath); ok {
		current := utils.ReadVersion(config.MavenRoot)
		dList, _ := os.ReadDir(config.MavenUntarFilePath)
		for _, d := range dList {
			if strings.Contains(d.Name(), "maven-") {
				version := strings.Split(d.Name(), "-")[1]
				if current == version {
					gprint.Yellow("%s <Current>", version)
					continue
				}
				gprint.Cyan(version)
			}
		}
	}
}

func (that *MavenVersion) RemoveTarFile(version string) {
	fPath := filepath.Join(config.MavenTarFilePath, fmt.Sprintf("maven-%s-bin.tar.gz", version))
	os.RemoveAll(fPath)
}

func (that *MavenVersion) RemoveVersion(version string) {
	if ok, _ := utils.PathIsExist(config.MavenUntarFilePath); ok {
		current := utils.ReadVersion(config.MavenRoot)
		dList, _ := os.ReadDir(config.MavenUntarFilePath)
		for _, d := range dList {
			if strings.Contains(d.Name(), "maven-") {
				v := strings.Split(d.Name(), "-")[1]
				if current != version && v == version {
					p := filepath.Join(config.MavenUntarFilePath, d.Name())
					os.RemoveAll(p)
					that.RemoveTarFile(version)
				}
			}
		}
	}
}

func (that *MavenVersion) RemoveUnused() {
	if ok, _ := utils.PathIsExist(config.MavenUntarFilePath); ok {
		current := utils.ReadVersion(config.MavenRoot)
		dList, _ := os.ReadDir(config.MavenUntarFilePath)
		for _, d := range dList {
			if strings.Contains(d.Name(), "maven-") {
				version := strings.Split(d.Name(), "-")[1]
				if current != version {
					p := filepath.Join(config.MavenUntarFilePath, d.Name())
					os.RemoveAll(p)
					that.RemoveTarFile(version)
				}
			}
		}
	}
}
