package vctrl

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/moqsien/gvc/pkgs/config"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Package struct {
	Url       string
	FileName  string
	Kind      string
	OS        string
	Arch      string
	Size      string
	Checksum  string
	CheckType string
}

type GoVersion struct {
	*downloader.Downloader
	Versions  map[string][]*Package
	Doc       *goquery.Document
	Conf      *config.Conf
	ParsedUrl *url.URL
}

func NewGoVersion() (gv *GoVersion) {
	return &GoVersion{
		Versions:   make(map[string][]*Package, 50),
		Conf:       config.New(),
		Downloader: &downloader.Downloader{},
	}
}

func (that *GoVersion) GetDoc() {
	if len(that.Conf.Config.Go.CompilerUrls) > 0 {
		that.Url = that.Conf.Config.Go.CompilerUrls[0]
		var err error
		if that.ParsedUrl, err = url.Parse(that.Url); err != nil {
			panic(err)
		}
		that.Timeout = 30 * time.Second
		if resp := that.GetUrl(); resp != nil {
			var err error
			that.Doc, err = goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				fmt.Println("[parse page errored] ", err)
			}
		}
	}
}

func (that *GoVersion) findPackages(table *goquery.Selection) (pkgs []*Package) {
	alg := strings.TrimSuffix(table.Find("thead").Find("th").Last().Text(), " Checksum")

	table.Find("tr").Not(".first").Each(func(j int, tr *goquery.Selection) {
		td := tr.Find("td")
		href := td.Eq(0).Find("a").AttrOr("href", "")
		if strings.HasPrefix(href, "/") { // relative paths
			href = fmt.Sprintf("%s://%s%s", that.ParsedUrl.Scheme, that.ParsedUrl.Host, href)
		}
		pkgs = append(pkgs, &Package{
			FileName:  td.Eq(0).Find("a").Text(),
			Url:       href,
			Kind:      strings.ToLower(td.Eq(1).Text()),
			OS:        utils.MapArchAndOS(td.Eq(2).Text()),
			Arch:      utils.MapArchAndOS(td.Eq(3).Text()),
			Size:      td.Eq(4).Text(),
			Checksum:  td.Eq(5).Text(),
			CheckType: alg,
		})
	})
	return pkgs
}

func (that *GoVersion) hasUnstableVersions() bool {
	return that.Doc.Find("#unstable").Length() > 0
}

func (that *GoVersion) StableVersions() (err error) {
	var divs *goquery.Selection
	if that.hasUnstableVersions() {
		divs = that.Doc.Find("#stable").NextUntil("#unstable")
	} else {
		divs = that.Doc.Find("#stable").NextUntil("#archive")
	}
	divs.Each(func(i int, div *goquery.Selection) {
		vname, ok := div.Attr("id")
		if !ok {
			return
		}
		vname = strings.TrimPrefix(vname, "go")
		that.Versions[vname] = that.findPackages(div.Find("table").First())
	})
	return nil
}

func (that *GoVersion) UnstableVersions() (err error) {
	that.Doc.Find("#unstable").NextUntil("#archive").Each(func(i int, div *goquery.Selection) {
		vname, ok := div.Attr("id")
		if !ok {
			return
		}
		vname = strings.TrimPrefix(vname, "go")
		that.Versions[vname] = that.findPackages(div.Find("table").First())
	})
	return nil
}

func (that *GoVersion) ArchivedVersions() (err error) {
	that.Doc.Find("#archive").Find("div.toggle").Each(func(i int, div *goquery.Selection) {
		vname, ok := div.Attr("id")
		if !ok {
			return
		}
		vname = strings.TrimPrefix(vname, "go")
		that.Versions[vname] = that.findPackages(div.Find("table").First())
	})
	return nil
}

func (that *GoVersion) AllVersions() (err error) {
	if that.Doc == nil {
		that.GetDoc()
	}
	err = that.StableVersions()
	if err != nil {
		return
	}
	err = that.ArchivedVersions()
	if err != nil {
		return
	}
	err = that.UnstableVersions()
	if err != nil {
		return
	}
	return
}

func (that *GoVersion) GetVersions() (vList []string) {
	for v := range that.Versions {
		vList = append(vList, v)
	}
	return vList
}

const (
	ShowAll      string = "1"
	ShowStable   string = "2"
	ShowUnstable string = "3"
)

func (that *GoVersion) ShowRemoteVersions(arg string) {
	var v *utils.VComparator
	if that.Doc == nil {
		that.GetDoc()
	}
	switch arg {
	case ShowAll:
		if err := that.AllVersions(); err == nil {
			v = utils.NewVComparator(that.GetVersions())
			fmt.Println(strings.Join(v.Order(), "  "))
		}
	case ShowStable:
		if err := that.StableVersions(); err == nil {
			v = utils.NewVComparator(that.GetVersions())
			fmt.Println(strings.Join(v.Order(), "  "))
		}
	case ShowUnstable:
		if err := that.UnstableVersions(); err == nil {
			v = utils.NewVComparator(that.GetVersions())
			fmt.Println(strings.Join(v.Order(), "  "))
		}
	default:
		fmt.Println("[Unknown show type] ", arg)
	}
}

func (that *GoVersion) findPackage(version string, kind ...string) (p *Package) {
	k := "archive"
	if len(kind) > 0 {
		k = kind[0]
	}
	that.AllVersions()
	if vList, ok := that.Versions[version]; ok {
		for _, v := range vList {
			if v.OS == runtime.GOOS && v.Arch == runtime.GOARCH && v.Kind == k {
				p = v
			}
		}
	}
	return
}

func (that *GoVersion) Download(version string) (r string) {
	p := that.findPackage(version)
	if p != nil {
		fpath := filepath.Join(config.GoTarFilesPath, p.FileName)
		that.Downloader.Url = p.Url
		that.Downloader.Timeout = 180 * time.Second
		if size := that.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size > 0 {
			that.CheckFile(p, fpath)
		}
	}
	return
}

func (that *GoVersion) CheckFile(p *Package, fpath string) (r bool) {
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Println("[Open file failed] ", err)
		return false
	}
	defer f.Close()

	var h hash.Hash
	switch strings.ToLower(p.CheckType) {
	case "sha256":
		h = sha256.New()
	case "sha1":
		h = sha1.New()
	default:
		fmt.Println("[Crypto] ", p.CheckType, " not supported.")
		return
	}

	if _, err = io.Copy(h, f); err != nil {
		fmt.Println("[Copy file failed] ", err)
		return
	}

	if p.Checksum != hex.EncodeToString(h.Sum(nil)) {
		fmt.Println("Checksum failed.")
		return
	}
	fmt.Println("Checksum successed.")
	return true
}

func (that *GoVersion) CheckEnv() {

}

func (that *GoVersion) UseVersion(version string) {

}
