package vctrl

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gookit/color"
	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/sorts"
)

type GradlePackage struct {
	Version  string
	Url      string
	Checksum string
	FileName string
}

type GradleVersion struct {
	Versions map[string]*GradlePackage
	sha      map[string]string
	Doc      *goquery.Document
	Conf     *config.GVConfig
	d        *downloader.Downloader
	c        *colly.Collector
	env      *utils.EnvsHandler
}

func NewGradleVersion() (gv *GradleVersion) {
	gv = &GradleVersion{
		Versions: make(map[string]*GradlePackage, 100),
		sha:      make(map[string]string, 100),
		Conf:     config.New(),
		d:        &downloader.Downloader{},
		c:        colly.NewCollector(),
		env:      utils.NewEnvsHandler(),
	}
	gv.initeDirs()
	return gv
}

func (that *GradleVersion) initeDirs() {
	if ok, _ := utils.PathIsExist(config.GradleRoot); !ok {
		os.RemoveAll(config.GradleRoot)
		if err := os.MkdirAll(config.GradleRoot, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.GradleTarFilePath); !ok {
		if err := os.MkdirAll(config.GradleTarFilePath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.GradleUntarFilePath); !ok {
		if err := os.MkdirAll(config.GradleUntarFilePath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
	if ok, _ := utils.PathIsExist(config.GradleInitFilePath); !ok {
		if err := os.MkdirAll(config.GradleInitFilePath, os.ModePerm); err != nil {
			fmt.Println("[mkdir Failed] ", err)
		}
	}
}

func (that *GradleVersion) getDoc() {
	gUrl := that.Conf.Gradle.OfficialUrl
	if !utils.VerifyUrls(gUrl) {
		return
	}
	that.c.OnResponse(func(r *colly.Response) {
		that.Doc, _ = goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
	})
	that.c.Visit(gUrl)
}

func (that *GradleVersion) getSha() {
	that.c = colly.NewCollector()
	cUrl := that.Conf.Gradle.OfficialCheckUrl
	if !utils.VerifyUrls(cUrl) {
		return
	}
	that.Doc = nil
	that.c.OnResponse(func(r *colly.Response) {
		that.Doc, _ = goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
	})
	that.c.Visit(cUrl)
	if that.Doc != nil {
		that.Doc.Find("h3.u-text-with-icon").Each(func(i int, s *goquery.Selection) {
			version := s.Find("a").AttrOr("id", "")
			if version == "" {
				return
			}
			shaCode := s.Next().Find("li").Eq(0).Find("code").Text()
			if shaCode != "" {
				that.sha[version] = shaCode
			}
		})
	}
}

func (that *GradleVersion) shaCode(version string) (code string) {
	if len(that.sha) == 0 {
		that.getSha()
	}
	for k, v := range that.sha {
		if strings.ReplaceAll(k, "v", "") == version {
			return v
		}
	}
	return
}

func (that *GradleVersion) getVersions() {
	if len(that.Versions) > 0 {
		return
	}
	if that.Doc == nil {
		that.getDoc()
	}
	if that.Doc != nil {
		that.Doc.Find("div.indent").Each(func(i int, s *goquery.Selection) {
			aLabel := s.Find("li").Eq(0).Find("a").Eq(0)
			p := &GradlePackage{}
			p.Url = aLabel.AttrOr("href", "")
			p.Version = aLabel.AttrOr("data-version", "")
			if p.Url == "" || p.Version == "" {
				return
			}
			p.Checksum = strings.TrimSpace(that.shaCode(p.Version))
			p.FileName = fmt.Sprintf("gradle-%s.zip", p.Version)
			that.Versions[p.Version] = p
		})
	}
}

func (that *GradleVersion) ShowVersions() {
	that.getVersions()
	vList := []string{}
	for k := range that.Versions {
		vList = append(vList, k)
	}
	res := sorts.SortGoVersion(vList)
	fmt.Println(strings.Join(res, "  "))
}

func (that *GradleVersion) download(version string) (r string) {
	that.getVersions()
	if p, ok := that.Versions[version]; ok {
		that.d.Url = p.Url
		that.d.Timeout = 30 * time.Minute
		fpath := filepath.Join(config.GradleTarFilePath, p.FileName)
		if size := that.d.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size > 0 {
			if ok := utils.CheckFile(fpath, "sha256", p.Checksum); ok {
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

func (that *GradleVersion) UseVersion(version string) {
	untarfile := filepath.Join(config.GradleUntarFilePath, fmt.Sprintf("gradle-%s", version))
	if ok, _ := utils.PathIsExist(untarfile); !ok {
		if tarfile := that.download(version); tarfile != "" {
			if err := archiver.Unarchive(tarfile, untarfile); err != nil {
				os.RemoveAll(untarfile)
				fmt.Println("[Unarchive failed] ", err)
				return
			}
		}
	}
	if ok, _ := utils.PathIsExist(config.GradleRoot); ok {
		os.RemoveAll(config.GradleRoot)
	}
	finder := utils.NewBinaryFinder(untarfile, "bin")
	dir := finder.String()
	if dir != "" {
		if err := utils.MkSymLink(dir, config.GradleRoot); err != nil {
			fmt.Println("[Create link failed] ", err)
			return
		}
		if !that.env.DoesEnvExist(utils.SUB_GRADLE) {
			that.CheckAndInitEnv()
		}
		utils.RecordVersion(version, dir)
		fmt.Println("Use", version, "succeeded!")
	}
}

func (that *GradleVersion) CheckAndInitEnv() {
	if runtime.GOOS != utils.Windows {
		gradleEnv := fmt.Sprintf(utils.GradleEnv,
			config.GradleRoot,
			config.JavaLocalRepoPath)
		that.env.UpdateSub(utils.SUB_GRADLE, gradleEnv)
	} else {
		envList := map[string]string{
			"GRADLE_HOME":      config.GradleRoot,
			"GRADLE_USER_HOME": config.JavaLocalRepoPath,
			"PATH":             filepath.Join(config.GradleRoot, "bin"),
		}
		that.env.SetEnvForWin(envList)
	}
}

func (that *GradleVersion) GenInitFile() {
	sf := filepath.Join(config.GradleInitFilePath, "init.gradle")
	osf := filepath.Join(config.GradleInitFilePath, "init.gradle.origin")
	if ok, _ := utils.PathIsExist(config.GradleInitFilePath); ok {
		if ok1, _ := utils.PathIsExist(osf); !ok1 {
			if ok2, _ := utils.PathIsExist(sf); ok2 {
				utils.CopyFile(sf, osf)
			}
		}
		os.WriteFile(sf, []byte(config.GradleInitFileContent), 0644)
	}
}

func (that *GradleVersion) ShowInstalled() {
	if ok, _ := utils.PathIsExist(config.GradleUntarFilePath); ok {
		current := utils.ReadVersion(config.GradleRoot)
		dList, _ := os.ReadDir(config.GradleUntarFilePath)
		for _, d := range dList {
			if strings.Contains(d.Name(), "gradle-") {
				version := strings.Split(d.Name(), "-")[1]
				if current == version {
					s := fmt.Sprintf("%s <Current>", version)
					color.Yellow.Println(s)
					continue
				}
				color.Cyan.Println(version)
			}
		}
	}
}

func (that *GradleVersion) RemoveTarFile(version string) {
	fPath := filepath.Join(config.GradleTarFilePath, fmt.Sprintf("gradle-%s.zip", version))
	os.RemoveAll(fPath)
}

func (that *GradleVersion) RemoveVersion(version string) {
	if ok, _ := utils.PathIsExist(config.GradleUntarFilePath); ok {
		current := utils.ReadVersion(config.GradleRoot)
		dList, _ := os.ReadDir(config.GradleUntarFilePath)
		for _, d := range dList {
			if strings.Contains(d.Name(), "gradle-") {
				v := strings.Split(d.Name(), "-")[1]
				if current != version && v == version {
					p := filepath.Join(config.GradleUntarFilePath, d.Name())
					os.RemoveAll(p)
					that.RemoveTarFile(version)
				}
			}
		}
	}
}

func (that *GradleVersion) RemoveUnused() {
	if ok, _ := utils.PathIsExist(config.GradleUntarFilePath); ok {
		current := utils.ReadVersion(config.GradleRoot)
		dList, _ := os.ReadDir(config.GradleUntarFilePath)
		for _, d := range dList {
			if strings.Contains(d.Name(), "gradle-") {
				version := strings.Split(d.Name(), "-")[1]
				if current != version {
					p := filepath.Join(config.GradleUntarFilePath, d.Name())
					os.RemoveAll(p)
					that.RemoveTarFile(version)
				}
			}
		}
	}
}
