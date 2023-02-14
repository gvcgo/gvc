package vctrl

import (
	"fmt"
	"net/url"
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
	Version   map[string][]*Package
	Doc       *goquery.Document
	Conf      *config.Conf
	ParsedUrl *url.URL
}

func NewGoVersion() (gv *GoVersion) {
	return &GoVersion{
		Version:    make(map[string][]*Package, 50),
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
		if resp := that.Download(); resp != nil {
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
		that.Version[vname] = that.findPackages(div.Find("table").First())
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
		that.Version[vname] = that.findPackages(div.Find("table").First())
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
		that.Version[vname] = that.findPackages(div.Find("table").First())
	})
	return nil
}

func (that *GoVersion) AllVersions() (err error) {
	that.GetDoc()
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

func (that *GoVersion) Run() {
	that.AllVersions()
	fmt.Println(that.Version)
}
