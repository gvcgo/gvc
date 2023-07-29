package vctrl

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hhatto/gocloc"
	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/urfave/cli/v2"
)

const (
	FlagByFile         = "by-file"
	FlagDebug          = "debug"
	FlagSkipDuplicated = "skip-duplicated"
	FlagShowLang       = "show-lang"
	FlagSortTag        = "sort-tag"
	FlagOutputType     = "output-type"
	FlagExcludeExt     = "exclude-ext"
	FlagIncludeLang    = "include-lang"
	FlagMatch          = "match"
	FlagNotMatch       = "not-match"
	FlagMatchDir       = "match-dir"
	FlagNotMatchDir    = "not-match-dir"
)

var (
	sortTag = map[string]struct{}{
		"name":    {},
		"files":   {},
		"blank":   {},
		"comment": {},
		"code":    {},
	}
	outputType = map[string]struct{}{
		"default":   {},
		"cloc-xml":  {},
		"sloccount": {},
		"json":      {},
	}
)

type Cloc struct {
	ctx    *cli.Context
	result *gocloc.Result
}

func NewCloc(ctx *cli.Context) *Cloc {
	return &Cloc{ctx: ctx}
}

func (that *Cloc) checkFlag() bool {
	if _, ok := sortTag[that.ctx.String(FlagSortTag)]; !ok {
		return false
	}
	if _, ok := outputType[that.ctx.String(FlagOutputType)]; !ok {
		return false
	}
	return true
}

func (that *Cloc) Run() {
	if that.ctx == nil {
		return
	}
	if !that.checkFlag() {
		return
	}
	dir, _ := os.Getwd()
	paths := []string{dir}
	if that.ctx.Args().Len() > 0 {
		cargs := that.ctx.Args()
		paths = append([]string{cargs.First()}, cargs.Tail()...)
	}
	languages := gocloc.NewDefinedLanguages()
	if that.ctx.Bool(FlagShowLang) {
		fmt.Println(languages.GetFormattedString())
		return
	}
	if that.ctx.Bool(FlagByFile) && that.ctx.String(FlagOutputType) == "files" {
		fmt.Println("`--sort files` option cannot be used in conjunction with the `--by-file` option")
		os.Exit(1)
	}
	clocOpts := gocloc.NewClocOptions()

	// setup option for exclude extensions
	for _, ext := range strings.Split(that.ctx.String(FlagExcludeExt), ",") {
		e, ok := gocloc.Exts[ext]
		if ok {
			clocOpts.ExcludeExts[e] = struct{}{}
		} else {
			clocOpts.ExcludeExts[ext] = struct{}{}
		}
	}

	// directory and file matching options
	if that.ctx.String(FlagMatch) != "" {
		clocOpts.ReMatch = regexp.MustCompile(that.ctx.String(FlagMatch))
	}
	if that.ctx.String(FlagNotMatch) != "" {
		clocOpts.ReNotMatch = regexp.MustCompile(that.ctx.String(FlagNotMatch))
	}
	if that.ctx.String(FlagMatchDir) != "" {
		clocOpts.ReMatchDir = regexp.MustCompile(that.ctx.String(FlagMatchDir))
	}
	if that.ctx.String(FlagNotMatchDir) != "" {
		clocOpts.ReNotMatchDir = regexp.MustCompile(that.ctx.String(FlagNotMatchDir))
	}

	// setup option for include languages
	for _, lang := range strings.Split(that.ctx.String(FlagIncludeLang), ",") {
		if _, ok := languages.Langs[lang]; ok {
			clocOpts.IncludeLangs[lang] = struct{}{}
		}
	}

	clocOpts.Debug = that.ctx.Bool(FlagDebug)
	clocOpts.SkipDuplicated = that.ctx.Bool(FlagSkipDuplicated)

	processor := gocloc.NewProcessor(languages, clocOpts)
	var err error
	that.result, err = processor.Analyze(paths)
	if err != nil {
		fmt.Printf("fail gocloc analyze. error: %v\n", err)
		return
	}
	that.WriteResult()
}

const (
	commonHeader           string = "files          blank        comment           code"
	defaultOutputSeparator string = "-------------------------------------------------------------------------" +
		"-------------------------------------------------------------------------" +
		"-------------------------------------------------------------------------"
	OutputTypeDefault   string = "default"
	OutputTypeClocXML   string = "cloc-xml"
	OutputTypeSloccount string = "sloccount"
	OutputTypeJSON      string = "json"
)

var rowLen = 79

func (that *Cloc) WriteHeader() {
	if that.result == nil {
		return
	}
	maxPathLen := that.result.MaxPathLength
	headerLen := 28
	header := "Language"

	if that.ctx.Bool(FlagByFile) {
		headerLen = maxPathLen + 1
		rowLen = maxPathLen + len(commonHeader) + 2
		header = "File"
	}
	if that.ctx.String(FlagOutputType) == OutputTypeDefault {
		tui.Cyan(fmt.Sprintf("%.[2]*[1]s\n", defaultOutputSeparator, rowLen))
		tui.Green(fmt.Sprintf("%-[2]*[1]s %[3]s\n", header, headerLen, commonHeader))
		tui.Cyan(fmt.Sprintf("%.[2]*[1]s\n", defaultOutputSeparator, rowLen))
	}
}

func (that *Cloc) WriteFooter() {
	if that.result == nil {
		return
	}
	total := that.result.Total
	maxPathLen := that.result.MaxPathLength

	if that.ctx.String(FlagOutputType) == OutputTypeDefault {
		tui.Cyan(fmt.Sprintf("%.[2]*[1]s\n", defaultOutputSeparator, rowLen))
		if that.ctx.Bool(FlagByFile) {
			tui.Magenta(fmt.Sprintf("%-[1]*[2]v %6[3]v %14[4]v %14[5]v %14[6]v\n",
				maxPathLen, "TOTAL", total.Total, total.Blanks, total.Comments, total.Code))
		} else {
			tui.Magenta(fmt.Sprintf("%-27v %6v %14v %14v %14v\n",
				"TOTAL", total.Total, total.Blanks, total.Comments, total.Code))
		}
		tui.Cyan(fmt.Sprintf("%.[2]*[1]s\n", defaultOutputSeparator, rowLen))
	}
}

func (that *Cloc) writeResultWithByFile() {
	clocFiles := that.result.Files
	total := that.result.Total
	maxPathLen := that.result.MaxPathLength

	var sortedFiles gocloc.ClocFiles
	for _, file := range clocFiles {
		sortedFiles = append(sortedFiles, *file)
	}
	switch that.ctx.String(FlagSortTag) {
	case "name":
		sortedFiles.SortByName()
	case "comment":
		sortedFiles.SortByComments()
	case "blank":
		sortedFiles.SortByBlanks()
	default:
		sortedFiles.SortByCode()
	}

	switch that.ctx.String(FlagOutputType) {
	case OutputTypeClocXML:
		t := gocloc.XMLTotalFiles{
			Code:    total.Code,
			Comment: total.Comments,
			Blank:   total.Blanks,
		}
		f := &gocloc.XMLResultFiles{
			Files: sortedFiles,
			Total: t,
		}
		xmlResult := gocloc.XMLResult{
			XMLFiles: f,
		}
		xmlResult.Encode()
	case OutputTypeSloccount:
		for _, file := range sortedFiles {
			p := ""
			if strings.HasPrefix(file.Name, "./") || string(file.Name[0]) == "/" {
				splitPaths := strings.Split(file.Name, string(os.PathSeparator))
				if len(splitPaths) >= 3 {
					p = splitPaths[1]
				}
			}
			fmt.Printf("%v\t%v\t%v\t%v\n",
				file.Code, file.Lang, p, file.Name)
		}
	case OutputTypeJSON:
		jsonResult := gocloc.NewJSONFilesResultFromCloc(total, sortedFiles)
		buf, err := json.Marshal(jsonResult)
		if err != nil {
			fmt.Println(err)
			panic("json marshal error")
		}
		os.Stdout.Write(buf)
	default:
		for _, file := range sortedFiles {
			clocFile := file
			fmt.Printf("%-[1]*[2]s %21[3]v %14[4]v %14[5]v\n",
				maxPathLen, file.Name, clocFile.Blanks, clocFile.Comments, clocFile.Code)
		}
	}
}

func (that *Cloc) WriteResult() {
	// write header
	that.WriteHeader()

	clocLangs := that.result.Languages
	total := that.result.Total

	if that.ctx.Bool(FlagByFile) {
		that.writeResultWithByFile()
	} else {
		var sortedLanguages gocloc.Languages
		for _, language := range clocLangs {
			if len(language.Files) != 0 {
				sortedLanguages = append(sortedLanguages, *language)
			}
		}
		switch that.ctx.String(FlagSortTag) {
		case "name":
			sortedLanguages.SortByName()
		case "files":
			sortedLanguages.SortByFiles()
		case "comment":
			sortedLanguages.SortByComments()
		case "blank":
			sortedLanguages.SortByBlanks()
		default:
			sortedLanguages.SortByCode()
		}

		switch that.ctx.String(FlagOutputType) {
		case OutputTypeClocXML:
			xmlResult := gocloc.NewXMLResultFromCloc(total, sortedLanguages, gocloc.XMLResultWithLangs)
			xmlResult.Encode()
		case OutputTypeJSON:
			jsonResult := gocloc.NewJSONLanguagesResultFromCloc(total, sortedLanguages)
			buf, err := json.Marshal(jsonResult)
			if err != nil {
				fmt.Println(err)
				panic("json marshal error")
			}
			os.Stdout.Write(buf)
		default:
			for _, language := range sortedLanguages {
				tui.Gray(fmt.Sprintf("%-27v %6v %14v %14v %14v\n",
					language.Name, len(language.Files), language.Blanks, language.Comments, language.Code))
			}
		}
	}

	// write footer
	that.WriteFooter()
}
