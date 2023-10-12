package vctrl

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mholt/archiver/v3"
	"github.com/moqsien/goutils/pkgs/gtea/confirm"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/gtea/gtable"
	"github.com/moqsien/goutils/pkgs/gtea/selector"
	"github.com/moqsien/goutils/pkgs/koanfer"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/sorts"
)

type GoPackage struct {
	Url       string
	AliUrl    string
	FileName  string
	Kind      string
	OS        string
	Arch      string
	Size      string
	Checksum  string
	CheckType string
}

type GoVersion struct {
	Versions  map[string][]*GoPackage
	Doc       *goquery.Document
	Conf      *config.GVConfig
	ParsedUrl *url.URL
	env       *utils.EnvsHandler
	fetcher   *request.Fetcher
}

func NewGoVersion() (gv *GoVersion) {
	gv = &GoVersion{
		Versions: make(map[string][]*GoPackage, 50),
		Conf:     config.New(),
		env:      utils.NewEnvsHandler(),
		fetcher:  request.NewFetcher(),
	}
	gv.initeDirs()
	gv.env.SetWinWorkDir(config.GVCDir)
	return
}

func (that *GoVersion) initeDirs() {
	utils.MakeDirs(config.DefaultGoRoot, config.GoTarFilesPath, config.GoUnTarFilesPath)
}

func (that *GoVersion) getDoc() {
	if len(that.Conf.Go.CompilerUrls) > 0 {
		itemList := selector.NewItemList()
		itemList.Add("from go.dev", that.Conf.Go.CompilerUrls[1])
		itemList.Add("from golang.google.cn", that.Conf.Go.CompilerUrls[0])
		sel := selector.NewSelector(
			itemList,
			selector.WithTitle("Choose a resource to download:"),
			selector.WithEnbleInfinite(true),
			selector.WidthEnableMulti(false),
			selector.WithWidth(20),
			selector.WithHeight(10),
		)
		sel.Run()
		val := sel.Value()[0]
		that.fetcher.Url = val.(string)

		var err error
		if that.ParsedUrl, err = url.Parse(that.fetcher.Url); err != nil {
			gprint.PrintError("%+v", err)
			os.Exit(1)
		}
		that.fetcher.Timeout = 30 * time.Second
		if resp := that.fetcher.Get(); resp != nil {
			var err error
			that.Doc, err = goquery.NewDocumentFromReader(resp.RawBody())
			if err != nil {
				gprint.PrintError(fmt.Sprintf("Parse page errored: %+v", err))
			}
			if that.Doc == nil {
				gprint.PrintError(fmt.Sprintf("Cannot parse html for %s", that.fetcher.Url))
				os.Exit(1)
			}
		}
	}
}

func (that *GoVersion) findPackages(table *goquery.Selection) (pkgs []*GoPackage) {
	alg := strings.TrimSuffix(table.Find("thead").Find("th").Last().Text(), " Checksum")

	table.Find("tr").Not(".first").Each(func(j int, tr *goquery.Selection) {
		td := tr.Find("td")
		href := td.Eq(0).Find("a").AttrOr("href", "")
		aliUrl := href
		if strings.HasPrefix(href, "/") { // relative paths
			href = fmt.Sprintf("%s://%s%s", that.ParsedUrl.Scheme, that.ParsedUrl.Host, href)
			fnameList := strings.Split(aliUrl, "/")
			fname := fnameList[len(fnameList)-1]
			if strings.Contains(fname, ".") {
				aliUrl = fmt.Sprintf("%s%s", that.Conf.Go.AliRepoUrl, fname)
			} else {
				aliUrl = ""
			}
		} else {
			aliUrl = ""
		}
		pkgs = append(pkgs, &GoPackage{
			FileName:  td.Eq(0).Find("a").Text(),
			Url:       href,
			AliUrl:    aliUrl,
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
	if that.Doc == nil {
		that.getDoc()
	}
	label := that.Doc.Find("#unstable")
	if label == nil {
		return false
	}
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
		that.getDoc()
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

func (that *GoVersion) filterVersionsForCurrentPlatform() (vList []string) {
	for key, value := range that.Versions {
	INNER:
		for _, p := range value {
			if p.OS == runtime.GOOS && p.Arch == runtime.GOARCH {
				vList = append(vList, key)
				break INNER
			}
		}
	}
	return
}

func (that *GoVersion) ShowRemoteVersions(arg string) {
	if that.Doc == nil {
		that.getDoc()
	}
	switch arg {
	case ShowAll:
		if err := that.AllVersions(); err == nil {
			fc := gprint.NewFadeColors(sorts.SortGoVersion(that.filterVersionsForCurrentPlatform()))
			fc.Println()
		}
	case ShowStable:
		if err := that.StableVersions(); err == nil {
			fc := gprint.NewFadeColors(sorts.SortGoVersion(that.filterVersionsForCurrentPlatform()))
			fc.Println()
		}
	case ShowUnstable:
		if err := that.UnstableVersions(); err == nil {
			fc := gprint.NewFadeColors(sorts.SortGoVersion(that.filterVersionsForCurrentPlatform()))
			fc.Println()
		}
	default:
		gprint.PrintWarning(fmt.Sprintf("Unknown show type: %s", arg))
	}
}

func (that *GoVersion) findPackage(version string, kind ...string) (p *GoPackage) {
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

func (that *GoVersion) download(version string) (r string) {
	p := that.findPackage(version)
	if p != nil {
		cfm := confirm.NewConfirm(confirm.WithTitle("Use mirrors.aliyun.com/golang for download acceleration?"))
		cfm.Run()
		if cfm.Result() {
			that.fetcher.Url = p.AliUrl
			if that.fetcher.Url == "" {
				that.fetcher.Url = p.Url
			}
		}
		that.fetcher.Timeout = 900 * time.Second
		that.fetcher.SetThreadNum(4)

		fName := fmt.Sprintf("go-%s-%s.%s%s", version, p.OS, p.Arch, utils.GetExt(p.FileName))
		fpath := filepath.Join(config.GoTarFilesPath, fName)
		if size := that.fetcher.GetAndSaveFile(fpath); size > 0 {
			if ok := that.checkFile(p, fpath); ok {
				return fpath
			} else {
				os.RemoveAll(fpath)
			}
		}
	} else {
		gprint.PrintError(fmt.Sprintf("Cannot find version: %s.", version))
	}
	return
}

func (that *GoVersion) checkFile(p *GoPackage, fpath string) (r bool) {
	return utils.CheckFile(fpath, p.CheckType, p.Checksum)
}

func (that *GoVersion) CheckAndInitEnv() {
	// if GOPROXY has already been set, then use the current one.
	gp := os.Getenv("GOPROXY")
	if gp == "" {
		gp = that.Conf.Go.Proxies[0]
	}
	gpath := os.Getenv("GOPATH")
	if gpath == "" {
		gpath = config.DefaultGoPath
	}
	if runtime.GOOS != utils.Windows {
		goEnv := fmt.Sprintf(utils.GoEnv,
			config.DefaultGoRoot,
			gpath,
			filepath.Join(gpath, "bin"),
			gp,
			fmt.Sprintf("%s:%s:$PATH", "$GOPATH/bin", "$GOROOT/bin"))
		that.env.UpdateSub(utils.SUB_GO, goEnv)
	} else {
		envarList := map[string]string{
			"GOROOT":  config.DefaultGoRoot,
			"GOPATH":  gpath,
			"GOBIN":   filepath.Join(gpath, "bin"),
			"GOPROXY": gp,
			"PATH": fmt.Sprintf("%s;%s", filepath.Join(gpath, "bin"),
				filepath.Join(config.DefaultGoRoot, "bin")),
		}
		that.env.SetEnvForWin(envarList)
	}
	if ok, _ := utils.PathIsExist(config.DefaultGoPath); !ok {
		os.MkdirAll(config.DefaultGoPath, os.ModePerm)
	}
}

func (that *GoVersion) UseVersion(version string) {
	untarfile := filepath.Join(config.GoUnTarFilesPath, version)
	if ok, _ := utils.PathIsExist(untarfile); !ok {
		if tarfile := that.download(version); tarfile != "" {
			if err := archiver.Unarchive(tarfile, untarfile); err != nil {
				os.RemoveAll(untarfile)
				gprint.PrintError(fmt.Sprintf("Unarchive failed: %+v.", err))
				return
			}
		} else {
			// version does not exist.
			os.RemoveAll(untarfile)
			return
		}
	}
	if ok, _ := utils.PathIsExist(config.DefaultGoRoot); ok {
		os.RemoveAll(config.DefaultGoRoot)
	}
	if err := utils.MkSymLink(filepath.Join(untarfile, "go"), config.DefaultGoRoot); err != nil {
		gprint.PrintError(fmt.Sprintf("Create link failed: %+v.", err))
		return
	}
	if !that.env.DoesEnvExist(utils.SUB_GO) {
		that.CheckAndInitEnv()
	}
	gprint.PrintSuccess(fmt.Sprintf("Use %s succeeded!", version))
}

func (that *GoVersion) getCurrent() (current string) {
	vFile := filepath.Join(config.DefaultGoRoot, "VERSION")
	if ok, _ := utils.PathIsExist(vFile); ok {
		if data, err := os.ReadFile(vFile); err == nil {
			return strings.TrimLeft(string(data), "go")
		}
	}
	return
}

func (that *GoVersion) ShowInstalled() {
	current := that.getCurrent()
	installedList, err := os.ReadDir(config.GoUnTarFilesPath)
	if err != nil {
		gprint.PrintError(fmt.Sprintf("Read dir failed: %+v", err))
		return
	}
	for _, v := range installedList {
		if v.IsDir() {
			if strings.Contains(current, v.Name()) {
				gprint.Yellow("%s <Current>", v.Name())
				continue
			}
			gprint.Cyan(v.Name())
		}
	}
}

func (that *GoVersion) parseTarFileName(name string) (v string) {
	vList := strings.Split(name, "-")
	if len(vList) > 1 {
		v = vList[1]
	}
	return
}

func (that *GoVersion) RemoveUnused() {
	current := that.getCurrent()
	installedList, err := os.ReadDir(config.GoUnTarFilesPath)
	if err != nil {
		gprint.PrintError(fmt.Sprintf("Read dir failed: %+v", err))
		return
	}
	tarFiles, _ := os.ReadDir(config.GoTarFilesPath)
	for _, v := range installedList {
		if strings.Contains(current, v.Name()) {
			continue
		}
		os.RemoveAll(filepath.Join(config.GoUnTarFilesPath, v.Name()))
		for _, vInfo := range tarFiles {
			if v.Name() == that.parseTarFileName(vInfo.Name()) {
				os.Remove(filepath.Join(config.GoTarFilesPath, vInfo.Name()))
			}
		}
	}
}

func (that *GoVersion) RemoveVersion(version string) {
	current := that.getCurrent()
	if !strings.Contains(current, version) {
		tarFiles, _ := os.ReadDir(config.GoTarFilesPath)
		os.RemoveAll(filepath.Join(config.GoUnTarFilesPath, version))
		for _, vInfo := range tarFiles {
			if version == that.parseTarFileName(vInfo.Name()) {
				os.Remove(filepath.Join(config.GoTarFilesPath, vInfo.Name()))
			}
		}
	}
}

// search libraries
func (that *GoVersion) SearchLibs(name string, sortby int) {
	that.fetcher.Url = fmt.Sprintf(that.Conf.Go.SearchUrl, name)
	resp := that.fetcher.Get()
	that.Doc, _ = goquery.NewDocumentFromReader(resp.RawBody())
	itemList := make([]sorts.Item, 0)
	that.Doc.Find(".SearchSnippet").Each(func(i int, s *goquery.Selection) {
		item := &sorts.GoLibrary{}
		item.Name = s.Find(".SearchSnippet-headerContainer").Find("a").AttrOr("href", "")
		item.Name = strings.Trim(item.Name, "/")
		item.Imported, _ = strconv.Atoi(s.Find(".SearchSnippet-infoLabel").Find("a").First().Find("strong").Text())
		if s.Find(".SearchSnippet-infoLabel").Find(".go-textSubtle").Eq(2).Find("strong").Length() > 1 {
			item.Version = s.Find(".SearchSnippet-infoLabel").Find(".go-textSubtle").Eq(2).Find("strong").Eq(0).Text()
		}
		item.Update = s.Find(".SearchSnippet-infoLabel").Find(".go-textSubtle").Eq(2).Find("span").First().Find("strong").Eq(0).Text()
		if strings.Contains(item.Update, "day") {
			s := strings.Split(item.Update, "day")[0]
			d, _ := strconv.Atoi(strings.Trim(s, " "))
			item.UpdateAt = time.Now().Add(-time.Duration(d) * 24 * time.Hour)
		} else if strings.Contains(item.Update, ",") {
			item.UpdateAt, _ = time.Parse("Jan2,2006", strings.ReplaceAll(item.Update, " ", ""))
		} else {
			item.UpdateAt = time.Now().UTC()
		}
		item.SortType = sortby
		itemList = append(itemList, item)
	})

	result := sorts.SortGoLibs(itemList)

	columns := []gtable.Column{
		{Title: "Url", Width: 60},
		{Title: "Version", Width: 40},
		{Title: "ImportedBy", Width: 15},
		{Title: "UpdatedAt", Width: 25},
	}

	rows := []gtable.Row{}

	for _, v := range result {
		rows = append(rows, gtable.Row{
			gprint.CyanStr(v.Name),
			gprint.GreenStr(v.Version),
			gprint.YellowStr("%d", v.Imported),
			v.Update,
		})
	}

	t := gtable.NewTable(
		gtable.WithColumns(columns),
		gtable.WithRows(rows),
		gtable.WithFocused(true),
		gtable.WithHeight(35),
		gtable.WithWidth(150),
	)
	t.Run()
}

/*
build programs in go for multi-platforms
*/
type GoBuildArchOS struct {
	ArchOSList []string `koanf:"arch_os_list"`
	Compress   bool     `koanf:"compress"`
	BuildArgs  []string `koanf:"build_args"`
}

func (that *GoVersion) getGoDistlist() []string {
	out, _ := utils.ExecuteSysCommand(true, "go", "tool", "dist", "list")
	commonlyUsed := map[string]struct{}{
		"darwin/amd64":  {},
		"darwin/arm64":  {},
		"linux/amd64":   {},
		"linux/arm64":   {},
		"windows/amd64": {},
		"windows/arm64": {},
	}
	commonlyUsedList := []string{}
	otherList := []string{}
	archOSList := strings.Split(out.String(), "\n")
	for _, v := range archOSList {
		v = strings.ReplaceAll(strings.Trim(v, "\r"), " ", "")
		if _, ok := commonlyUsed[v]; ok {
			commonlyUsedList = append(commonlyUsedList, v)
		} else {
			otherList = append(otherList, v)
		}
	}
	return append(commonlyUsedList, otherList...)
}

func (that *GoVersion) ShowGoDistlist() {
	result := that.getGoDistlist()
	fc := gprint.NewFadeColors(result)
	fc.Println()
}

func (that *GoVersion) zip(src, dst, binName string) (err error) {
	fr, err := os.Open(src)
	if err != nil {
		return
	}
	defer fr.Close()

	info, err := fr.Stat()
	if err != nil || info.IsDir() {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	fw, err := os.Create(dst)
	if err != nil {
		return
	}
	defer fw.Close()
	header.Name = binName
	header.Method = zip.Deflate
	zw := zip.NewWriter(fw)
	writer, err := zw.CreateHeader(header)
	if err != nil {
		return
	}
	defer zw.Close()

	if _, err = io.Copy(writer, fr); err != nil {
		return
	}
	return nil
}

func (that *GoVersion) findCompiledBinary(binaryStoreDir string) (bPath string) {
	if fList, err := os.ReadDir(binaryStoreDir); err == nil {
		for _, f := range fList {
			if !f.IsDir() {
				return filepath.Join(binaryStoreDir, f.Name())
			}
		}
	}
	return
}

func (that *GoVersion) getBuildArgs(buildArgs []string, binaryStoreDir string) (r []string) {
	hasOuput := false
	for _, v := range buildArgs {
		if v == "-o" {
			hasOuput = true
		} else if v == "." {
			continue
		}
		r = append(r, v)
	}
	if !hasOuput {
		r = append(r, "-o", binaryStoreDir)
	}
	return
}

func (that *GoVersion) build(buildArgs []string, buildBaseDir, archOS string, toGzip bool) {
	gprint.PrintInfo(fmt.Sprintf("Compiling for %s...", archOS))
	dirName := strings.ReplaceAll(archOS, "/", "-")
	infoList := strings.Split(archOS, "/")
	if len(infoList) == 2 {
		pOs, pArch := infoList[0], infoList[1]
		binaryStoreDir := filepath.Join(buildBaseDir, dirName)
		if ok, _ := utils.PathIsExist(binaryStoreDir); !ok {
			if err := os.MkdirAll(binaryStoreDir, os.ModePerm); err != nil {
				gprint.PrintError("%+v", err)
				return
			}
		}
		os.Setenv("GOOS", pOs)
		os.Setenv("GOARCH", pArch)
		cmdArgs := []string{"go", "build"}

		var targetDir string
		if len(buildArgs) > 0 {
			lastArg := buildArgs[len(buildArgs)-1]
			if ok, _ := utils.PathIsExist(lastArg); ok {
				targetDir = lastArg
				buildArgs = buildArgs[:len(buildArgs)-1]
			}
		}

		if !strings.Contains(strings.Join(buildArgs, " "), "-ldflags") {
			cmdArgs = append(cmdArgs, "-ldflags", `-s -w`)
		}

		bArgs := that.getBuildArgs(buildArgs, binaryStoreDir)
		cmdArgs = append(cmdArgs, bArgs...)

		if targetDir != "" {
			cmdArgs = append(cmdArgs, targetDir)
		}

		if _, err := utils.ExecuteSysCommand(false, cmdArgs...); err != nil {
			gprint.PrintError("%+v", err)
		} else if toGzip {
			gprint.PrintSuccess(fmt.Sprintf("Compilation for %s succeeded.", archOS))
			binPath := that.findCompiledBinary(binaryStoreDir)
			nList := strings.Split(binPath, string(filepath.Separator))
			binName := nList[len(nList)-1]
			binSuffix := path.Ext(binPath)
			name := strings.TrimSuffix(binName, binSuffix)
			tarFilePath := strings.Join([]string{buildBaseDir, fmt.Sprintf(`%s_%s.zip`, name, dirName)}, string(filepath.Separator))
			if ok, _ := utils.PathIsExist(tarFilePath); ok {
				os.RemoveAll(tarFilePath)
			}

			if err := that.zip(binPath, tarFilePath, binName); err != nil {
				gprint.PrintError("%+v", err)
			} else {
				gprint.PrintSuccess(fmt.Sprintf("Compression for %s succeeded.", archOS))
			}
		}
	} else {
		gprint.PrintError(archOS)
	}
}

// parse args by executing shell commands
func (that *GoVersion) handleBuildArgs(buildArgs ...string) (args []string) {
	var reg = regexp.MustCompile(`(\$\(.+?\))`)
	for _, a := range buildArgs {
		toExpand := reg.FindAll([]byte(a), -1)
		for _, b := range toExpand {
			if len(b) <= 0 {
				continue
			}

			cmd := strings.TrimLeft(strings.TrimRight(string(b), ")"), "$(")
			cmdArgs := strings.Split(cmd, " ")
			if output, err := utils.ExecuteSysCommand(true, cmdArgs...); err == nil {
				result := strings.TrimRight(output.String(), "\n")
				a = strings.Replace(a, string(b), result, 1)
			} else {
				gprint.PrintError("%+v", err)
				os.Exit(1)
			}
		}
		args = append(args, a)
	}
	return
}

func (that *GoVersion) Build(args ...string) {
	goRoot := os.Getenv("GOROOT")
	if ok, _ := utils.PathIsExist(goRoot); !ok {
		gprint.PrintError("Cannot find a go compiler.")
		gprint.PrintInfo(`You can install a go compiler using gvc. See help info by "gvc go help".`)
		return
	}

	if ok, _ := utils.PathIsExist("go.mod"); !ok {
		gprint.PrintError("Cannot find go.mod file. Please check your present working directory.")
		return
	}

	buildDir := "build"
	if ok, _ := utils.PathIsExist(buildDir); !ok {
		if err := os.MkdirAll(buildDir, os.ModePerm); err != nil {
			gprint.PrintError("%+v", err)
			return
		}
	}

	buildConfig := filepath.Join(buildDir, "build.json")
	kfer, _ := koanfer.NewKoanfer(buildConfig)
	bConf := &GoBuildArchOS{BuildArgs: []string{}}
	if len(args) > 0 && len(bConf.BuildArgs) == 0 {
		for idx, v := range args {
			value := v
			if value == "-ldflags" && len(args) > idx+1 {
				args[idx+1] = args[idx+1] + " -s -w"
			}
			if strings.Contains(value, "#(") {
				value = strings.ReplaceAll(value, "#(", "$(")
			}
			bConf.BuildArgs = append(bConf.BuildArgs, value)
		}
	}
	if ok, _ := utils.PathIsExist(buildConfig); ok {
		kfer.Load(bConf)
		kfer.Save(bConf)
	} else {
		itemList := selector.NewItemList()
		itemList.Add("Only for current platform", []string{fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)})
		itemList.Add("Commanly used pc platforms[Mac|Win|Linux/amd64|arm64]", []string{
			"darwin/amd64", "darwin/arm64",
			"linux/amd64", "linux/arm64",
			"windows/amd64", "windows/arm64",
		})
		for _, osArch := range that.getGoDistlist() {
			if osArch != "" {
				itemList.Add(osArch, []string{osArch})
			}
		}
		sel := selector.NewSelector(
			itemList,
			selector.WithTitle("Choose arch and os for compilation:"),
			selector.WidthEnableMulti(true),
			selector.WithEnbleInfinite(true),
			selector.WithWidth(40),
			selector.WithHeight(20),
		)
		sel.Run()
		list := sel.Value()
		bConf.ArchOSList = []string{}
		for _, val := range list {
			bConf.ArchOSList = append(bConf.ArchOSList, val.([]string)...)
		}

		cfm := confirm.NewConfirm(confirm.WithTitle("To compress binaries or not?"))
		cfm.Run()
		bConf.Compress = cfm.Result()
		kfer.Save(bConf)
	}

	alreadyBuilt := map[string]struct{}{}
	for _, archOS := range bConf.ArchOSList {
		if _, ok := alreadyBuilt[archOS]; ok {
			continue
		}
		buildArgs := that.handleBuildArgs(bConf.BuildArgs...)
		that.build(buildArgs, buildDir, archOS, bConf.Compress)
		alreadyBuilt[archOS] = struct{}{}
	}
}

// Rename local go module
func (that *GoVersion) getOldModuleName(moduleDir string) string {
	var (
		modFileName = "go.mod"
		keyword     = "module"
	)
	if eList, err := os.ReadDir(moduleDir); err == nil {
		for _, entry := range eList {
			if entry.Name() == modFileName {
				// open the file
				file, err := os.Open(filepath.Join(moduleDir, entry.Name()))
				if err != nil {
					gprint.PrintError("%+v", err)
					return ""
				}
				defer file.Close()
				fileScanner := bufio.NewScanner(file)
				for fileScanner.Scan() {
					t := fileScanner.Text()
					if strings.Contains(t, keyword) {
						sList := strings.Split(t, keyword)
						return strings.TrimSpace(sList[1])
					}
				}
				if err := fileScanner.Err(); err != nil {
					gprint.PrintError("%+v", err)
				}
			}
		}
	}
	gprint.PrintError(fmt.Sprintf("Can not find module name in [%s].", filepath.Join(moduleDir, modFileName)))
	return ""
}

func (that *GoVersion) renameModule(pathStr, oldName, newName string, isDir bool) {
	if isDir {
		if eList, err := os.ReadDir(pathStr); err == nil {
			for _, entry := range eList {
				that.renameModule(filepath.Join(pathStr, entry.Name()), oldName, newName, entry.IsDir())
			}
		}
	} else {
		if strings.HasSuffix(pathStr, "go.mod") || strings.HasSuffix(pathStr, ".go") {
			if content, err := os.ReadFile(pathStr); err == nil {
				newStr := strings.ReplaceAll(string(content), oldName, newName)
				os.WriteFile(pathStr, []byte(newStr), os.ModePerm)
			}
		}
	}
}

func (that *GoVersion) RenameLocalModule(moduleDir, newName string) {
	oldName := that.getOldModuleName(moduleDir)
	if oldName == "" {
		return
	}
	that.renameModule(moduleDir, oldName, newName, true)
}
