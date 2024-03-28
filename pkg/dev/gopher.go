package dev

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/gvcgo/goutils/pkgs/gtea/confirm"
	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/gtea/selector"
	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/goutils/pkgs/koanfer"
)

/*
1. go build
*/
func isGolangInstalled() bool {
	_, err := gutils.ExecuteSysCommand(true, "", "go", "version")
	return err == nil
}

type GoBuildArchOS struct {
	ArchOSList []string `koanf:"arch_os_list"`
	Compress   bool     `koanf:"compress"`
	BuildArgs  []string `koanf:"build_args"`
}

func getGoDistlist() []string {
	out, _ := gutils.ExecuteSysCommand(true, "go", "tool", "dist", "list")
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

func ShowGoDistlist() {
	result := getGoDistlist()
	fc := gprint.NewFadeColors(result)
	fc.Println()
}

func zipDir(src, dst, binName string) (err error) {
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

func findCompiledBinary(binaryStoreDir string) (bPath string) {
	gprint.PrintInfo("binary dir: %s", binaryStoreDir)
	if fList, err := os.ReadDir(binaryStoreDir); err == nil {
		for _, f := range fList {
			if !f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
				return filepath.Join(binaryStoreDir, f.Name())
			}
		}
	}
	return
}

func getBuildArgs(buildArgs []string, binaryStoreDir string) (r []string) {
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

func build(buildArgs []string, buildBaseDir, archOS string, toGzip bool) {
	gprint.PrintInfo(fmt.Sprintf("Compiling for %s...", archOS))
	dirName := strings.ReplaceAll(archOS, "/", "-")
	infoList := strings.Split(archOS, "/")
	if len(infoList) == 2 {
		pOs, pArch := infoList[0], infoList[1]
		binaryStoreDir := filepath.Join(buildBaseDir, dirName)
		if ok, _ := gutils.PathIsExist(binaryStoreDir); !ok {
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
			if ok, _ := gutils.PathIsExist(lastArg); ok {
				targetDir = lastArg
				buildArgs = buildArgs[:len(buildArgs)-1]
			}
		}

		if !strings.Contains(strings.Join(buildArgs, " "), "-ldflags") {
			cmdArgs = append(cmdArgs, "-ldflags", `-s -w`)
		}

		bArgs := getBuildArgs(buildArgs, binaryStoreDir)
		cmdArgs = append(cmdArgs, bArgs...)

		if targetDir != "" {
			cmdArgs = append(cmdArgs, targetDir)
		}

		if _, err := gutils.ExecuteSysCommand(false, "", cmdArgs...); err != nil {
			gprint.PrintError("%+v", err)
		} else if toGzip {
			gprint.PrintSuccess(fmt.Sprintf("Compilation for %s succeeded.", archOS))
			binPath := findCompiledBinary(binaryStoreDir)
			binName := filepath.Base(binPath)
			binSuffix := path.Ext(binPath)
			name := binName
			if binName != binSuffix {
				name = strings.TrimSuffix(binName, binSuffix)
			}
			tarFilePath := strings.Join([]string{buildBaseDir, fmt.Sprintf(`%s_%s.zip`, name, dirName)}, string(filepath.Separator))
			if ok, _ := gutils.PathIsExist(tarFilePath); ok {
				os.RemoveAll(tarFilePath)
			}

			if err := zipDir(binPath, tarFilePath, binName); err != nil {
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
func handleBuildArgs(buildArgs ...string) (args []string) {
	reg := regexp.MustCompile(`(\$\(.+?\))`)
	for _, a := range buildArgs {
		toExpand := reg.FindAll([]byte(a), -1)
		for _, b := range toExpand {
			if len(b) <= 0 {
				continue
			}

			cmd := strings.TrimLeft(strings.TrimRight(string(b), ")"), "$(")
			cmdArgs := strings.Split(cmd, " ")
			if output, err := gutils.ExecuteSysCommand(true, "", cmdArgs...); err == nil {
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

func Build(args ...string) {
	if ok := isGolangInstalled(); !ok {
		gprint.PrintError("Cannot find a go compiler.")
		return
	}

	if ok, _ := gutils.PathIsExist("go.mod"); !ok {
		gprint.PrintError("Cannot find go.mod. Please check your present working directory.")
		return
	}

	buildDir := "build"
	if ok, _ := gutils.PathIsExist(buildDir); !ok {
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
	if ok, _ := gutils.PathIsExist(buildConfig); ok {
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
		for _, osArch := range getGoDistlist() {
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
		buildArgs := handleBuildArgs(bConf.BuildArgs...)
		build(buildArgs, buildDir, archOS, bConf.Compress)
		alreadyBuilt[archOS] = struct{}{}
	}
}

/*
Renames local go module.
*/
func getOldModuleName(moduleDir string) string {
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

func renameModule(pathStr, oldName, newName string, isDir bool) {
	if isDir {
		if eList, err := os.ReadDir(pathStr); err == nil {
			for _, entry := range eList {
				renameModule(filepath.Join(pathStr, entry.Name()), oldName, newName, entry.IsDir())
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

func RenameLocalModule(moduleDir, newName string) {
	oldName := getOldModuleName(moduleDir)
	if oldName == "" {
		return
	}
	renameModule(moduleDir, oldName, newName, true)
}

/*
TODO: install binaries.
2. grpc

3. goctrl gf

4. stress test

5. fmt lint etc.

6. neobox
*/
